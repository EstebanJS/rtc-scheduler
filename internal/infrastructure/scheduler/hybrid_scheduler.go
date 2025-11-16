// internal/infrastructure/scheduler/hybrid_scheduler.go
package scheduler

import (
	"fmt"
	"time"

	"rtc-scheduler/internal/domain/repositories"
)

// HybridScheduler combina AtScheduler y SystemdTimerScheduler para máxima compatibilidad
type HybridScheduler struct {
	atScheduler    *AtScheduler
	timerScheduler *SystemdTimerScheduler
	testMode       bool
}

// Verificar que implementa la interfaz
var _ repositories.SchedulerRepository = (*HybridScheduler)(nil)

// NewHybridScheduler crea una nueva instancia del scheduler híbrido
func NewHybridScheduler() *HybridScheduler {
	return &HybridScheduler{
		atScheduler:    NewAtScheduler(),
		timerScheduler: NewSystemdTimerScheduler(),
		testMode:       false,
	}
}

// NewHybridSchedulerWithTestMode crea una instancia en modo prueba
func NewHybridSchedulerWithTestMode(testMode bool) *HybridScheduler {
	return &HybridScheduler{
		atScheduler:    NewAtSchedulerWithTestMode(testMode),
		timerScheduler: NewSystemdTimerSchedulerWithTestMode(testMode),
		testMode:       testMode,
	}
}

// ScheduleShutdown elige el mejor scheduler disponible
func (s *HybridScheduler) ScheduleShutdown(t time.Time) error {
	// Prioridad 1: AtScheduler (si está disponible y filesystem es writable)
	if s.atScheduler.IsAvailable() && s.atScheduler.isFilesystemWritable() {
		return s.atScheduler.ScheduleShutdown(t)
	}

	// Prioridad 2: SystemdTimerScheduler (si está disponible)
	if s.timerScheduler.IsAvailable() {
		return s.timerScheduler.ScheduleShutdown(t)
	}

	// Ningún scheduler disponible
	return fmt.Errorf("no suitable scheduler available: at command not available or filesystem read-only, and systemd-run not available")
}

// CancelShutdown intenta cancelar en todos los schedulers disponibles
func (s *HybridScheduler) CancelShutdown() error {
	var lastErr error

	// Intentar cancelar en AtScheduler
	if s.atScheduler.IsAvailable() {
		if err := s.atScheduler.CancelShutdown(); err != nil {
			lastErr = err
		}
	}

	// Intentar cancelar en SystemdTimerScheduler
	if s.timerScheduler.IsAvailable() {
		if err := s.timerScheduler.CancelShutdown(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// ListScheduledJobs combina trabajos de todos los schedulers
func (s *HybridScheduler) ListScheduledJobs() ([]*repositories.ShutdownJob, error) {
	var allJobs []*repositories.ShutdownJob

	// Obtener trabajos de AtScheduler
	if s.atScheduler.IsAvailable() {
		if jobs, err := s.atScheduler.ListScheduledJobs(); err == nil {
			allJobs = append(allJobs, jobs...)
		}
	}

	// Obtener trabajos de SystemdTimerScheduler
	if s.timerScheduler.IsAvailable() {
		if jobs, err := s.timerScheduler.ListScheduledJobs(); err == nil {
			allJobs = append(allJobs, jobs...)
		}
	}

	return allJobs, nil
}

// IsAvailable verifica si al menos un scheduler está disponible
func (s *HybridScheduler) IsAvailable() bool {
	// Está disponible si:
	// 1. AtScheduler está disponible Y filesystem es writable, O
	// 2. SystemdTimerScheduler está disponible
	return (s.atScheduler.IsAvailable() && s.atScheduler.isFilesystemWritable()) || s.timerScheduler.IsAvailable()
}

// GetSchedulerStatus retorna información detallada sobre el estado de los schedulers
func (s *HybridScheduler) GetSchedulerStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// Estado del AtScheduler
	atAvailable := s.atScheduler.IsAvailable()
	atWritable := s.atScheduler.isFilesystemWritable()
	status["at_scheduler"] = map[string]interface{}{
		"available":     atAvailable,
		"filesystem_writable": atWritable,
		"usable":        atAvailable && atWritable,
	}

	// Estado del SystemdTimerScheduler
	timerAvailable := s.timerScheduler.IsAvailable()
	status["systemd_timer_scheduler"] = map[string]interface{}{
		"available": timerAvailable,
		"usable":    timerAvailable,
	}

	// Scheduler activo
	var activeScheduler string
	if atAvailable && atWritable {
		activeScheduler = "at_scheduler"
	} else if timerAvailable {
		activeScheduler = "systemd_timer_scheduler"
	} else {
		activeScheduler = "none"
	}
	status["active_scheduler"] = activeScheduler

	return status
}

// GetJobDetails obtiene detalles de un trabajo específico
func (s *HybridScheduler) GetJobDetails(jobID string) (string, error) {
	// Intentar primero en AtScheduler
	if s.atScheduler.IsAvailable() {
		if details, err := s.atScheduler.GetJobDetails(jobID); err == nil {
			return details, nil
		}
	}

	// Intentar en SystemdTimerScheduler
	if s.timerScheduler.IsAvailable() {
		if details, err := s.timerScheduler.GetJobDetails(jobID); err == nil {
			return details, nil
		}
	}

	return "", fmt.Errorf("job details not found for ID: %s", jobID)
}

// ScheduleAt programa un comando usando el mejor scheduler disponible
func (s *HybridScheduler) ScheduleAt(t time.Time, command string) error {
	// Prioridad 1: AtScheduler
	if s.atScheduler.IsAvailable() && s.atScheduler.isFilesystemWritable() {
		return s.atScheduler.ScheduleAt(t, command)
	}

	// Prioridad 2: SystemdTimerScheduler
	if s.timerScheduler.IsAvailable() {
		return s.timerScheduler.ScheduleAt(t, command)
	}

	return fmt.Errorf("no suitable scheduler available for custom command scheduling")
}

// ParseJobID parsea el ID del trabajo de la salida del scheduler activo
func (s *HybridScheduler) ParseJobID(output string) (string, error) {
	// Intentar con AtScheduler primero
	if s.atScheduler.IsAvailable() && s.atScheduler.isFilesystemWritable() {
		return s.atScheduler.ParseJobID(output)
	}

	// Intentar con SystemdTimerScheduler
	if s.timerScheduler.IsAvailable() {
		return s.timerScheduler.ParseJobID(output)
	}

	return "", fmt.Errorf("no suitable scheduler available for parsing job ID")
}

// CountScheduledJobs cuenta trabajos de todos los schedulers
func (s *HybridScheduler) CountScheduledJobs() (int, error) {
	jobs, err := s.ListScheduledJobs()
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}