// internal/application/usecases/show_status.go
package usecases

import (
	"fmt"
	"time"

	"rtc-scheduler/internal/domain/repositories"
	"rtc-scheduler/pkg/logger"
)

type ShowStatusInput struct{}

type ShowStatusOutput struct {
	ServiceInstalled   bool
	ServiceEnabled     bool
	ServiceRunning     bool
	ConfigExists       bool
	WakeTime           string
	ShutdownTime       string
	Enabled            bool
	RTCWakeAlarm       string
	RTCCurrentTime     string
	SystemTime         string
	ScheduledJobs      []*repositories.ShutdownJob
	Message            string
}

// ShowStatusUseCase muestra el estado completo del sistema
type ShowStatusUseCase struct {
	rtcRepo       repositories.RTCRepository
	configRepo    repositories.ConfigRepository
	serviceRepo   repositories.ServiceRepository
	schedulerRepo repositories.SchedulerRepository
	logger        logger.Logger
}

func NewShowStatusUseCase(
	rtc repositories.RTCRepository,
	config repositories.ConfigRepository,
	service repositories.ServiceRepository,
	scheduler repositories.SchedulerRepository,
	log logger.Logger,
) *ShowStatusUseCase {
	return &ShowStatusUseCase{
		rtcRepo:       rtc,
		configRepo:    config,
		serviceRepo:   service,
		schedulerRepo: scheduler,
		logger:        log,
	}
}

func (uc *ShowStatusUseCase) Execute(input *ShowStatusInput) (*ShowStatusOutput, error) {
	uc.logger.Info("Gathering system status")

	output := &ShowStatusOutput{}

	// Verificar servicio
	output.ServiceInstalled = uc.serviceRepo.IsInstalled()
	if output.ServiceInstalled {
		status, err := uc.serviceRepo.Status()
		if err == nil {
			output.ServiceEnabled = status.IsEnabled
			output.ServiceRunning = status.IsRunning
		}
	}

	// Verificar configuración
	output.ConfigExists = uc.configRepo.Exists()
	if output.ConfigExists {
		config, err := uc.configRepo.Load()
		if err == nil {
			output.WakeTime = config.WakeTime
			output.ShutdownTime = config.ShutdownTime
			output.Enabled = config.Enabled
		}
	}

	// Estado RTC
	if uc.rtcRepo.IsAvailable() {
		if wakeTime, err := uc.rtcRepo.GetWakeAlarm(); err == nil {
			output.RTCWakeAlarm = wakeTime.Format("2006-01-02 15:04:05")
		} else {
			output.RTCWakeAlarm = "Not set"
		}

		if rtcTime, err := uc.rtcRepo.GetCurrentTime(); err == nil {
			output.RTCCurrentTime = rtcTime.Format("2006-01-02 15:04:05")
		}
	} else {
		output.RTCWakeAlarm = "RTC not available"
		output.RTCCurrentTime = "RTC not available"
	}

	// Hora del sistema
	output.SystemTime = time.Now().Format("2006-01-02 15:04:05")

	// Tareas programadas
	if uc.schedulerRepo.IsAvailable() {
		jobs, err := uc.schedulerRepo.ListScheduledJobs()
		if err == nil {
			output.ScheduledJobs = jobs
		}
	}

	// Generar mensaje de resumen
	output.Message = uc.generateStatusMessage(output)

	uc.logger.Info("Status gathered successfully")
	return output, nil
}

func (uc *ShowStatusUseCase) generateStatusMessage(output *ShowStatusOutput) string {
	var msg string

	msg += "📊 RTC Scheduler Status\n"
	msg += "═══════════════════════════════════════\n\n"

	// Servicio
	msg += "⚙️  Service:\n"
	if output.ServiceInstalled {
		status := "❌ Stopped"
		if output.ServiceRunning {
			status = "✅ Running"
		}
		msg += fmt.Sprintf("   Installed: ✅ Yes\n")
		msg += fmt.Sprintf("   Status: %s\n", status)
		msg += fmt.Sprintf("   Enabled: %s\n", map[bool]string{true: "✅ Yes", false: "❌ No"}[output.ServiceEnabled])
	} else {
		msg += "   Installed: ❌ No\n"
	}
	msg += "\n"

	// Configuración
	msg += "🔧 Configuration:\n"
	if output.ConfigExists {
		msg += fmt.Sprintf("   Exists: ✅ Yes\n")
		msg += fmt.Sprintf("   Wake Time: %s\n", output.WakeTime)
		shutdownStatus := "Not configured"
		if output.ShutdownTime != "" {
			shutdownStatus = output.ShutdownTime
		}
		msg += fmt.Sprintf("   Shutdown Time: %s\n", shutdownStatus)
		msg += fmt.Sprintf("   Enabled: %s\n", map[bool]string{true: "✅ Yes", false: "❌ No"}[output.Enabled])
	} else {
		msg += "   Exists: ❌ No\n"
	}
	msg += "\n"

	// RTC
	msg += "🕐 RTC (Hardware Clock):\n"
	msg += fmt.Sprintf("   Available: %s\n", map[bool]string{true: "✅ Yes", false: "❌ No"}[output.RTCCurrentTime != "RTC not available"])
	msg += fmt.Sprintf("   Current Time: %s\n", output.RTCCurrentTime)
	msg += fmt.Sprintf("   Wake Alarm: %s\n", output.RTCWakeAlarm)
	msg += "\n"

	// Sistema
	msg += "💻 System:\n"
	msg += fmt.Sprintf("   Current Time: %s\n", output.SystemTime)
	msg += "\n"

	// Tareas programadas
	msg += "⏰ Scheduled Jobs:\n"
	if len(output.ScheduledJobs) > 0 {
		for i, job := range output.ScheduledJobs {
			msg += fmt.Sprintf("   %d. %s - %s\n", i+1, job.ScheduledAt.Format("2006-01-02 15:04:05"), job.Command)
		}
	} else {
		msg += "   None\n"
	}

	msg += "\n═══════════════════════════════════════"

	return msg
}