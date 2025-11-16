// internal/presentation/cli/handlers.go
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"rtc-scheduler/internal/application/usecases"
)

// handleInstall maneja la instalación del servicio
func (c *CLI) handleInstall(wakeTime, shutdownTime string) error {
	if wakeTime == "" || shutdownTime == "" {
		return fmt.Errorf("❌ Wake time and shutdown time are required for installation")
	}

	c.logger.Info("Installing service", "wake_time", wakeTime, "shutdown_time", shutdownTime)

	// Obtener ruta del ejecutable
	execPath, err := c.getExecutablePath()
	if err != nil {
		return fmt.Errorf("❌ Failed to get executable path: %w", err)
	}

	// Ejecutar caso de uso
	input := &usecases.InstallServiceInput{
		WakeTime:     wakeTime,
		ShutdownTime: shutdownTime,
		ExecutablePath: execPath,
	}

	output, err := c.installUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Installation failed: %w", err)
	}

	fmt.Println("✅", output.Message)
	fmt.Println()
	fmt.Println("💡 Useful commands:")
	fmt.Println("   sudo rtc-scheduler -status          # Show status")
	fmt.Println("   sudo rtc-scheduler -disable         # Temporarily disable")
	fmt.Println("   sudo rtc-scheduler -enable          # Re-enable")
	fmt.Println("   sudo rtc-scheduler -uninstall       # Complete uninstall")

	return nil
}

// handleUninstall maneja la desinstalación del servicio
func (c *CLI) handleUninstall() error {
	c.logger.Info("Uninstalling service")

	input := &usecases.UninstallServiceInput{}
	output, err := c.uninstallUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Uninstallation failed: %w", err)
	}

	fmt.Println("✅", output.Message)
	return nil
}

// handleStatus muestra el estado del sistema
func (c *CLI) handleStatus() error {
	c.logger.Info("Showing status")

	input := &usecases.ShowStatusInput{}
	output, err := c.statusUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Failed to get status: %w", err)
	}

	fmt.Println(output.Message)
	return nil
}

// handleClear limpia la alarma de encendido
func (c *CLI) handleClear() error {
	c.logger.Info("Clearing wake alarm")

	input := &usecases.ClearAlarmInput{}
	output, err := c.clearUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Failed to clear alarm: %w", err)
	}

	fmt.Println("✅", output.Message)
	return nil
}

// handleEnable habilita el servicio
func (c *CLI) handleEnable() error {
	c.logger.Info("Enabling service")

	input := &usecases.EnableServiceInput{}
	output, err := c.enableUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Failed to enable service: %w", err)
	}

	fmt.Println("✅", output.Message)
	return nil
}

// handleDisable deshabilita el servicio
func (c *CLI) handleDisable() error {
	c.logger.Info("Disabling service")

	input := &usecases.DisableServiceInput{}
	output, err := c.disableUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Failed to disable service: %w", err)
	}

	fmt.Println("✅", output.Message)
	return nil
}

// handleRunService ejecuta desde el servicio systemd
func (c *CLI) handleRunService() error {
	c.logger.Info("Running from service")

	input := &usecases.RunServiceInput{}
	output, err := c.runServiceUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Service execution failed: %w", err)
	}

	fmt.Println("✅", output.Message)
	return nil
}

// handleManualSchedule maneja la programación manual (una sola vez)
func (c *CLI) handleManualSchedule(wakeTime, shutdownTime string, testMode bool) error {
	c.logger.Info("Manual scheduling", "wake_time", wakeTime, "shutdown_time", shutdownTime, "test_mode", testMode)

	input := &usecases.SchedulePowerInput{
		WakeTime:     wakeTime,
		ShutdownTime: shutdownTime,
		TestMode:     testMode,
	}

	output, err := c.scheduleUC.Execute(input)
	if err != nil {
		return fmt.Errorf("❌ Scheduling failed: %w", err)
	}

	fmt.Println("✅", output.Message)
	if output.Schedule != nil {
		fmt.Printf("   Next wake: %s\n", output.Schedule.WakeTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Next shutdown: %s\n", output.Schedule.ShutdownTime.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// getExecutablePath obtiene la ruta del ejecutable actual
func (c *CLI) getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Resolver enlaces simbólicos
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}

	// Convertir a ruta absoluta
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	// Verificar que existe
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return "", fmt.Errorf("executable not found: %s", execPath)
	}

	return execPath, nil
}