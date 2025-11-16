// internal/infrastructure/config/json_config.go
package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"rtc-scheduler/internal/domain/entities"
	"rtc-scheduler/internal/domain/repositories"
)

var (
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrInvalidConfig  = errors.New("invalid configuration format")
)

// configDTO es la estructura para serialización JSON
type configDTO struct {
	WakeTime     string `json:"wake_time"`
	ShutdownTime string `json:"shutdown_time"`
	Enabled      bool   `json:"enabled"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// JSONConfigRepository implementa ConfigRepository usando archivos JSON
type JSONConfigRepository struct {
	filePath string
}

// Verificar que implementa la interfaz
var _ repositories.ConfigRepository = (*JSONConfigRepository)(nil)

// NewJSONConfigRepository crea una nueva instancia
func NewJSONConfigRepository(filePath string) *JSONConfigRepository {
	return &JSONConfigRepository{
		filePath: filePath,
	}
}

// Save guarda la configuración en un archivo JSON
func (r *JSONConfigRepository) Save(config *entities.Config) error {
	// Convertir a DTO
	dto := &configDTO{
		WakeTime:     config.WakeTime,
		ShutdownTime: config.ShutdownTime,
		Enabled:      config.Enabled,
		CreatedAt:    config.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    config.UpdatedAt.Format(time.RFC3339),
	}

	// Serializar a JSON con formato legible
	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return err
	}

	// Escribir archivo
	return os.WriteFile(r.filePath, data, 0644)
}

// Load carga la configuración desde el archivo JSON
func (r *JSONConfigRepository) Load() (*entities.Config, error) {
	// Verificar que el archivo existe
	if !r.Exists() {
		return nil, ErrConfigNotFound
	}

	// Leer archivo
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	// Deserializar JSON
	var dto configDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, ErrInvalidConfig
	}

	// Parsear timestamps
	createdAt, err := time.Parse(time.RFC3339, dto.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, dto.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Convertir a entidad
	config := &entities.Config{
		WakeTime:     dto.WakeTime,
		ShutdownTime: dto.ShutdownTime,
		Enabled:      dto.Enabled,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	return config, nil
}

// Delete elimina el archivo de configuración
func (r *JSONConfigRepository) Delete() error {
	if !r.Exists() {
		return nil // Ya está eliminado
	}
	return os.Remove(r.filePath)
}

// Exists verifica si el archivo de configuración existe
func (r *JSONConfigRepository) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !os.IsNotExist(err)
}

// GetFilePath retorna la ruta del archivo (útil para debugging)
func (r *JSONConfigRepository) GetFilePath() string {
	return r.filePath
}

// CreateDefault crea una configuración por defecto
func (r *JSONConfigRepository) CreateDefault() error {
	defaultConfig, err := entities.NewConfig("08:00", "22:00", false)
	if err != nil {
		return err
	}
	return r.Save(defaultConfig)
}