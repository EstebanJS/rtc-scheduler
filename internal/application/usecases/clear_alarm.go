// internal/application/usecases/clear_alarm.go
package usecases

import (
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

type ClearAlarmInput struct{}

type ClearAlarmOutput struct {
	AlarmCleared bool
	Message      string
}

// ClearAlarmUseCase maneja la limpieza de alarmas RTC
type ClearAlarmUseCase struct {
	rtcRepo       repositories.RTCRepository
	schedulerRepo repositories.SchedulerRepository
	logger        logger.Logger
}

func NewClearAlarmUseCase(
	rtc repositories.RTCRepository,
	scheduler repositories.SchedulerRepository,
	log logger.Logger,
) *ClearAlarmUseCase {
	return &ClearAlarmUseCase{
		rtcRepo:       rtc,
		schedulerRepo: scheduler,
		logger:        log,
	}
}

func (uc *ClearAlarmUseCase) Execute(input *ClearAlarmInput) (*ClearAlarmOutput, error) {
	uc.logger.Info("Clearing wake alarm")

	// Limpiar alarma RTC
	if err := uc.rtcRepo.ClearWakeAlarm(); err != nil {
		uc.logger.Error("Failed to clear RTC wake alarm", "error", err)
		return nil, err
	}

	// Cancelar tareas de apagado programadas
	if err := uc.schedulerRepo.CancelShutdown(); err != nil {
		uc.logger.Warn("Failed to cancel scheduled shutdowns", "error", err)
	}

	uc.logger.Info("Wake alarm cleared successfully")

	return &ClearAlarmOutput{
		AlarmCleared: true,
		Message:      "Wake alarm cleared successfully",
	}, nil
}