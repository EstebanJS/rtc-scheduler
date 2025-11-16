// internal/application/usecases/disable_service.go
package usecases

import (
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

type DisableServiceInput struct{}

type DisableServiceOutput struct {
	ServiceDisabled bool
	Message         string
}

// DisableServiceUseCase maneja la deshabilitación del servicio
type DisableServiceUseCase struct {
	configRepo    repositories.ConfigRepository
	serviceRepo   repositories.ServiceRepository
	schedulerRepo repositories.SchedulerRepository
	rtcRepo       repositories.RTCRepository
	logger        logger.Logger
}

func NewDisableServiceUseCase(
	config repositories.ConfigRepository,
	service repositories.ServiceRepository,
	scheduler repositories.SchedulerRepository,
	rtc repositories.RTCRepository,
	log logger.Logger,
) *DisableServiceUseCase {
	return &DisableServiceUseCase{
		configRepo:    config,
		serviceRepo:   service,
		schedulerRepo: scheduler,
		rtcRepo:       rtc,
		logger:        log,
	}
}

func (uc *DisableServiceUseCase) Execute(input *DisableServiceInput) (*DisableServiceOutput, error) {
	uc.logger.Info("Disabling service")

	// Verificar que el servicio esté instalado
	if !uc.serviceRepo.IsInstalled() {
		return nil, ErrServiceNotInstalled
	}

	// Deshabilitar el servicio
	if err := uc.serviceRepo.Disable(); err != nil {
		uc.logger.Error("Failed to disable service", "error", err)
		return nil, err
	}

	// Limpiar alarmas RTC
	if err := uc.rtcRepo.ClearWakeAlarm(); err != nil {
		uc.logger.Warn("Failed to clear RTC wake alarm", "error", err)
	}

	// Cancelar tareas programadas
	if err := uc.schedulerRepo.CancelShutdown(); err != nil {
		uc.logger.Warn("Failed to cancel scheduled shutdowns", "error", err)
	}

	// Actualizar configuración
	if uc.configRepo.Exists() {
		config, err := uc.configRepo.Load()
		if err == nil {
			config.Enabled = false
			if err := uc.configRepo.Save(config); err != nil {
				uc.logger.Warn("Failed to update config", "error", err)
			}
		}
	}

	uc.logger.Info("Service disabled successfully")

	return &DisableServiceOutput{
		ServiceDisabled: true,
		Message:         "Service disabled successfully",
	}, nil
}