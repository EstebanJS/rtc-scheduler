// internal/application/usecases/run_service.go
package usecases

import (
	"rtc-scheduler/internal/domain/repositories"
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

	// Verificar que haya configuración
	if !uc.configRepo.Exists() {
		uc.logger.Warn("No configuration found, skipping service execution")
		return &RunServiceOutput{
			Executed: false,
			Message:  "No configuration found",
		}, nil
	}

	// Cargar configuración
	config, err := uc.configRepo.Load()
	if err != nil {
		uc.logger.Error("Failed to load configuration", "error", err)
		return nil, err
	}

	// Verificar que esté habilitado
	if !config.Enabled {
		uc.logger.Info("Service is disabled, skipping execution")
		return &RunServiceOutput{
			Executed: false,
			Message:  "Service is disabled",
		}, nil
	}

	// Convertir configuración a schedule
	schedule, err := config.ParseToSchedule()
	if err != nil {
		uc.logger.Error("Failed to parse schedule from config", "error", err)
		return nil, err
	}

	uc.logger.Info("Executing schedule",
		"wake_time", schedule.WakeTime,
		"shutdown_time", schedule.ShutdownTime,
	)

	// Configurar alarma RTC
	if err := uc.rtcRepo.SetWakeAlarm(schedule.WakeTime); err != nil {
		uc.logger.Error("Failed to set RTC wake alarm", "error", err)
		return nil, err
	}

	// Programar apagado
	if err := uc.schedulerRepo.ScheduleShutdown(schedule.ShutdownTime); err != nil {
		// Limpiar alarma si falla el scheduling
		uc.rtcRepo.ClearWakeAlarm()
		uc.logger.Error("Failed to schedule shutdown", "error", err)
		return nil, err
	}

	uc.logger.Info("Service execution completed successfully")

	return &RunServiceOutput{
		Executed: true,
		Message:  "Schedule configured successfully",
	}, nil
}