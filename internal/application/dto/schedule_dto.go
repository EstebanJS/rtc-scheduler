// Data Transfer Objects

package dto

import "time"

type ScheduleDTO struct {
	ID      string
	Time    time.Time
	Action  string
	Enabled bool
}