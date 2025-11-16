// internal/presentation/cli/commands.go
package cli

import (
	"flag"
	"fmt"
	"os"

	"rtc-scheduler/internal/application/usecases"
	"rtc-scheduler/pkg/logger"
)

// Command representa un comando CLI
type Command interface {
	Execute() error
	Validate() error
}

// CLI maneja la interfaz de línea de comandos
type CLI struct {
	installUC    *usecases.InstallServiceUseCase
	uninstallUC  *usecases.UninstallServiceUseCase
	scheduleUC   *usecases.SchedulePowerUseCase
	statusUC     *usecases.ShowStatusUseCase
	enableUC     *usecases.EnableServiceUseCase
	disableUC    *usecases.DisableServiceUseCase
	clearUC      *usecases.ClearAlarmUseCase
	runServiceUC *usecases.RunServiceUseCase
	logger       logger.Logger
}

// NewCLI crea una nueva instancia de CLI
func NewCLI(
	installUC *usecases.InstallServiceUseCase,
	uninstallUC *usecases.UninstallServiceUseCase,
	scheduleUC *usecases.SchedulePowerUseCase,
	statusUC *usecases.ShowStatusUseCase,
	enableUC *usecases.EnableServiceUseCase,
	disableUC *usecases.DisableServiceUseCase,
	clearUC *usecases.ClearAlarmUseCase,
	runServiceUC *usecases.RunServiceUseCase,
	log logger.Logger,
) *CLI {
	return &CLI{
		installUC:    installUC,
		uninstallUC:  uninstallUC,
		scheduleUC:   scheduleUC,
		statusUC:     statusUC,
		enableUC:     enableUC,
		disableUC:    disableUC,
		clearUC:      clearUC,
		runServiceUC: runServiceUC,
		logger:       log,
	}
}

// Run ejecuta la aplicación CLI
func (c *CLI) Run() error {
	// Definir flags
	install := flag.Bool("install", false, "Install service with wake and shutdown times")
	uninstall := flag.Bool("uninstall", false, "Uninstall service")
	status := flag.Bool("status", false, "Show current status")
	clear := flag.Bool("clear", false, "Clear wake alarm")
	enable := flag.Bool("enable", false, "Enable service")
	disable := flag.Bool("disable", false, "Disable service")
	runService := flag.Bool("run-service", false, "Run from service (internal use)")
	test := flag.Bool("test", false, "Test mode (no real shutdown)")
	version := flag.Bool("version", false, "Show version")

	wakeTime := flag.String("wake", "", "Wake time (HH:MM)")
	shutdownTime := flag.String("shutdown", "", "Shutdown time (HH:MM)")

	flag.Parse()

	// Mostrar versión
	if *version {
		fmt.Println("rtc-scheduler v1.0.0")
		return nil
	}

	// Verificar permisos de root (excepto para status y version)
	if os.Geteuid() != 0 && !*status && !*version {
		return fmt.Errorf("❌ This program must be run as root (sudo)")
	}

	// Routing de comandos
	switch {
	case *install:
		return c.handleInstall(*wakeTime, *shutdownTime)
	case *uninstall:
		return c.handleUninstall()
	case *status:
		return c.handleStatus()
	case *clear:
		return c.handleClear()
	case *enable:
		return c.handleEnable()
	case *disable:
		return c.handleDisable()
	case *runService:
		return c.handleRunService()
	default:
		// Modo manual (programación única)
		if *wakeTime == "" || *shutdownTime == "" {
			c.showUsage()
			return fmt.Errorf("wake and shutdown times are required for manual scheduling")
		}
		return c.handleManualSchedule(*wakeTime, *shutdownTime, *test)
	}
}

// showUsage muestra la ayuda de uso
func (c *CLI) showUsage() {
	fmt.Println("RTC Scheduler - Power management for Linux systems")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  rtc-scheduler [flags]")
	fmt.Println()
	fmt.Println("SERVICE MANAGEMENT:")
	fmt.Println("  -install -wake HH:MM -shutdown HH:MM    Install and enable service")
	fmt.Println("  -uninstall                              Uninstall service")
	fmt.Println("  -enable                                 Enable service")
	fmt.Println("  -disable                                Disable service")
	fmt.Println("  -status                                 Show status")
	fmt.Println()
	fmt.Println("MANUAL SCHEDULING:")
	fmt.Println("  -wake HH:MM -shutdown HH:MM             Schedule once")
	fmt.Println("  -wake HH:MM -shutdown HH:MM -test       Schedule once (test mode)")
	fmt.Println()
	fmt.Println("MAINTENANCE:")
	fmt.Println("  -clear                                  Clear wake alarm")
	fmt.Println("  -version                                Show version")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  sudo ./rtc-scheduler -install -wake 08:00 -shutdown 22:00")
	fmt.Println("  sudo ./rtc-scheduler -status")
	fmt.Println("  sudo ./rtc-scheduler -wake 08:00 -shutdown 22:00 -test")
}