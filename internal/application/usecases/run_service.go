// internal/application/usecases/run_service.go
package usecases

import (
	"errors"
	"fmt"
	"os"

	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/internal/infrastructure/scheduler"
	"rtc-scheduler/pkg/logger"
)

type RunServiceInput struct{}

type RunServiceOutput struct {
	Executed     bool
	Message      string
}

// RunServiceUseCase maneja la ejecución desde el servicio systemd
type RunServiceUseCase struct {
	configRepo    repositories.ConfigRepository
	rtcRepo       repositories.RTCRepository
	schedulerRepo repositories.SchedulerRepository
	logger        logger.Logger
}

func NewRunServiceUseCase(
	config repositories.ConfigRepository,
	rtc repositories.RTCRepository,
	scheduler repositories.SchedulerRepository,
	log logger.Logger,
) *RunServiceUseCase {
	return &RunServiceUseCase{
		configRepo:    config,
		rtcRepo:       rtc,
		schedulerRepo: scheduler,
		logger:        log,
	}
}

func (uc *RunServiceUseCase) Execute(input *RunServiceInput) (*RunServiceOutput, error) {
	uc.logger.Info("Running service execution")

	// Validación inicial de dependencias
	fmt.Fprintf(os.Stderr, "DEBUG: Validating dependencies...\n")
	if !uc.rtcRepo.IsAvailable() {
		errMsg := "RTC device is not available"
		fmt.Fprintln(os.Stderr, "ERROR:", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	fmt.Fprintf(os.Stderr, "DEBUG: RTC device is available\n")

	if !uc.schedulerRepo.IsAvailable() {
		errMsg := "Scheduler (at command) is not available"
		fmt.Fprintln(os.Stderr, "ERROR:", errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	fmt.Fprintf(os.Stderr, "DEBUG: Scheduler is available\n")

	// Verificar que haya configuración
	if !uc.configRepo.Exists() {
		errMsg := "No configuration found, skipping service execution"
		uc.logger.Warn(errMsg)
		// Also write to stderr for systemd logging
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errMsg)
		return &RunServiceOutput{
			Executed: false,
			Message:  "No configuration found",
		}, nil
	}

	// Cargar configuración
	fmt.Fprintf(os.Stderr, "DEBUG: Loading configuration...\n")
	config, err := uc.configRepo.Load()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load configuration: %v", err)
		uc.logger.Error("Failed to load configuration", "error", err)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errMsg)
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "DEBUG: Configuration loaded successfully\n")

	// Verificar que esté habilitado
	if !config.Enabled {
		msg := "Service is disabled, skipping execution"
		uc.logger.Info(msg)
		fmt.Fprintf(os.Stderr, "INFO: %s\n", msg)
		return &RunServiceOutput{
			Executed: false,
			Message:  "Service is disabled",
		}, nil
	}

	// Convertir configuración a schedule
	fmt.Fprintf(os.Stderr, "DEBUG: Parsing schedule from config...\n")
	schedule, err := config.ParseToSchedule()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse schedule from config: %v", err)
		uc.logger.Error("Failed to parse schedule from config", "error", err)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errMsg)
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "DEBUG: Schedule parsed successfully\n")

	fmt.Fprintf(os.Stderr, "DEBUG: Executing schedule - Wake: %s, Shutdown: %s\n",
		schedule.WakeTime.Format("2006-01-02 15:04:05"),
		schedule.ShutdownTime.Format("2006-01-02 15:04:05"))

	uc.logger.Info("Executing schedule",
		"wake_time", schedule.WakeTime,
		"shutdown_time", schedule.ShutdownTime,
	)

	// Configurar alarma RTC
	fmt.Fprintf(os.Stderr, "DEBUG: Setting RTC wake alarm...\n")
	if err := uc.rtcRepo.SetWakeAlarm(schedule.WakeTime); err != nil {
		errMsg := fmt.Sprintf("Failed to set RTC wake alarm: %v", err)
		uc.logger.Error("Failed to set RTC wake alarm", "error", err)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errMsg)
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "DEBUG: RTC wake alarm set successfully\n")

	// Programar apagado
	fmt.Fprintf(os.Stderr, "DEBUG: Scheduling shutdown...\n")
	if err := uc.schedulerRepo.ScheduleShutdown(schedule.ShutdownTime); err != nil {
		// Verificar si es un error de filesystem read-only (modo degradado)
		if errors.Is(err, scheduler.ErrFilesystemReadOnly) {
			// Modo degradado: RTC funciona, pero shutdown no se programa
			warnMsg := "Filesystem is read-only, operating in degraded mode: RTC wake alarm configured but shutdown scheduling unavailable"
			uc.logger.Warn(warnMsg, "error", err)
			fmt.Fprintf(os.Stderr, "WARNING: %s\n", warnMsg)
			fmt.Fprintf(os.Stderr, "INFO: RTC wake alarm remains configured for automatic startup\n")

			uc.logger.Info("Service execution completed in degraded mode (RTC only)")

			return &RunServiceOutput{
				Executed: true,
				Message:  "RTC wake alarm configured (shutdown scheduling unavailable due to read-only filesystem)",
			}, nil
		}

		// Para otros errores, limpiar alarma y fallar
		_ = uc.rtcRepo.ClearWakeAlarm()
		errMsg := fmt.Sprintf("Failed to schedule shutdown: %v", err)
		uc.logger.Error("Failed to schedule shutdown", "error", err)
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errMsg)
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "DEBUG: Shutdown scheduled successfully\n")

	uc.logger.Info("Service execution completed successfully")

	return &RunServiceOutput{
		Executed: true,
		Message:  "Schedule configured successfully",
	}, nil
}