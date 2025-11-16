// internal/domain/entities/schedule.go
package entities

import (
	"errors"
	"time"
)

var (
	ErrInvalidWakeTime     = errors.New("wake time must be in the future")
	ErrInvalidShutdownTime = errors.New("shutdown time must be after wake time")
)

type Schedule struct {
	WakeTime     time.Time
	ShutdownTime time.Time
	Enabled      bool
	CreatedAt    time.Time
}

// NewSchedule crea un nuevo schedule con validación
func NewSchedule(wakeTime, shutdownTime time.Time) (*Schedule, error) {
	schedule := &Schedule{
		WakeTime:     wakeTime,
		ShutdownTime: shutdownTime,
		Enabled:      true,
		CreatedAt:    time.Now(),
	}

	if err := schedule.Validate(); err != nil {
		return nil, err
	}

	return schedule, nil
}

// Validate verifica que el schedule sea válido
func (s *Schedule) Validate() error {
	now := time.Now()

	if s.WakeTime.Before(now) {
		return ErrInvalidWakeTime
	}

	if s.ShutdownTime.Before(s.WakeTime) {
		return ErrInvalidShutdownTime
	}

	return nil
}

// IsActive verifica si el schedule está activo
func (s *Schedule) IsActive() bool {
	return s.Enabled
}

// NextWakeTime retorna la próxima hora de encendido
func (s *Schedule) NextWakeTime() time.Time {
	return s.WakeTime
}

// NextShutdownTime retorna la próxima hora de apagado
func (s *Schedule) NextShutdownTime() time.Time {
	return s.ShutdownTime
}