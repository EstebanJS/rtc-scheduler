// internal/application/usecases/uninstall_service.go
package usecases

import (
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

type UninstallServiceInput struct{}

type UninstallServiceOutput struct {
	ServiceUninstalled bool
	ConfigDeleted      bool
	AlarmsCleared      bool
	Message            string
}

// UninstallServiceUseCase maneja la desinstalación del servicio
type UninstallServiceUseCase struct {
	rtcRepo       repositories.RTCRepository
	configRepo    repositories.ConfigRepository
	serviceRepo   repositories.ServiceRepository
	schedulerRepo repositories.SchedulerRepository
	logger        logger.Logger
}

func NewUninstallServiceUseCase(
	rtc repositories.RTCRepository,
	config repositories.ConfigRepository,
	service repositories.ServiceRepository,
	scheduler repositories.SchedulerRepository,
	log logger.Logger,
) *UninstallServiceUseCase {
	return &UninstallServiceUseCase{
		rtcRepo:       rtc,
		configRepo:    config,
		serviceRepo:   service,
		schedulerRepo: scheduler,
		logger:        log,
	}
}

func (uc *UninstallServiceUseCase) Execute(input *UninstallServiceInput) (*UninstallServiceOutput, error) {
	uc.logger.Info("Starting service uninstallation")

	output := &UninstallServiceOutput{}

	// 1. Detener y deshabilitar servicio
	if uc.serviceRepo.IsInstalled() {
		if err := uc.serviceRepo.Stop(); err != nil {
			uc.logger.Warn("Failed to stop service", "error", err)
		}

		if err := uc.serviceRepo.Disable(); err != nil {
			uc.logger.Warn("Failed to disable service", "error", err)
		}

		if err := uc.serviceRepo.Uninstall(); err != nil {
			uc.logger.Error("Failed to uninstall service", "error", err)
			return nil, err
		}
		output.ServiceUninstalled = true
		uc.logger.Info("Service uninstalled successfully")
	} else {
		uc.logger.Info("Service was not installed")
		output.ServiceUninstalled = true // Considerado exitoso
	}

	// 2. Limpiar alarmas RTC
	if err := uc.rtcRepo.ClearWakeAlarm(); err != nil {
		uc.logger.Warn("Failed to clear RTC wake alarm", "error", err)
	} else {
		output.AlarmsCleared = true
		uc.logger.Info("RTC wake alarm cleared")
	}

	// 3. Cancelar tareas de apagado programadas
	if err := uc.schedulerRepo.CancelShutdown(); err != nil {
		uc.logger.Warn("Failed to cancel scheduled shutdowns", "error", err)
	} else {
		uc.logger.Info("Scheduled shutdowns cancelled")
	}

	// 4. Eliminar configuración
	if uc.configRepo.Exists() {
		if err := uc.configRepo.Delete(); err != nil {
			uc.logger.Warn("Failed to delete configuration", "error", err)
		} else {
			output.ConfigDeleted = true
			uc.logger.Info("Configuration deleted")
		}
	} else {
		output.ConfigDeleted = true // Ya estaba eliminado
	}

	output.Message = "Service uninstallation completed"
	uc.logger.Info("Service uninstallation completed successfully")

	return output, nil
}