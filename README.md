# Go System Monitor

A powerful terminal-based system monitoring tool built in Go that provides real-time metrics and alerts for your system.

![System Monitor Screenshot](screenshot.png)

## Features

- **Real-time system metrics**: CPU, memory, disk, network usage, and running processes
- **Interactive UI**: Tab-based navigation with keyboard shortcuts
- **Process management**: Sort processes by CPU usage, memory usage, PID, or name
- **Alert system**: Configurable alerts for resource usage thresholds
- **Visual indicators**: Color-coded metrics showing resource health
- **Clean interface**: Well-organized layout with proper separation of concerns

## Installation

### Prerequisites

- Go 1.15 or higher

### Install from source

```bash
# Clone the repository
git clone https://github.com/yourusername/go-system-monitor.git
cd go-system-monitor

# Build the executable
go build -o sysmon

# Make it executable
chmod +x sysmon

# Optional: move to a directory in your PATH
sudo mv sysmon /usr/local/bin/
```

## Usage

### Basic Usage

Simply run the executable:

```bash
./sysmon
```

Or if you moved it to your PATH:

```bash
sysmon
```

### Keyboard Controls

- **Tab / ← →**: Navigate between tabs
- **1-4**: When on Processes tab, sort by different criteria
  - **1**: Sort by CPU usage
  - **2**: Sort by memory usage
  - **3**: Sort by PID
  - **4**: Sort by name
- **r**: Manually refresh data
- **q / Esc**: Quit the application

## Configuration

By default, the application uses these alert thresholds:
- CPU Usage: 85%
- Memory Usage: 85%
- Disk Usage: 90%
- Swap Usage: 80%

You can customize these by editing the configuration file at `~/.config/sysmon/config.json`.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) for system metrics collection
- [bubbletea](https://github.com/charmbracelet/bubbletea) for the terminal UI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal styling
