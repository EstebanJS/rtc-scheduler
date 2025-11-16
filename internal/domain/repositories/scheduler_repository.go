package repositories

import "time"

type SchedulerRepository interface {
	ScheduleShutdown(t time.Time) error
	CancelShutdown() error
	ListScheduledJobs() ([]*ShutdownJob, error)
	IsAvailable() bool
}

type ShutdownJob struct {
	ID          string
	ScheduledAt time.Time
	Command     string
}