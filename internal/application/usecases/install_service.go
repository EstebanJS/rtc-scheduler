// internal/application/usecases/install_service.go
package usecases

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"rtc-scheduler/internal/domain/entities"
	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

var (
	ErrServiceAlreadyInstalled = errors.New("service is already installed")
	ErrInvalidExecutablePath   = errors.New("invalid executable path")
)

// InstallServiceInput representa los datos de entrada
type InstallServiceInput struct {
	WakeTime     string
	ShutdownTime string
	ExecutablePath string
}

// InstallServiceOutput representa el resultado
type InstallServiceOutput struct {
	ServiceInstalled bool
	ConfigCreated    bool
	Message          string
}

// InstallServiceUseCase maneja la instalación del servicio
type InstallServiceUseCase struct {
	configRepo  repositories.ConfigRepository
	serviceRepo repositories.ServiceRepository
	logger      logger.Logger
}

// NewInstallServiceUseCase crea una nueva instancia
func NewInstallServiceUseCase(
	config repositories.ConfigRepository,
	service repositories.ServiceRepository,
	log logger.Logger,
) *InstallServiceUseCase {
	return &InstallServiceUseCase{
		configRepo:  config,
		serviceRepo: service,
		logger:      log,
	}
}

// Execute ejecuta la instalación del servicio
func (uc *InstallServiceUseCase) Execute(input *InstallServiceInput) (*InstallServiceOutput, error) {
	uc.logger.Info("Starting service installation",
		"wake_time", input.WakeTime,
		"shutdown_time", input.ShutdownTime,
	)

	// 1. Verificar que no esté ya instalado
	if uc.serviceRepo.IsInstalled() {
		uc.logger.Warn("Service is already installed")
		return nil, ErrServiceAlreadyInstalled
	}

	// 2. Validar y crear configuración
	if err := uc.createConfiguration(input); err != nil {
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	// 3. Obtener ruta del ejecutable actual
	execPath, err := uc.getExecutablePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// 4. Instalar servicio systemd
	if err := uc.serviceRepo.Install(execPath); err != nil {
		// Limpiar configuración si falla la instalación
		uc.configRepo.Delete()
		return nil, fmt.Errorf("failed to install service: %w", err)
	}

	// 5. Habilitar servicio
	if err := uc.serviceRepo.Enable(); err != nil {
		uc.logger.Warn("Failed to enable service", "error", err)
		// No fallar aquí, el servicio está instalado
	}

	// 6. Iniciar servicio (esto programará el primer ciclo)
	if err := uc.serviceRepo.Start(); err != nil {
		uc.logger.Warn("Failed to start service", "error", err)
		// No fallar aquí, el servicio está instalado
	}

	uc.logger.Info("Service installed successfully")

	return &InstallServiceOutput{
		ServiceInstalled: true,
		ConfigCreated:    true,
		Message:          "Service installed and enabled successfully",
	}, nil
}

// createConfiguration crea y guarda la configuración
func (uc *InstallServiceUseCase) createConfiguration(input *InstallServiceInput) error {
	// Validar formato de horarios
	config, err := entities.NewConfig(input.WakeTime, input.ShutdownTime, true)
	if err != nil {
		uc.logger.Error("Invalid configuration", "error", err)
		return err
	}

	// Guardar configuración directamente
	if err := uc.configRepo.Save(config); err != nil {
		uc.logger.Error("Failed to save configuration", "error", err)
		return err
	}

	uc.logger.Info("Configuration created successfully")
	return nil
}

// getExecutablePath obtiene la ruta del ejecutable actual
func (uc *InstallServiceUseCase) getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Resolver enlaces simbólicos
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}

	// Convertir a ruta absoluta
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	// Verificar que el archivo existe
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return "", ErrInvalidExecutablePath
	}

	uc.logger.Debug("Executable path resolved", "path", execPath)
	return execPath, nil
}

// Rollback deshace la instalación en caso de error
func (uc *InstallServiceUseCase) Rollback() error {
	uc.logger.Info("Rolling back installation")

	// Eliminar configuración
	_ = uc.configRepo.Delete() //nolint:errcheck

	// Desinstalar servicio
	if err := uc.serviceRepo.Uninstall(); err != nil {
		uc.logger.Warn("Failed to uninstall service during rollback", "error", err)
	}

	return nil
}