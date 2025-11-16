// internal/infrastructure/scheduler/systemd_timer_scheduler.go
package scheduler

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"rtc-scheduler/internal/domain/repositories"
)

var (
	ErrSystemdRunNotAvailable = errors.New("systemd-run command is not available")
	ErrSystemdTimerFailed     = errors.New("failed to create systemd timer")
)

// SystemdTimerScheduler implementa SchedulerRepository usando systemd-run
type SystemdTimerScheduler struct {
	testMode bool
}

// Verificar que implementa la interfaz
var _ repositories.SchedulerRepository = (*SystemdTimerScheduler)(nil)

// NewSystemdTimerScheduler crea una nueva instancia
func NewSystemdTimerScheduler() *SystemdTimerScheduler {
	return &SystemdTimerScheduler{
		testMode: false,
	}
}

// NewSystemdTimerSchedulerWithTestMode crea una instancia en modo prueba
func NewSystemdTimerSchedulerWithTestMode(testMode bool) *SystemdTimerScheduler {
	return &SystemdTimerScheduler{
		testMode: testMode,
	}
}

// ScheduleShutdown programa un apagado del sistema usando systemd-run
func (s *SystemdTimerScheduler) ScheduleShutdown(t time.Time) error {
	if !s.IsAvailable() {
		return ErrSystemdRunNotAvailable
	}

	// Calcular minutos hasta el apagado
	duration := time.Until(t)
	if duration < 0 {
		return ErrInvalidTime
	}

	minutes := int(duration.Minutes())
	if minutes < 1 {
		minutes = 1 // Mínimo 1 minuto
	}

	// Preparar comando
	var command string
	if s.testMode {
		command = "/usr/bin/wall 'TEST MODE: Shutdown time reached'"
	} else {
		// Usar systemctl poweroff con ruta absoluta correcta
		command = "/usr/bin/systemctl poweroff"
	}

	// Crear timer con systemd-run
	// --on-calendar: ejecutar en un tiempo específico
	// --timer-property: propiedades del timer
	// --unit: nombre único para el timer
	timerName := fmt.Sprintf("rtc-scheduler-shutdown-%d", time.Now().Unix())

	args := []string{
		"--on-active", fmt.Sprintf("%dm", minutes), // Ejecutar X minutos después de activarse
		"--timer-property", "AccuracySec=1s",       // Alta precisión
		"--setenv", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", // Asegurar PATH correcto
		"--unit", timerName,
		"--description", "RTC Scheduler Shutdown Timer",
		"--service-type", "oneshot",
		"sh", "-c", command, // Usar sh -c para ejecutar el comando correctamente
	}

	cmd := exec.Command("systemd-run", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create systemd timer: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// CancelShutdown cancela todos los timers de apagado programados
func (s *SystemdTimerScheduler) CancelShutdown() error {
	if !s.IsAvailable() {
		return nil // Si systemd-run no está disponible, no hay nada que cancelar
	}

	// Listar timers activos
	timers, err := s.listActiveTimers()
	if err != nil {
		return err
	}

	// Cancelar cada timer que sea de rtc-scheduler
	for _, timer := range timers {
		if s.isShutdownTimer(timer) {
			if err := s.cancelTimer(timer); err != nil {
				// Continuar aunque falle cancelar uno
				continue
			}
		}
	}

	return nil
}

// ListScheduledJobs lista todas las tareas programadas (timers activos)
func (s *SystemdTimerScheduler) ListScheduledJobs() ([]*repositories.ShutdownJob, error) {
	if !s.IsAvailable() {
		return nil, ErrSystemdRunNotAvailable
	}

	timers, err := s.listActiveTimers()
	if err != nil {
		return nil, err
	}

	var jobs []*repositories.ShutdownJob
	for _, timer := range timers {
		if s.isShutdownTimer(timer) {
			job := s.parseTimerToJob(timer)
			if job != nil {
				jobs = append(jobs, job)
			}
		}
	}

	return jobs, nil
}

// IsAvailable verifica si systemd-run está disponible
func (s *SystemdTimerScheduler) IsAvailable() bool {
	cmd := exec.Command("which", "systemd-run")
	err := cmd.Run()
	return err == nil
}

// listActiveTimers lista todos los timers activos de systemd
func (s *SystemdTimerScheduler) listActiveTimers() ([]string, error) {
	cmd := exec.Command("systemctl", "list-timers", "--all", "--no-pager", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list timers: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var timers []string

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Extraer el nombre del timer de la línea
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			timerName := fields[5]
			timers = append(timers, timerName)
		}
	}

	return timers, nil
}

// isShutdownTimer verifica si un timer es de rtc-scheduler
func (s *SystemdTimerScheduler) isShutdownTimer(timerName string) bool {
	return strings.HasPrefix(timerName, "rtc-scheduler-shutdown-")
}

// cancelTimer cancela un timer específico
func (s *SystemdTimerScheduler) cancelTimer(timerName string) error {
	cmd := exec.Command("systemctl", "stop", timerName)
	return cmd.Run()
}

// parseTimerToJob convierte un timer en un ShutdownJob
func (s *SystemdTimerScheduler) parseTimerToJob(timerName string) *repositories.ShutdownJob {
	// Extraer timestamp del nombre del timer
	re := regexp.MustCompile(`rtc-scheduler-shutdown-(\d+)`)
	matches := re.FindStringSubmatch(timerName)

	if len(matches) < 2 {
		return nil
	}

	timestampStr := matches[1]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil
	}

	// Obtener detalles del timer
	cmd := exec.Command("systemctl", "show", timerName, "--property", "TimersCalendar")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	// Parsear la salida para obtener la hora programada
	// Formato: TimersCalendar=OnCalendar=2025-11-16 22:00:00 UTC; ...
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "OnActive=") {
			// Extraer la duración y calcular la hora
			re := regexp.MustCompile(`OnActive=(\d+)m`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				minutes, err := strconv.Atoi(matches[1])
				if err == nil {
					scheduledTime := time.Unix(timestamp, 0).Add(time.Duration(minutes) * time.Minute)
					return &repositories.ShutdownJob{
						ID:          timerName,
						ScheduledAt: scheduledTime,
						Command:     "shutdown",
					}
				}
			}
		}
	}

	return nil
}

// GetJobDetails obtiene detalles completos de un timer
func (s *SystemdTimerScheduler) GetJobDetails(timerName string) (string, error) {
	cmd := exec.Command("systemctl", "show", timerName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get timer details: %w", err)
	}
	return string(output), nil
}

// ScheduleAt programa un comando a una hora específica usando systemd-run
func (s *SystemdTimerScheduler) ScheduleAt(t time.Time, command string) error {
	if !s.IsAvailable() {
		return ErrSystemdRunNotAvailable
	}

	duration := time.Until(t)
	if duration < 0 {
		return ErrInvalidTime
	}

	minutes := int(duration.Minutes())
	if minutes < 1 {
		minutes = 1
	}

	timerName := fmt.Sprintf("rtc-scheduler-custom-%d", time.Now().Unix())

	args := []string{
		"--on-active", fmt.Sprintf("%dm", minutes),
		"--timer-property", "AccuracySec=1s",
		"--setenv", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", // Asegurar PATH correcto
		"--unit", timerName,
		"--description", "RTC Scheduler Custom Timer",
		"--service-type", "oneshot",
		"sh", "-c", command, // Usar sh -c para ejecutar el comando correctamente
	}

	cmd := exec.Command("systemd-run", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create systemd timer: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ParseJobID extrae el ID del timer de la salida de systemd-run
func (s *SystemdTimerScheduler) ParseJobID(output string) (string, error) {
	// Buscar patrón en la salida
	re := regexp.MustCompile(`Running timer as unit: (.+)\.timer`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 2 {
		return "", errors.New("could not parse timer unit from output")
	}

	return strings.TrimSpace(matches[1]), nil
}

// CountScheduledJobs cuenta cuántos timers están programados
func (s *SystemdTimerScheduler) CountScheduledJobs() (int, error) {
	jobs, err := s.ListScheduledJobs()
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}