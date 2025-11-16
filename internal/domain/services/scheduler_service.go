// internal/domain/services/scheduler_service.go
package services

import (
	"fmt"
	"time"

	"rtc-scheduler/internal/domain/repositories"
)

type SchedulerService struct {
	rtcRepo repositories.RTCRepository
}

func NewSchedulerService(rtcRepo repositories.RTCRepository) *SchedulerService {
	return &SchedulerService{rtcRepo: rtcRepo}
}

func (s *SchedulerService) SchedulePower(t time.Time, action string) error {
	if action != "on" {
		return fmt.Errorf("unsupported action: %s (only 'on' is supported)", action)
	}

	if err := s.rtcRepo.SetWakeAlarm(t); err != nil {
		return fmt.Errorf("failed to set wake alarm: %w", err)
	}

	return nil
}

func (s *SchedulerService) ClearPowerSchedule() error {
	if err := s.rtcRepo.ClearWakeAlarm(); err != nil {
		return fmt.Errorf("failed to clear wake alarm: %w", err)
	}
	return nil
}

func (s *SchedulerService) GetPowerSchedule() (time.Time, error) {
	return s.rtcRepo.GetWakeAlarm()
}