# RTC Power Scheduler v1.0.9 - Release Assets

This release fixes **systemd timer PATH environment issues** on Raspberry Pi systems. Adds proper PATH environment to service file so systemd-run can find systemctl and other system commands.

## 📦 Release Files

### Binaries
- **`rtc-scheduler-v1.0.9-linux-arm64.tar.gz`** - Compressed archive containing the binary and checksum
- **`rtc-scheduler-v1.0.9-linux-arm64`** - Standalone binary (Linux ARM64)
- **`rtc-scheduler-v1.0.9-linux-arm64.sha256`** - SHA256 checksum for verification

## 🔍 Verification

Verify the integrity of the downloaded binary:

```bash
# Extract the archive
tar -xzf rtc-scheduler-v1.0.9-linux-arm64.tar.gz

# Verify checksum
sha256sum -c rtc-scheduler-v1.0.9-linux-arm64.sha256

# Make executable and test
chmod +x rtc-scheduler-v1.0.9-linux-arm64
./rtc-scheduler-v1.0.9-linux-arm64 -version
```

Expected output:
```
rtc-scheduler version 1.0.9 (built 2025-11-16T16:57:40Z)
```

## 🔧 Troubleshooting

### Service Startup Issues

If the systemd service fails to start, check the logs:

```bash
# Check service status
sudo systemctl status rtc-scheduler

# View detailed logs
sudo journalctl -u rtc-scheduler -n 50 --no-pager

# Check if atd service is running (optional)
sudo systemctl status atd

# Ensure atd is enabled (optional)
sudo systemctl enable atd
sudo systemctl start atd
```

### Robust Error Handling Features

**v1.0.3 introduces intelligent fallback mechanisms:**

#### Graceful Degradation Mode
- **Filesystem read-only**: RTC alarm works, shutdown scheduling logs warning but continues
- **at command unavailable**: Automatically falls back to systemd timers
- **Partial failures**: System maintains core functionality even with component failures

#### Scheduler Fallback Priority
1. **AtScheduler** (preferred - if available and filesystem writable)
2. **SystemdTimerScheduler** (fallback - works in read-only environments)
3. **None** (RTC-only mode with clear warnings)

### Common Issues & Solutions

1. **Filesystem read-only (Raspberry Pi common issue)**:
   ```bash
   # Check filesystem status
   mount | grep " / "

   # If read-only, check for filesystem errors
   sudo dmesg | grep -i "read-only"

   # Remount as read-write (if safe)
   sudo mount -o remount,rw /
   ```

2. **Permission denied on RTC device**:
   ```bash
   sudo chmod 666 /sys/class/rtc/rtc0/wakealarm
   ```

3. **at command not available**:
   ```bash
   sudo apt update && sudo apt install at
   ```
   **Note**: v1.0.8 works without `at` using systemd timers as fallback with corrected systemctl command paths

4. **Service fails to start**:
   ```bash
   # Check service file syntax
   sudo systemd-analyze verify /etc/systemd/system/rtc-scheduler.service

   # Reload systemd configuration
   sudo systemctl daemon-reload

   # Check systemd-run availability
   which systemd-run
   ```

5. **Systemd timer issues**:
   ```bash
   # List active timers
   systemctl list-timers | grep rtc-scheduler

   # Check timer status
   systemctl status rtc-scheduler-shutdown-*
   ```

## 🚀 Quick Install

```bash
# Download and extract
wget https://github.com/EstebanJS/rtc-scheduler/releases/download/v1.0.9/rtc-scheduler-v1.0.9-linux-arm64.tar.gz
tar -xzf rtc-scheduler-v1.0.9-linux-arm64.tar.gz

# Install system-wide
sudo cp rtc-scheduler-v1.0.9-linux-arm64 /usr/local/bin/rtc-scheduler
sudo chmod +x /usr/local/bin/rtc-scheduler

# Verify installation
rtc-scheduler -version
```

## 📋 System Requirements

- **Architecture**: ARM64 (aarch64)
- **OS**: Linux with systemd
- **Dependencies** (with intelligent fallbacks):
  - **Primary**: `at` command (`sudo apt install at`) - preferred for reliability
  - **Fallback**: `systemd-run` (included in systemd) - works in read-only environments
  - **Required**: RTC hardware at `/sys/class/rtc/rtc0`
  - **Required**: Systemd services

**Note**: v1.0.3 automatically selects the best available scheduler. Works even without `at` command using systemd timers.

## 🛠️ Usage Examples

```bash
# Install service with daily schedule
sudo rtc-scheduler -install -wake 08:00 -shutdown 22:00

# Check status
rtc-scheduler -status

# Manual one-time scheduling
sudo rtc-scheduler -wake 08:00 -shutdown 22:00

# Test mode (safe testing)
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test
```

## 🔧 Build Information

- **Version**: v1.0.9
- **Build Time**: 2025-11-16T16:57:40Z
- **Go Version**: 1.21
- **Build Flags**: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` (cross-compiled, statically linked)
- **Architecture**: Linux ARM64 (aarch64)

## 🐛 Robust Error Handling & Debugging Features

**v1.0.9 provides proper PATH environment in systemd service with all robust error handling features:**

### Core Reliability Features
- **Graceful Degradation**: System maintains RTC functionality even when scheduling fails
- **Intelligent Scheduler Selection**: Automatic fallback from `at` to `systemd-run` based on availability
- **Filesystem Read-Only Detection**: Prevents failures in embedded systems with read-only filesystems
- **Proper PATH Environment**: Service file includes `Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin` so systemd-run can find systemctl and system commands
- **Proper ARM64 Cross-Compilation**: Statically linked binary with `CGO_ENABLED=0` for Raspberry Pi compatibility
- **Detailed Error Classification**: Specific error types with actionable recovery instructions

### Enhanced Debugging
- **Detailed stderr logging** during service execution
- **Dependency validation** before service startup
- **Step-by-step execution tracing** to identify failure points
- **Scheduler status reporting** showing which scheduler is active and why
- **Clear error messages** for configuration, RTC, and scheduler issues

### Operational Modes
- **Full Mode**: RTC + AtScheduler (optimal performance)
- **Hybrid Mode**: RTC + SystemdTimerScheduler (read-only filesystem compatibility)
- **Degraded Mode**: RTC only (scheduler unavailable, but core functionality preserved)

To debug service issues, check systemd logs:
```bash
sudo journalctl -u rtc-scheduler -n 20 --no-pager

# Check scheduler status
rtc-scheduler -status

# Test manual scheduling
sudo rtc-scheduler -wake 08:00 -shutdown 22:00 -test
```

## 📄 License

This release is licensed under the MIT License. See the main repository for full license text.

## 🆘 Support

For issues or questions:
- GitHub Issues: https://github.com/EstebanJS/rtc-scheduler/issues
- Documentation: https://github.com/EstebanJS/rtc-scheduler#readme