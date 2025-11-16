// internal/infrastructure/rtc/linux_rtc.go
package rtc

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"rtc-scheduler/internal/domain/repositories"
)

const (
	wakeAlarmPath = "/sys/class/rtc/rtc0/wakealarm"
	timePath      = "/sys/class/rtc/rtc0/since_epoch"
)

type LinuxRTC struct {
	wakeAlarmPath string
	timePath      string
}

// Verificar que implementa la interfaz
var _ repositories.RTCRepository = (*LinuxRTC)(nil)

func NewLinuxRTC() *LinuxRTC {
	return &LinuxRTC{
		wakeAlarmPath: wakeAlarmPath,
		timePath:      timePath,
	}
}

func (r *LinuxRTC) SetWakeAlarm(t time.Time) error {
	// Limpiar alarma anterior
	if err := r.ClearWakeAlarm(); err != nil {
		return fmt.Errorf("failed to clear previous alarm: %w", err)
	}

	// Configurar nueva alarma
	timestamp := strconv.FormatInt(t.Unix(), 10)
	if err := os.WriteFile(r.wakeAlarmPath, []byte(timestamp), 0644); err != nil {
		return fmt.Errorf("failed to set wake alarm: %w", err)
	}

	return nil
}

func (r *LinuxRTC) GetWakeAlarm() (time.Time, error) {
	data, err := os.ReadFile(r.wakeAlarmPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read wake alarm: %w", err)
	}

	alarmStr := strings.TrimSpace(string(data))
	if alarmStr == "" || alarmStr == "0" {
		return time.Time{}, fmt.Errorf("no wake alarm set")
	}

	timestamp, err := strconv.ParseInt(alarmStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid wake alarm timestamp: %w", err)
	}

	return time.Unix(timestamp, 0), nil
}

func (r *LinuxRTC) ClearWakeAlarm() error {
	if err := os.WriteFile(r.wakeAlarmPath, []byte("0"), 0644); err != nil {
		return fmt.Errorf("failed to clear wake alarm: %w", err)
	}
	return nil
}

func (r *LinuxRTC) GetCurrentTime() (time.Time, error) {
	data, err := os.ReadFile(r.timePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read RTC time: %w", err)
	}

	timestamp, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RTC timestamp: %w", err)
	}

	return time.Unix(timestamp, 0), nil
}

func (r *LinuxRTC) IsAvailable() bool {
	// Verificar que los archivos del dispositivo existen
	if _, err := os.Stat(r.wakeAlarmPath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(r.timePath); os.IsNotExist(err) {
		return false
	}
	return true
}