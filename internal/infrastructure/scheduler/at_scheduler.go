// internal/infrastructure/scheduler/at_scheduler.go
package scheduler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"rtc-scheduler/internal/domain/repositories"
)

var (
	ErrAtNotAvailable     = errors.New("'at' command is not available")
	ErrInvalidTime        = errors.New("invalid time for scheduling")
	ErrFilesystemReadOnly = errors.New("filesystem is read-only, cannot schedule shutdown")
)

// AtScheduler implementa SchedulerRepository usando el comando 'at'
type AtScheduler struct {
	testMode bool
}

// Verificar que implementa la interfaz
var _ repositories.SchedulerRepository = (*AtScheduler)(nil)

// NewAtScheduler crea una nueva instancia
func NewAtScheduler() *AtScheduler {
	return &AtScheduler{
		testMode: false,
	}
}

// NewAtSchedulerWithTestMode crea una instancia en modo prueba
func NewAtSchedulerWithTestMode(testMode bool) *AtScheduler {
	return &AtScheduler{
		testMode: testMode,
	}
}

// ScheduleShutdown programa un apagado del sistema
func (s *AtScheduler) ScheduleShutdown(t time.Time) error {
	if !s.IsAvailable() {
		return ErrAtNotAvailable
	}

	// Verificar que el filesystem permita escritura
	if !s.isFilesystemWritable() {
		return ErrFilesystemReadOnly
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
		command = "echo 'TEST MODE: Shutdown time reached' | wall"
	} else {
		command = "shutdown -h now"
	}

	// Ejecutar 'at'
	cmd := exec.Command("at", fmt.Sprintf("now + %d minutes", minutes))
	cmd.Stdin = strings.NewReader(command + "\n")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to schedule shutdown: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// CancelShutdown cancela todos los apagados programados
func (s *AtScheduler) CancelShutdown() error {
	if !s.IsAvailable() {
		return nil // Si 'at' no está disponible, no hay nada que cancelar
	}

	// Listar todos los trabajos
	jobs, err := s.ListScheduledJobs()
	if err != nil {
		return err
	}

	// Cancelar cada trabajo que sea un shutdown
	for _, job := range jobs {
		if s.isShutdownJob(job) {
			if err := s.cancelJob(job.ID); err != nil {
				// Continuar aunque falle cancelar uno
				continue
			}
		}
	}

	return nil
}

// ListScheduledJobs lista todas las tareas programadas
func (s *AtScheduler) ListScheduledJobs() ([]*repositories.ShutdownJob, error) {
	if !s.IsAvailable() {
		return nil, ErrAtNotAvailable
	}

	// Ejecutar 'atq' para listar trabajos
	cmd := exec.Command("atq")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	// Parsear salida
	jobs := s.parseAtqOutput(string(output))
	return jobs, nil
}

// IsAvailable verifica si el comando 'at' está disponible
func (s *AtScheduler) IsAvailable() bool {
	cmd := exec.Command("which", "at")
	err := cmd.Run()
	if err != nil {
		return false
	}

	// Verificar que el demonio atd está corriendo
	cmd = exec.Command("systemctl", "is-active", "atd")
	err = cmd.Run()
	return err == nil
}

// isFilesystemWritable verifica si el filesystem permite escritura en el directorio de 'at'
func (s *AtScheduler) isFilesystemWritable() bool {
	// Intentar crear un archivo temporal en el directorio de 'at'
	atDir := "/var/spool/cron/atjobs"
	testFile := filepath.Join(atDir, ".rtc_scheduler_test")

	// Verificar que el directorio existe
	if _, err := os.Stat(atDir); os.IsNotExist(err) {
		return false
	}

	// Intentar crear archivo temporal
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()

	// Intentar borrar el archivo temporal
	err = os.Remove(testFile)
	return err == nil
}

// parseAtqOutput parsea la salida del comando 'atq'
func (s *AtScheduler) parseAtqOutput(output string) []*repositories.ShutdownJob {
	var jobs []*repositories.ShutdownJob

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		job := s.parseAtqLine(line)
		if job != nil {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// parseAtqLine parsea una línea de salida de 'atq'
// Formato: 1	Sun Nov 17 22:00:00 2025 a user
func (s *AtScheduler) parseAtqLine(line string) *repositories.ShutdownJob {
	// Regex para parsear la línea
	re := regexp.MustCompile(`^(\d+)\s+\w+\s+(\w+\s+\d+\s+\d+:\d+:\d+\s+\d+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 3 {
		return nil
	}

	jobID := matches[1]
	timeStr := matches[2]

	// Parsear la fecha
	scheduledAt, err := time.Parse("Jan 2 15:04:05 2006", timeStr)
	if err != nil {
		return nil
	}

	return &repositories.ShutdownJob{
		ID:          jobID,
		ScheduledAt: scheduledAt,
		Command:     "shutdown", // Asumimos que es un shutdown
	}
}

// isShutdownJob verifica si un trabajo es un comando de apagado
func (s *AtScheduler) isShutdownJob(job *repositories.ShutdownJob) bool {
	// Obtener detalles del trabajo
	cmd := exec.Command("at", "-c", job.ID)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Buscar el comando 'shutdown' en la salida
	return strings.Contains(string(output), "shutdown")
}

// cancelJob cancela un trabajo específico
func (s *AtScheduler) cancelJob(jobID string) error {
	cmd := exec.Command("atrm", jobID)
	return cmd.Run()
}

// EnsureAtdRunning asegura que el demonio atd esté corriendo
func (s *AtScheduler) EnsureAtdRunning() error {
	// Verificar si está corriendo
	cmd := exec.Command("systemctl", "is-active", "atd")
	if err := cmd.Run(); err == nil {
		return nil // Ya está corriendo
	}

	// Intentar iniciar
	cmd = exec.Command("systemctl", "start", "atd")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start atd service: %w", err)
	}

	// Habilitar para arranque automático
	cmd = exec.Command("systemctl", "enable", "atd")
	cmd.Run() // Ignorar error aquí

	return nil
}

// GetJobDetails obtiene detalles completos de un trabajo
func (s *AtScheduler) GetJobDetails(jobID string) (string, error) {
	cmd := exec.Command("at", "-c", jobID)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get job details: %w", err)
	}
	return string(output), nil
}

// ScheduleAt programa un comando a una hora específica (método auxiliar)
func (s *AtScheduler) ScheduleAt(t time.Time, command string) error {
	if !s.IsAvailable() {
		return ErrAtNotAvailable
	}

	// Verificar que el filesystem permita escritura
	if !s.isFilesystemWritable() {
		return ErrFilesystemReadOnly
	}

	duration := time.Until(t)
	if duration < 0 {
		return ErrInvalidTime
	}

	minutes := int(duration.Minutes())
	if minutes < 1 {
		minutes = 1
	}

	cmd := exec.Command("at", fmt.Sprintf("now + %d minutes", minutes))
	cmd.Stdin = strings.NewReader(command + "\n")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to schedule command: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ParseJobID extrae el ID del trabajo de la salida de 'at'
func (s *AtScheduler) ParseJobID(output string) (string, error) {
	// Buscar patrón "job <ID>"
	re := regexp.MustCompile(`job\s+(\d+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 2 {
		return "", errors.New("could not parse job ID from output")
	}

	return matches[1], nil
}

// CountScheduledJobs cuenta cuántos trabajos están programados
func (s *AtScheduler) CountScheduledJobs() (int, error) {
	jobs, err := s.ListScheduledJobs()
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}