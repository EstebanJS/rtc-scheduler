// cmd/rtc-scheduler/main.go
package main

import (
	"fmt"
	"os"

	"rtc-scheduler/internal/application/usecases"
	"rtc-scheduler/internal/infrastructure/config"
	"rtc-scheduler/internal/infrastructure/rtc"
	"rtc-scheduler/internal/infrastructure/scheduler"
	"rtc-scheduler/internal/infrastructure/systemd"
	"rtc-scheduler/internal/presentation/cli"
	"rtc-scheduler/pkg/logger"
)

const (
	configFilePath = "/etc/rtc-scheduler.json"
	version        = "1.0.11"
)

var (
	buildTime string
	gitCommit string
)

func main() {
	// Inicializar logger
	log := logger.New()

	// Verificar versión
	if len(os.Args) > 1 && (os.Args[1] == "-version" || os.Args[1] == "--version") {
		fmt.Printf("rtc-scheduler version %s", version)
		if buildTime != "" {
			fmt.Printf(" (built %s", buildTime)
			if gitCommit != "" {
				fmt.Printf(", commit %s", gitCommit)
			}
			fmt.Printf(")")
		}
		fmt.Println()
		os.Exit(0)
	}

	// Crear contenedor de dependencias
	container := initializeDependencies(log)

	// Crear CLI
	cliApp := cli.NewCLI(
		container.installUC,
		container.uninstallUC,
		container.scheduleUC,
		container.statusUC,
		container.enableUC,
		container.disableUC,
		container.clearUC,
		container.runServiceUC,
		log,
	)

	// Ejecutar aplicación
	if err := cliApp.Run(); err != nil {
		log.Error("Application error", "error", err)
		os.Exit(1)
	}
}

// DependencyContainer contiene todas las dependencias de la aplicación
type DependencyContainer struct {
	// Repositories
	rtcRepo       *rtc.LinuxRTC
	configRepo    *config.JSONConfigRepository
	serviceRepo   *systemd.SystemdService
	schedulerRepo *scheduler.HybridScheduler

	// Use Cases
	installUC    *usecases.InstallServiceUseCase
	uninstallUC  *usecases.UninstallServiceUseCase
	scheduleUC   *usecases.SchedulePowerUseCase
	statusUC     *usecases.ShowStatusUseCase
	enableUC     *usecases.EnableServiceUseCase
	disableUC    *usecases.DisableServiceUseCase
	clearUC      *usecases.ClearAlarmUseCase
	runServiceUC *usecases.RunServiceUseCase
}

// initializeDependencies inicializa todas las dependencias (Dependency Injection)
func initializeDependencies(log logger.Logger) *DependencyContainer {
	// Inicializar repositorios (Infrastructure Layer)
	rtcRepo := rtc.NewLinuxRTC()
	configRepo := config.NewJSONConfigRepository(configFilePath)
	serviceRepo := systemd.NewSystemdService()
	schedulerRepo := scheduler.NewHybridScheduler()

	// Verificar que componentes críticos estén disponibles
	if !rtcRepo.IsAvailable() {
		log.Warn("RTC device not available at /sys/class/rtc/rtc0")
		log.Warn("The program will continue but power scheduling will not work")
	}

	if !schedulerRepo.IsAvailable() {
		log.Warn("No suitable scheduler available")
		log.Warn("Neither 'at' command nor systemd-run found")
		log.Warn("Install at with: sudo apt install at")
		log.Warn("Shutdown scheduling will not work without scheduling capabilities")
	}

	// Inicializar casos de uso (Application Layer)
	installUC := usecases.NewInstallServiceUseCase(
		configRepo,
		serviceRepo,
		log,
	)

	uninstallUC := usecases.NewUninstallServiceUseCase(
		rtcRepo,
		configRepo,
		serviceRepo,
		schedulerRepo,
		log,
	)

	scheduleUC := usecases.NewSchedulePowerUseCase(
		rtcRepo,
		schedulerRepo,
		log,
	)

	statusUC := usecases.NewShowStatusUseCase(
		rtcRepo,
		configRepo,
		serviceRepo,
		schedulerRepo,
		log,
	)

	enableUC := usecases.NewEnableServiceUseCase(
		configRepo,
		serviceRepo,
		log,
	)

	disableUC := usecases.NewDisableServiceUseCase(
		configRepo,
		serviceRepo,
		schedulerRepo,
		rtcRepo,
		log,
	)

	clearUC := usecases.NewClearAlarmUseCase(
		rtcRepo,
		schedulerRepo,
		log,
	)

	runServiceUC := usecases.NewRunServiceUseCase(
		configRepo,
		rtcRepo,
		schedulerRepo,
		log,
	)

	return &DependencyContainer{
		rtcRepo:       rtcRepo,
		configRepo:    configRepo,
		serviceRepo:   serviceRepo,
		schedulerRepo: schedulerRepo,
		installUC:     installUC,
		uninstallUC:   uninstallUC,
		scheduleUC:    scheduleUC,
		statusUC:      statusUC,
		enableUC:      enableUC,
		disableUC:     disableUC,
		clearUC:       clearUC,
		runServiceUC:  runServiceUC,
	}
}