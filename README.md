<img width="1149" height="413" alt="Screenshot_20251031_144205" src="https://github.com/user-attachments/assets/1f63e97c-2352-484a-af0a-b2e121e562cb" /># Go System Monitor

A powerful terminal-based system monitoring tool built in Go that provides real-time metrics and alerts for your system.

![System Monitor Screenshot](screenshot.png)

![Uploading Screenshot_20251031_144205.png<img width="1149" height="413" alt="Screenshot_20251031_144205-1" src="https://github.com/user-attachments/assets/188083cb-db3e-489b-9f37-8a4bb56f1fcc" /><img width="1148" height="401" alt="Screenshot_20251031_144228-1" src="https://github.com/user-attachments/assets/da729c18-f9b1-4d41-b2e0-b0c6f1b5fc81" />
<img width="1148" height="401" alt="Screenshot_20251031_144228" src="https://github.com/user-attachments/assets/f40a925a-36d8-4536-8694-6804636d66b0" />


## Features

### Core Monitoring
- **Real-time system metrics**: CPU, memory, disk, network usage, and running processes
- **Historical data tracking**: Sparklines showing CPU and memory trends over time
- **Real-time I/O rates**: Network and disk speeds in MB/s (not just cumulative totals)
- **Comprehensive process details**: Thread count, command line, parent PID, memory breakdown
- **Interactive UI**: Tab-based navigation with keyboard shortcuts

### Advanced Features
- **Process scrolling**: Navigate through hundreds of processes with keyboard/mouse
- **Process filtering**: Search and filter processes by name, user, or command
- **Process sorting**: Sort by CPU, memory, PID, or name
- **Alert system**: Configurable alerts for resource usage thresholds
- **Mouse support**: Click tabs, scroll with mouse wheel
- **Visual indicators**: Color-coded metrics showing resource health
- **Number formatting**: Thousands separators for better readability

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

#### General Navigation
- **Tab / ← →**: Navigate between tabs
- **r**: Manually refresh data
- **c**: Toggle compact mode
- **f**: Toggle fullscreen
- **s**: Toggle status bar
- **?**: Show/hide help overlay
- **q / Esc**: Quit the application

#### Process Table (Processes Tab)
**Scrolling**:
- **↑ / k**: Scroll up one line
- **↓ / j**: Scroll down one line
- **PgUp / Ctrl+u**: Scroll up one page
- **PgDn / Ctrl+d**: Scroll down one page
- **Home / g**: Jump to top
- **End / G**: Jump to bottom

**Sorting**:
- **1**: Sort by CPU usage
- **2**: Sort by memory usage
- **3**: Sort by PID
- **4**: Sort by name

### Mouse Controls
- **Click tabs**: Switch between different views
- **Mouse wheel**: Scroll through process list
- **Click anywhere**: Basic interaction support

## Configuration

### Default Alert Thresholds
- CPU Usage: 85%
- Memory Usage: 85%
- Disk Usage: 90%
- Swap Usage: 80%

### Configuration File
Location: `~/.config/sysmon/config.json`

Example configuration:
```json
{
  "cpu_threshold": 85.0,
  "memory_threshold": 85.0,
  "disk_threshold": 90.0,
  "swap_threshold": 80.0,
  "refresh_interval_ms": 1000,
  "max_processes": 15,
  "max_alerts_to_keep": 100,
  "default_sorting_mode": "cpu"
}
```

### Command Line Options
```bash
./sysmon -cpu 90        # Set CPU threshold
./sysmon -mem 85        # Set memory threshold
./sysmon -disk 95       # Set disk threshold
./sysmon -swap 80       # Set swap threshold
./sysmon -version       # Show version
./sysmon -help          # Show help
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## What's New in v2.0

### Major Improvements
- ✅ **Real-time I/O rates** - See actual network/disk speeds, not cumulative totals
- ✅ **Process scrolling** - Navigate through all processes, not just top 15
- ✅ **Historical data** - Sparklines showing CPU/memory trends
- ✅ **Complete process info** - Thread count, command line, parent PID, and more
- ✅ **Process filtering** - Search and filter processes easily
- ✅ **Mouse support** - Click tabs and scroll with mouse
- ✅ **Better error handling** - Graceful degradation, no crashes
- ✅ **Enhanced UI** - Number formatting, better help, visual indicators

### Bug Fixes
- Fixed missing history fields (sparklines now work)
- Fixed missing cardConfig (combined view works)
- Removed duplicate code
- Fixed thread count display
- Improved error resilience

See [CHANGELOG.md](CHANGELOG.md) for complete details.

## Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) for system metrics collection
- [bubbletea](https://github.com/charmbracelet/bubbletea) for the terminal UI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) for terminal styling
- [bubbles](https://github.com/charmbracelet/bubbles) for TUI components
