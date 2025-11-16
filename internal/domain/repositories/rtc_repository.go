// Interfaces (contratos)

package repositories

import "time"

type RTCRepository interface {
	SetWakeAlarm(t time.Time) error
	GetWakeAlarm() (time.Time, error)
	ClearWakeAlarm() error
	GetCurrentTime() (time.Time, error)
	IsAvailable() bool
}