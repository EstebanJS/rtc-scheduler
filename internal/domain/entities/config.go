// internal/domain/entities/config.go
package entities

import (
	"errors"
	"time"
)

var (
	ErrEmptyWakeTime     = errors.New("wake time cannot be empty")
	ErrEmptyShutdownTime = errors.New("shutdown time cannot be empty")
	ErrInvalidTimeFormat = errors.New("invalid time format, use HH:MM")
)

// Config representa la configuración del sistema
type Config struct {
	WakeTime     string
	ShutdownTime string
	Enabled      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewConfig crea una nueva configuración con validación
func NewConfig(wakeTime, shutdownTime string, enabled bool) (*Config, error) {
	config := &Config{
		WakeTime:     wakeTime,
		ShutdownTime: shutdownTime,
		Enabled:      enabled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate verifica que la configuración sea válida
func (c *Config) Validate() error {
	if c.WakeTime == "" {
		return ErrEmptyWakeTime
	}

	if c.ShutdownTime == "" {
		return ErrEmptyShutdownTime
	}

	// Validar formato HH:MM
	if !isValidTimeFormat(c.WakeTime) {
		return ErrInvalidTimeFormat
	}

	if !isValidTimeFormat(c.ShutdownTime) {
		return ErrInvalidTimeFormat
	}

	return nil
}

// Update actualiza el timestamp de modificación
func (c *Config) Update() {
	c.UpdatedAt = time.Now()
}

// Disable deshabilita la configuración
func (c *Config) Disable() {
	c.Enabled = false
	c.Update()
}

// Enable habilita la configuración
func (c *Config) Enable() {
	c.Enabled = true
	c.Update()
}

// isValidTimeFormat verifica si el formato de hora es válido (HH:MM)
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// ParseToSchedule convierte la configuración en un Schedule
func (c *Config) ParseToSchedule() (*Schedule, error) {
	now := time.Now()

	// Parsear wake time
	wakeTime, err := time.ParseInLocation("15:04", c.WakeTime, time.Local)
	if err != nil {
		return nil, err
	}
	wakeTime = time.Date(now.Year(), now.Month(), now.Day(),
		wakeTime.Hour(), wakeTime.Minute(), 0, 0, time.Local)
	if wakeTime.Before(now) {
		wakeTime = wakeTime.Add(24 * time.Hour)
	}

	// Parsear shutdown time
	shutdownTime, err := time.ParseInLocation("15:04", c.ShutdownTime, time.Local)
	if err != nil {
		return nil, err
	}
	shutdownTime = time.Date(now.Year(), now.Month(), now.Day(),
		shutdownTime.Hour(), shutdownTime.Minute(), 0, 0, time.Local)
	if shutdownTime.Before(now) {
		shutdownTime = shutdownTime.Add(24 * time.Hour)
	}

	// Ajustar si shutdown es antes de wake
	if shutdownTime.Before(wakeTime) {
		shutdownTime = shutdownTime.Add(24 * time.Hour)
	}

	return NewSchedule(wakeTime, shutdownTime)
}