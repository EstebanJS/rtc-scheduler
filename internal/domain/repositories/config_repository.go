package repositories

import "rtc-scheduler/internal/domain/entities"

type ConfigRepository interface {
	Load() (*entities.Config, error)
	Save(*entities.Config) error
	Delete() error
	Exists() bool
	CreateDefault() error
}