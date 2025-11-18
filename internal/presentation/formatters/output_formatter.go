// internal/presentation/formatters/output_formatter.go
package formatters

import (
	"fmt"

	"rtc-scheduler/internal/application/usecases"
)

// OutputFormatter formatea la salida para la consola
type OutputFormatter struct{}

// NewOutputFormatter crea una nueva instancia
func NewOutputFormatter() *OutputFormatter {
	return &OutputFormatter{}
}

// PrintInstallSuccess imprime mensaje de instalación exitosa
func (f *OutputFormatter) PrintInstallSuccess(wake, shutdown string) {
	fmt.Println("✅ Service installed and enabled successfully")
	fmt.Println()
	fmt.Println("📅 Schedule Configuration:")
	fmt.Printf("   Wake time:     %s\n", wake)
	fmt.Printf("   Shutdown time: %s\n", shutdown)
	fmt.Println()
	fmt.Println("💡 Useful Commands:")
	fmt.Println("   rtc-scheduler -status          # View status")
	fmt.Println("   sudo rtc-scheduler -disable    # Pause temporarily")
	fmt.Println("   sudo rtc-scheduler -enable     # Reactivate")
	fmt.Println("   sudo rtc-scheduler -uninstall  # Uninstall completely")
	fmt.Println()
	fmt.Println("🔄 The system will automatically:")
	fmt.Println("   • Shutdown at", shutdown, "every day")
	fmt.Println("   • Power on at", wake, "every day")
}

// PrintUninstallSuccess imprime mensaje de desinstalación exitosa
func (f *OutputFormatter) PrintUninstallSuccess(output *usecases.UninstallServiceOutput) {
	fmt.Println("✅ Service uninstalled successfully")
	fmt.Println()
	fmt.Println("Cleanup summary:")
	if output.ServiceUninstalled {
		fmt.Println("   ✓ Systemd service removed")
	}
	if output.ConfigDeleted {
		fmt.Println("   ✓ Configuration deleted")
	}
	if output.AlarmsCleared {
		fmt.Println("   ✓ RTC alarms cleared")
	}
	fmt.Println()
	fmt.Println("Note: The program binary was not removed.")
	fmt.Println("To remove it: sudo rm /usr/local/bin/rtc-scheduler")
}

// PrintStatus imprime el estado completo del sistema
func (f *OutputFormatter) PrintStatus(output *usecases.ShowStatusOutput) {
	fmt.Println(output.Message)
}


// PrintScheduleSuccess imprime mensaje de programación exitosa
func (f *OutputFormatter) PrintScheduleSuccess(output *usecases.SchedulePowerOutput, testMode bool) {
	fmt.Println("📅 Power schedule configured:")
	fmt.Printf("   Next wake:     %s\n", output.Schedule.WakeTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Next shutdown: %s\n", output.Schedule.ShutdownTime.Format("2006-01-02 15:04:05"))
	fmt.Println()

	fmt.Println("✅ RTC wake alarm configured")

	if testMode {
		fmt.Println("✅ Shutdown scheduled (TEST MODE - will not actually shutdown)")
	} else {
		fmt.Println("✅ Shutdown scheduled")
	}

	fmt.Println()
	fmt.Println("⚠️  Note: This is a one-time schedule.")
	fmt.Println("   For recurring schedules, use: sudo rtc-scheduler -install -wake HH:MM -shutdown HH:MM")
}

// PrintError imprime un mensaje de error formateado
func (f *OutputFormatter) PrintError(err error) {
	fmt.Printf("❌ Error: %v\n", err)
}
