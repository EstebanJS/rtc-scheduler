// internal/application/usecases/schedule_power.go
package usecases

import (
	"fmt"

	"rtc-scheduler/internal/domain/entities"
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

type SchedulePowerInput struct {
	WakeTime     string
	ShutdownTime string
	TestMode     bool
}

type SchedulePowerOutput struct {
	Schedule     *entities.Schedule
	Message      string
	TestMode     bool
}

// SchedulePowerUseCase maneja la programación de encendido/apagado
type SchedulePowerUseCase struct {
	rtcRepo       repositories.RTCRepository
	schedulerRepo repositories.SchedulerRepository
	logger        logger.Logger
}

func NewSchedulePowerUseCase(
	rtc repositories.RTCRepository,
	scheduler repositories.SchedulerRepository,
	log logger.Logger,
) *SchedulePowerUseCase {
	return &SchedulePowerUseCase{
		rtcRepo:       rtc,
		schedulerRepo: scheduler,
		logger:        log,
	}
}

func (uc *SchedulePowerUseCase) Execute(input *SchedulePowerInput) (*SchedulePowerOutput, error) {
	uc.logger.Info("Starting power scheduling",
		"wake_time", input.WakeTime,
		"shutdown_time", input.ShutdownTime,
		"test_mode", input.TestMode,
	)

	// Crear configuración
	config, err := entities.NewConfig(input.WakeTime, input.ShutdownTime, true)
	if err != nil {
		uc.logger.Error("Invalid configuration", "error", err)
		return nil, fmt.Errorf("invalid time configuration: %w", err)
	}

	// Convertir a schedule
	schedule, err := config.ParseToSchedule()
	if err != nil {
		uc.logger.Error("Failed to parse schedule", "error", err)
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	// Configurar alarma RTC para encendido
	if err := uc.rtcRepo.SetWakeAlarm(schedule.WakeTime); err != nil {
		uc.logger.Error("Failed to set RTC wake alarm", "error", err)
		return nil, fmt.Errorf("failed to set wake alarm: %w", err)
	}

	// Programar apagado
	if err := uc.schedulerRepo.ScheduleShutdown(schedule.ShutdownTime); err != nil {
		// Si falla, limpiar la alarma RTC
		uc.rtcRepo.ClearWakeAlarm() //nolint:errcheck
		uc.logger.Error("Failed to schedule shutdown", "error", err)
		return nil, fmt.Errorf("failed to schedule shutdown: %w", err)
	}

	mode := "production"
	if input.TestMode {
		mode = "test"
	}

	uc.logger.Info("Power scheduling completed successfully", "mode", mode)

	return &SchedulePowerOutput{
		Schedule: schedule,
		Message:  fmt.Sprintf("Power scheduling configured successfully (%s mode)", mode),
		TestMode: input.TestMode,
	}, nil
}