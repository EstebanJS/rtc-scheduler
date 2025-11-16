// internal/infrastructure/systemd/systemd_service.go
package systemd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"rtc-scheduler/internal/domain/repositories"
)

const (
	serviceName     = "rtc-scheduler.service"
	systemdPath     = "/etc/systemd/system"
	serviceFileName = serviceName
)

var (
	ErrServiceNotInstalled = errors.New("service is not installed")
	ErrSystemdNotAvailable = errors.New("systemd is not available")
)

// SystemdService implementa ServiceRepository usando systemd
type SystemdService struct {
	servicePath string
}

// Verificar que implementa la interfaz
var _ repositories.ServiceRepository = (*SystemdService)(nil)

// NewSystemdService crea una nueva instancia
func NewSystemdService() *SystemdService {
	return &SystemdService{
		servicePath: fmt.Sprintf("%s/%s", systemdPath, serviceFileName),
	}
}

// Install instala el servicio systemd
func (s *SystemdService) Install(executablePath string) error {
	// Verificar que systemd está disponible
	if !s.isSystemdAvailable() {
		return ErrSystemdNotAvailable
	}

	// Crear contenido del archivo de servicio
	serviceContent := s.generateServiceContent(executablePath)

	// Escribir archivo
	if err := os.WriteFile(s.servicePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Recargar systemd
	if err := s.daemonReload(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	return nil
}

// Uninstall desinstala el servicio
func (s *SystemdService) Uninstall() error {
	if !s.IsInstalled() {
		return nil // Ya está desinstalado
	}

	// Detener el servicio si está corriendo
	s.Stop()

	// Deshabilitar el servicio
	s.Disable()

	// Eliminar archivo
	if err := os.Remove(s.servicePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Recargar systemd
	return s.daemonReload()
}

// Enable habilita el servicio para arranque automático
func (s *SystemdService) Enable() error {
	if !s.IsInstalled() {
		return ErrServiceNotInstalled
	}
	return s.runSystemctl("enable", serviceName)
}

// Disable deshabilita el servicio
func (s *SystemdService) Disable() error {
	if !s.IsInstalled() {
		return nil
	}
	return s.runSystemctl("disable", serviceName)
}

// Start inicia el servicio
func (s *SystemdService) Start() error {
	if !s.IsInstalled() {
		return ErrServiceNotInstalled
	}
	return s.runSystemctl("start", serviceName)
}

// Stop detiene el servicio
func (s *SystemdService) Stop() error {
	if !s.IsInstalled() {
		return nil
	}
	return s.runSystemctl("stop", serviceName)
}

// Restart reinicia el servicio
func (s *SystemdService) Restart() error {
	if !s.IsInstalled() {
		return ErrServiceNotInstalled
	}
	return s.runSystemctl("restart", serviceName)
}

// Status obtiene el estado del servicio
func (s *SystemdService) Status() (*repositories.ServiceStatus, error) {
	status := &repositories.ServiceStatus{
		Name: serviceName,
	}

	if !s.IsInstalled() {
		status.Error = ErrServiceNotInstalled
		return status, nil
	}

	// Verificar si está corriendo
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	status.IsRunning = err == nil && strings.TrimSpace(string(output)) == "active"

	// Verificar si está habilitado
	cmd = exec.Command("systemctl", "is-enabled", serviceName)
	output, err = cmd.Output()
	status.IsEnabled = err == nil && strings.TrimSpace(string(output)) == "enabled"

	return status, nil
}

// IsInstalled verifica si el servicio está instalado
func (s *SystemdService) IsInstalled() bool {
	_, err := os.Stat(s.servicePath)
	return !os.IsNotExist(err)
}

// generateServiceContent genera el contenido del archivo de servicio
func (s *SystemdService) generateServiceContent(executablePath string) string {
	return fmt.Sprintf(`[Unit]
Description=RTC Power Schedule Manager
Documentation=https://github.com/yourusername/rtc-scheduler
After=network.target time-sync.target
Wants=atd.service
# systemd-run is available in most systemd installations, no need for Wants

[Service]
Type=oneshot
User=root
ExecStart=%s -run-service
RemainAfterExit=yes
StandardOutput=journal
StandardError=journal
Restart=no

# Permisos necesarios para RTC y scheduling
PrivateTmp=yes
NoNewPrivileges=no
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/sys/class/rtc/rtc0 /etc/rtc-scheduler.json /var/spool/cron/atjobs
CapabilityBoundingSet=CAP_SYS_ADMIN

# Environment
Environment=SYSTEMD_LOG_LEVEL=info
Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

[Install]
WantedBy=multi-user.target
`, executablePath)
}

// runSystemctl ejecuta un comando systemctl
func (s *SystemdService) runSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s failed: %s", strings.Join(args, " "), string(output))
	}
	return nil
}

// daemonReload recarga la configuración de systemd
func (s *SystemdService) daemonReload() error {
	return s.runSystemctl("daemon-reload")
}

// isSystemdAvailable verifica si systemd está disponible en el sistema
func (s *SystemdService) isSystemdAvailable() bool {
	cmd := exec.Command("systemctl", "--version")
	return cmd.Run() == nil
}

// GetLogs obtiene los logs del servicio
func (s *SystemdService) GetLogs(lines int) (string, error) {
	if !s.IsInstalled() {
		return "", ErrServiceNotInstalled
	}

	cmd := exec.Command("journalctl", "-u", serviceName, "-n", fmt.Sprintf("%d", lines), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return string(output), nil
}