# 🕐 RTC Power Scheduler

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)](https://www.linux.org/)

A robust power management system for Linux systems using RTC hardware and systemd services. Features automatic daily wake-up/shutdown cycles with intelligent fallback mechanisms.

**Built with Hexagonal Architecture & SOLID Principles** ✨

## 📋 Table of Contents

- [🎯 Features](#-features)
- [📋 Requirements](#-requirements)
- [🚀 Quick Start](#-quick-start)
- [📖 Usage](#-usage)
- [🏗️ Architecture](#️-architecture)
- [🧪 Development](#-development)
- [📊 How It Works](#-how-it-works)
- [🛠️ Troubleshooting](#️-troubleshooting)
- [🤝 Contributing](#-contributing)
- [📞 Support](#-support)
- [🙏 Acknowledgments](#-acknowledgments)

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
- ✅ **Intelligent Fallbacks**: Works with or without `at` command
- ✅ **Cross-Platform**: ARM64 and AMD64 support

## 📋 Requirements

### System Requirements
- **Operating System**: Linux with systemd (Ubuntu, Debian, Raspbian, etc.)
- **Architecture**: AMD64 or ARM64
- **RTC Hardware**: Real-Time Clock device at `/sys/class/rtc/rtc0`

### Software Dependencies
- **Go 1.21+**: Required only for building from source
- **System Packages** (with intelligent fallbacks):
  - **Primary**: `at` command (`sudo apt install at`) - preferred for reliability
  - **Fallback**: `systemd-run` (included in systemd) - works in read-only environments

### Privileges
- **Root access** required for installation and scheduling operations
- **Regular user** can check status and view logs

## 🚀 Quick Start

### Option 1: Download Pre-built Binary

```bash
# Download latest release
wget https://github.com/EstebanJS/rtc-scheduler/releases/latest/download/rtc-scheduler-linux-arm64.tar.gz
tar -xzf rtc-scheduler-linux-arm64.tar.gz

# Install system-wide
sudo cp rtc-scheduler /usr/local/bin/
sudo chmod +x /usr/local/bin/rtc-scheduler
```

### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/EstebanJS/rtc-scheduler
cd rtc-scheduler

# Build and install
make install
```

### Install Daily Schedule

```bash
# Install service (wake at 8am, shutdown at 10pm daily)
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00

# Check status
rtc-scheduler -status
```

**🎉 That's it! Your system will now automatically wake up and shut down daily.**

## 📖 Usage

### 🔧 Service Management

| Command | Description | Example |
|---------|-------------|---------|
| `sudo rtc-scheduler -install` | Install with daily schedule | `sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00` |
| `sudo rtc-scheduler -enable` | Enable service (keeps config) | `sudo rtc-scheduler -enable` |
| `sudo rtc-scheduler -disable` | Disable service (keeps config) | `sudo rtc-scheduler -disable` |
| `sudo rtc-scheduler -uninstall` | Remove service completely | `sudo rtc-scheduler -uninstall` |

### ⏰ Manual Scheduling (One-time)

| Command | Description | Example |
|---------|-------------|---------|
| `sudo rtc-scheduler -wake HH:MM -shutdown HH:MM` | One-time power cycle | `sudo rtc-scheduler -wake 08:00 -shutdown 22:00` |
| `sudo rtc-scheduler -wake HH:MM -shutdown HH:MM -test` | Test mode (safe) | `sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test` |

### 📊 Status & Information

| Command | Description | Requires Sudo |
|---------|-------------|---------------|
| `rtc-scheduler -status` | Show comprehensive status | ❌ No |
| `rtc-scheduler -version` | Show version information | ❌ No |
| `sudo rtc-scheduler -clear` | Clear wake alarm | ✅ Yes |

### 💡 Complete Examples

```bash
# 🏠 Install recurring schedule (8am wake, 10pm shutdown)
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00

# 📊 Check comprehensive status
rtc-scheduler -status

# ⏸️ Temporarily disable (keeps configuration)
sudo rtc-scheduler -disable

# ▶️ Re-enable service
sudo rtc-scheduler -enable

# 🗑️ Uninstall completely
sudo rtc-scheduler -uninstall

# 🔄 One-time manual schedule
sudo rtc-scheduler -wake 08:00 -shutdown 22:00

# 🧪 Test mode (doesn't actually shutdown)
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test

# 🧹 Clear wake alarm
sudo rtc-scheduler -clear

# ℹ️ Show version information
rtc-scheduler -version
```

## 🏗️ Architecture

This project implements **Hexagonal Architecture** (Ports & Adapters) with clean separation of concerns and SOLID principles.

### 📁 Project Structure

```
rtc-scheduler/
├── cmd/rtc-scheduler/           # 🚀 Application entry point
├── internal/
│   ├── domain/                  # 🎯 Core business logic
│   │   ├── entities/           # 📦 Domain models (Config, Schedule)
│   │   ├── repositories/       # 🔌 Interfaces (ports/contracts)
│   │   └── services/           # ⚙️ Domain services
│   ├── application/            # 🎪 Application layer
│   │   ├── usecases/          # 🎬 Business use cases
│   │   └── dto/               # 📋 Data transfer objects
│   ├── infrastructure/         # 🔧 External implementations (adapters)
│   │   ├── rtc/               # 🕐 RTC hardware access
│   │   ├── config/            # 💾 JSON configuration storage
│   │   ├── systemd/           # 🔄 Systemd service management
│   │   └── scheduler/         # ⏰ Command scheduling (at/systemd-run)
│   └── presentation/           # 💻 User interface adapters
│       ├── cli/               # ⌨️ Command-line interface
│       └── formatters/        # 📄 Output formatting
├── pkg/                        # 📚 Shared packages
│   ├── logger/                # 📝 Structured logging
│   └── errors/                # ⚠️ Custom error types
├── configs/                    # ⚙️ Default configuration files
├── Makefile                    # 🔨 Build automation
└── rtc_scheduler_test.go       # 🧪 Basic tests
```

### 🎨 Design Patterns

| Pattern | Usage | Benefit |
|---------|-------|---------|
| **Repository** | Abstract data access | Database/hardware independence |
| **Dependency Injection** | Loose coupling | Testable, maintainable code |
| **Factory** | Object creation | Encapsulated instantiation |
| **Command** | CLI commands | Extensible command structure |
| **Strategy** | Scheduling strategies | Runtime algorithm selection |
| **Hybrid Scheduler** | Intelligent fallbacks | Robust operation across environments |

### ✨ SOLID Principles

| Principle | Implementation | Benefit |
|-----------|----------------|---------|
| **S**ingle Responsibility | Each module has one purpose | Focused, maintainable code |
| **O**pen/Closed | Extensible via interfaces | New features without modification |
| **L**iskov Substitution | Interchangeable implementations | Reliable polymorphism |
| **I**nterface Segregation | Small, focused interfaces | Reduced coupling |
| **D**ependency Inversion | Depend on abstractions | Flexible, testable architecture |

## 🧪 Development & Testing

### 🚀 Build Commands

```bash
# Build the application
make build

# Build and install system-wide
make install

# Build for different architectures
make build-linux-amd64    # AMD64 Linux
make build-linux-arm64    # ARM64 Linux (Raspberry Pi)
```

### 🧪 Testing & Quality

```bash
# Run tests with race detection
make test

# Generate coverage report
make coverage

# Run static analysis
make vet

# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Download/update dependencies
make deps

# Clean build artifacts
make clean

# Show all available targets
make help
```

### 📦 Releases

Pre-built binaries are available for download from [GitHub Releases](https://github.com/EstebanJS/rtc-scheduler/releases):

- **Linux AMD64**: `rtc-scheduler-linux-amd64.tar.gz`
- **Linux ARM64**: `rtc-scheduler-linux-arm64.tar.gz` (Raspberry Pi compatible)

```bash
# Download and install latest release
wget https://github.com/EstebanJS/rtc-scheduler/releases/latest/download/rtc-scheduler-linux-arm64.tar.gz
tar -xzf rtc-scheduler-linux-arm64.tar.gz
sudo cp rtc-scheduler /usr/local/bin/
sudo chmod +x /usr/local/bin/rtc-scheduler
```


## 📊 How It Works

### 🔄 Power Management Cycle

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   System Boot   │ -> │  Service Runs   │ -> │  RTC Alarm Set  │
│   (Daily)       │    │  (Config Load)  │    │  (Wake Time)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       │                       │
         │                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ System Shutdown │ <- │   Timer Exec    │ <- │  Timer Created  │
│   (Scheduled)   │    │   (at/systemd)  │    │  (Shutdown)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 🚀 Service Installation Flow

1. **📦 Installation**: `sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00`
   - Creates JSON config: `/etc/rtc-scheduler.json`
   - Generates systemd service unit
   - Enables and starts automatic service

2. **⚡ Automatic Execution**: At boot, systemd runs the service
   - Loads configuration from JSON
   - Sets hardware RTC wake alarm
   - Creates shutdown timer (at/systemd-run)
   - Service completes (timers remain active)

3. **🔄 Daily Power Cycle**:
   - **Evening**: System shuts down at scheduled time
   - **Morning**: RTC alarm wakes system automatically
   - **Repeat**: Cycle continues daily without intervention

### ⏰ Manual Scheduling

For one-time operations: `sudo rtc-scheduler -wake HH:MM -shutdown HH:MM`
- ✅ Immediately configures RTC wake alarm
- ✅ Creates shutdown timer using available scheduler
- ✅ No permanent service installation required
- ✅ Safe testing with `-test` flag

## 🛠️ Troubleshooting

### 🔍 Common Issues & Solutions

| Issue | Symptoms | Solution |
|-------|----------|----------|
| **RTC not available** | `rtc-scheduler -status` shows RTC errors | Check hardware: `ls -l /sys/class/rtc/rtc0` |
| **'at' command not working** | Scheduling fails with 'at' errors | Install: `sudo apt install at && sudo systemctl enable atd` |
| **Service not starting** | `systemctl status` shows failed state | Check logs: `sudo journalctl -u rtc-scheduler -n 20` |
| **Permission denied** | RTC access fails | Run with sudo or fix permissions: `sudo chmod 666 /sys/class/rtc/rtc0/wakealarm` |
| **Command not found** | systemd-run fails to execute | Ensure PATH is set correctly (fixed in v1.0.11+) |

### 🔧 Detailed Diagnostics

#### RTC Hardware Check
```bash
# Check if RTC device exists
ls -l /sys/class/rtc/rtc0

# Check kernel messages for RTC
dmesg | grep -i rtc

# Test RTC wakealarm capability
sudo sh -c 'echo 0 > /sys/class/rtc/rtc0/wakealarm'
```

#### 'at' Command Setup
```bash
# Install at package
sudo apt update && sudo apt install at

# Enable and start atd service
sudo systemctl enable atd
sudo systemctl start atd

# Check atd status
sudo systemctl status atd
```

#### Service Debugging
```bash
# Check service status
sudo systemctl status rtc-scheduler

# View recent logs
sudo journalctl -u rtc-scheduler -n 20 --no-pager

# Check service file syntax
sudo systemd-analyze verify /etc/systemd/system/rtc-scheduler.service

# Reload systemd configuration
sudo systemctl daemon-reload
```

#### Manual Testing
```bash
# Test RTC functionality
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test

# Test service installation
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00
```

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### 🚀 Development Setup

```bash
# Clone the repository
git clone https://github.com/EstebanJS/rtc-scheduler
cd rtc-scheduler

# Set up development environment
make deps          # Download dependencies
make test          # Run tests
make build         # Build the application
make lint          # Run linter (optional)
```

### 📋 Contribution Guidelines

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### 🐛 Reporting Issues

- Use [GitHub Issues](https://github.com/EstebanJS/rtc-scheduler/issues) for bugs and feature requests
- Provide detailed steps to reproduce
- Include system information and logs when possible

## 📞 Support

### 📧 Getting Help

- 📖 **Documentation**: This README and inline code comments
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/EstebanJS/rtc-scheduler/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/EstebanJS/rtc-scheduler/discussions)
- 💬 **General Questions**: [GitHub Discussions](https://github.com/EstebanJS/rtc-scheduler/discussions)

### 📊 System Compatibility

| Platform | Status | Notes |
|----------|--------|-------|
| **Ubuntu 20.04+** | ✅ Fully Supported | Tested extensively |
| **Debian 11+** | ✅ Fully Supported | Primary development platform |
| **Raspbian/Raspberry Pi OS** | ✅ Fully Supported | ARM64 cross-compilation |
| **Other systemd Linux** | ✅ Should Work | Report issues if found |

## 🙏 Acknowledgments

### 🎯 Architectural Inspiration
- **Hexagonal Architecture**: Alistair Cockburn's Ports & Adapters pattern
- **Clean Architecture**: Robert C. Martin's layered architecture principles
- **SOLID Principles**: Robert C. Martin's design principles

### 🛠️ Technical Foundations
- **Go Programming Language**: Excellent standard library and tooling
- **Linux RTC Subsystem**: Reliable hardware clock access
- **systemd**: Modern Linux service management
- **Open Source Community**: Countless libraries and tools

### 🤝 Special Thanks
- **Beta Testers**: For valuable feedback and real-world testing
- **Contributors**: For code improvements and bug fixes
- **Linux Community**: For maintaining excellent documentation

---

**Made with ❤️ for the Linux community**