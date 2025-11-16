# RTC Power Scheduler

A comprehensive power management system for Linux systems using RTC (Real-Time Clock) hardware and systemd services. Built with hexagonal architecture and SOLID principles.

## 🎯 Features

- ✅ **Automatic Power Scheduling**: Daily wake-up and shutdown cycles
- ✅ **RTC Hardware Integration**: Reliable wake-up using hardware clock
- ✅ **Systemd Service**: Automatic recurring schedules
- ✅ **Rich CLI Interface**: Full command-line control
- ✅ **Test Mode**: Safe testing without actual shutdowns
- ✅ **Hexagonal Architecture**: Clean separation of concerns
- ✅ **SOLID Principles**: Maintainable, extensible code
- ✅ **Zero Dependencies**: Pure Go standard library
- ✅ **Comprehensive Logging**: Structured logging throughout
- ✅ **Input Validation**: Robust time format validation

## 📋 Requirements

- **Linux System**: Any Linux distribution with systemd
- **RTC Hardware**: Real-Time Clock device at `/sys/class/rtc/rtc0`
- **Go 1.21+**: For building from source
- **System Packages**:
  - `at` command: `sudo apt install at`
  - `systemctl` (systemd)
- **Privileges**: Root access required for most operations

## 🚀 Quick Start

### 1. Build and Install

```bash
# Clone the repository
git clone https://github.com/yourusername/rtc-scheduler
cd rtc-scheduler

# Build and install system-wide
make install

# Or build manually:
make build
sudo cp bin/rtc-scheduler /usr/local/bin/
sudo chmod +x /usr/local/bin/rtc-scheduler
```

### 2. Install Recurring Schedule

```bash
# Install service (wake at 8am, shutdown at 10pm daily)
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00
```

### 3. Check Status

```bash
# View current configuration and status (no sudo needed)
rtc-scheduler -status
```

## 📖 Usage

### Service Management

```bash
# Install with daily schedule (recommended)
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00

# Enable/disable service (keeps configuration)
sudo rtc-scheduler -enable
sudo rtc-scheduler -disable

# Uninstall completely
sudo rtc-scheduler -uninstall
```

### Manual Scheduling (One-time)

```bash
# Schedule one-time power cycle
sudo rtc-scheduler -wake 08:00 -shutdown 22:00

# Test mode (doesn't actually shutdown)
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test
```

### Status and Maintenance

```bash
# Show comprehensive status (no sudo needed)
rtc-scheduler -status

# Clear wake alarm
sudo rtc-scheduler -clear

# Show version information
rtc-scheduler -version
```

### CLI Examples

```bash
# Install recurring schedule (8am wake, 10pm shutdown)
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00

# Check comprehensive status
rtc-scheduler -status

# Temporarily disable (keeps configuration)
sudo rtc-scheduler -disable

# Re-enable service
sudo rtc-scheduler -enable

# Uninstall completely
sudo rtc-scheduler -uninstall

# One-time manual schedule
sudo rtc-scheduler -wake 08:00 -shutdown 22:00

# Test mode (doesn't actually shutdown)
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test

# Clear wake alarm
sudo rtc-scheduler -clear

# Show version information
rtc-scheduler -version
```

## 🏗️ Architecture

This project implements **Hexagonal Architecture** (Ports & Adapters) with clear separation of concerns:

```
rtc-scheduler/
├── cmd/rtc-scheduler/           # Application entry point
├── internal/
│   ├── domain/                  # Core business logic
│   │   ├── entities/           # Domain models (Config, Schedule)
│   │   ├── repositories/       # Interfaces (ports/contracts)
│   │   └── services/           # Domain services
│   ├── application/            # Application layer
│   │   ├── usecases/          # Business use cases
│   │   └── dto/               # Data transfer objects
│   ├── infrastructure/         # External implementations (adapters)
│   │   ├── rtc/               # RTC hardware access
│   │   ├── config/            # JSON configuration storage
│   │   ├── systemd/           # Systemd service management
│   │   └── scheduler/         # 'at' command scheduling
│   └── presentation/           # User interface adapters
│       ├── cli/               # Command-line interface
│       └── formatters/        # Output formatting
├── pkg/                        # Shared packages
│   ├── logger/                # Structured logging
│   └── errors/                # Custom error types
├── configs/                    # Default configuration files
├── Makefile                    # Build automation
└── rtc_scheduler_test.go       # Basic tests
```

### Design Patterns Used

- **Repository Pattern**: Abstract data access
- **Dependency Injection**: Loose coupling
- **Factory Pattern**: Object creation
- **Command Pattern**: CLI commands
- **Strategy Pattern**: Different scheduling strategies

### SOLID Principles

- **S**ingle Responsibility: Each module has one reason to change
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Implementations are interchangeable
- **I**nterface Segregation: Small, focused interfaces
- **D**ependency Inversion: Depend on abstractions

## 🧪 Development & Testing

```bash
# Build the application
make build

# Run tests with race detection
make test

# Generate coverage report
make coverage

# Format code
make fmt

# Run static analysis
make vet

# Download/update dependencies
make deps

# Run linter (if golangci-lint is installed)
make lint

# Clean build artifacts
make clean

# Show all available targets
make help
```


## 📊 How It Works

### Service Installation Flow

1. **Installation**: `rtc-scheduler -install -wake HH:MM -shutdown HH:MM`
   - Creates JSON configuration file (`/etc/rtc-scheduler.json`)
   - Installs systemd service unit
   - Enables and starts the service

2. **Automatic Execution**: At system boot, systemd runs the service
   - Service reads configuration from JSON file
   - Sets RTC wake alarm for next scheduled time
   - Schedules system shutdown using `at` command
   - Service exits (job is scheduled)

3. **Power Cycle**:
   - System shuts down at scheduled time (via `at`)
   - RTC alarm wakes system at scheduled time
   - System boots, systemd service runs again
   - Cycle repeats daily

### Manual Scheduling

For one-time schedules: `rtc-scheduler -wake HH:MM -shutdown HH:MM`
- Immediately sets RTC wake alarm
- Schedules shutdown using `at`
- No service installation required

## 🛠️ Troubleshooting

### RTC not available

```bash
# Check if RTC exists
ls -l /sys/class/rtc/rtc0

# Check dmesg for RTC messages
dmesg | grep rtc
```

### 'at' command not working

```bash
# Install at package
sudo apt install at

# Enable and start atd service
sudo systemctl enable atd
sudo systemctl start atd
```

### Service not starting

```bash
# Check service status
sudo systemctl status rtc-scheduler

# View logs
sudo journalctl -u rtc-scheduler -n 50
```

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/EstebanJS/rtc-scheduler
cd rtc-scheduler

# Install Go 1.21+
# Build and test
make deps
make test
make build

# Run linter (optional)
make lint
```

## 📞 Support

For issues and questions, please open an issue on GitHub.

## 🙏 Acknowledgments

- **Hexagonal Architecture**: Inspired by Alistair Cockburn
- **Clean Architecture**: Robert C. Martin
- **SOLID Principles**: Robert C. Martin
- **Go Community**: For the excellent standard library
- **Linux RTC Subsystem**: For reliable hardware clock access