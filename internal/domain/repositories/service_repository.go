package repositories

type ServiceRepository interface {
	Install(executablePath string) error
	Uninstall() error
	Enable() error
	Disable() error
	Start() error
	Stop() error
	Status() (*ServiceStatus, error)
	IsInstalled() bool
}

type ServiceStatus struct {
	Name      string
	IsRunning bool
	IsEnabled bool
	Error     error
}