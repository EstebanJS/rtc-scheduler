// internal/application/usecases/enable_service.go
package usecases

import (
	"errors"
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

var (
	ErrServiceNotInstalled = errors.New("service is not installed")
)

type EnableServiceInput struct{}

type EnableServiceOutput struct {
	ServiceEnabled bool
	Message        string
}

// EnableServiceUseCase maneja la habilitación del servicio
type EnableServiceUseCase struct {
	configRepo repositories.ConfigRepository
	serviceRepo repositories.ServiceRepository
	logger      logger.Logger
}

func NewEnableServiceUseCase(
	config repositories.ConfigRepository,
	service repositories.ServiceRepository,
	log logger.Logger,
) *EnableServiceUseCase {
	return &EnableServiceUseCase{
		configRepo: config,
		serviceRepo: service,
		logger:     log,
	}
}

func (uc *EnableServiceUseCase) Execute(input *EnableServiceInput) (*EnableServiceOutput, error) {
	uc.logger.Info("Enabling service")

	// Verificar que el servicio esté instalado
	if !uc.serviceRepo.IsInstalled() {
		return nil, ErrServiceNotInstalled
	}

	// Habilitar el servicio
	if err := uc.serviceRepo.Enable(); err != nil {
		uc.logger.Error("Failed to enable service", "error", err)
		return nil, err
	}

	// Actualizar configuración
	if uc.configRepo.Exists() {
		config, err := uc.configRepo.Load()
		if err == nil {
			config.Enabled = true
			if err := uc.configRepo.Save(config); err != nil {
				uc.logger.Warn("Failed to update config", "error", err)
			}
		}
	}

	uc.logger.Info("Service enabled successfully")

	return &EnableServiceOutput{
		ServiceEnabled: true,
		Message:        "Service enabled successfully",
	}, nil
}