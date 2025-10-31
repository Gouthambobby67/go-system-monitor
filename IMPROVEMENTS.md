# Codebase Analysis: Improvements & Missing Features

## 🔴 Critical Issues (Must Fix)

### 1. **Missing Time Series History Fields** ⚠️ HIGH PRIORITY
**Problem**: Code references `metrics.CPU.History` and `metrics.Memory.History` in `ui/combined_view.go` (lines 60, 69), but these fields don't exist in `CPUInfo` or `MemoryInfo` structs.

**Impact**: Sparklines won't render, causes runtime panic.

**Fix Required**:
```go
// In system/metrics.go, add to CPUInfo:
type CPUInfo struct {
    // ... existing fields
    History system.TimeSeries  // ADD THIS
}

// In system/metrics.go, add to MemoryInfo:
type MemoryInfo struct {
    // ... existing fields
    History system.TimeSeries  // ADD THIS
}

// Update Collect() to maintain history (keep last N points)
```

### 2. **Incomplete Function Implementation**
**Problem**: `collectCPUInfo()` in `system/metrics.go` line 226 is missing opening brace `{`.

**Impact**: Code won't compile.

### 3. **Missing cardConfig Field**
**Problem**: `ui/combined_view.go` uses `d.cardConfig` but `cardConfig` field doesn't exist in `Dashboard` struct.

**Impact**: Code won't compile, runtime panic when accessing combined view.

**Fix Required**:
```go
// In ui/dashboard.go, add to Dashboard struct:
type Dashboard struct {
    // ... existing fields
    cardConfig CardConfig  // ADD THIS
}

// In NewDashboard(), initialize:
cardConfig: DefaultCardConfig(),
```

### 4. **Duplicate/Dead Code**
**Problem**: `main_update.go` exists but seems redundant. Either consolidate or document purpose.

**Fix**: Remove or merge with `main.go`.

---

## 🟡 High Priority Improvements

### 5. **Network I/O Rates (Not Just Totals)**
**Current**: Shows cumulative bytes since boot.

**Needed**: Calculate and display bytes/second rates (like `iftop` or `nethogs`).

**Implementation**:
- Store previous counters
- Calculate delta between samples
- Display: `123.4 MB/s` instead of just totals

### 6. **Disk I/O Rates**
**Current**: Shows total read/write bytes only.

**Needed**: Read/write speeds (MB/s, IOPS).

### 7. **Process Table Scrolling**
**Current**: Process table is fixed, shows only top N processes.

**Needed**: 
- Vertical scrolling (↑↓ keys)
- Search/filter functionality
- Process details panel (on Enter/click)

### 8. **Process Details Missing**
**Current**: Only shows PID, CPU%, MEM%, NAME, USER, STATUS.

**Needed**:
- Command line arguments
- Thread count (currently shows "-")
- File descriptors count
- Parent PID
- Nice/Priority
- Start time
- Working directory
- Memory RSS vs VSS breakdown

---

## 🟢 Medium Priority Features

### 9. **Process Filtering & Search**
**Features**:
- Filter by name pattern (`/firefox`)
- Filter by user (`u:username`)
- Filter by CPU/Memory thresholds
- Real-time search as you type

**UI**: Add search bar at top of process view.

### 10. **Process Management Actions**
**Features**:
- Kill process (K key)
- Kill with signal selection (K, then SIGTERM/SIGKILL)
- Nice/renice process
- Only show own processes option

**Security**: Warn before killing, require confirmation for SIGKILL.

### 11. **Historical Data & Export**
**Features**:
- Save metrics snapshots to file
- Export to CSV/JSON
- Load historical data for comparison
- Chart view with time axis

**Implementation**:
```go
// New command: sysmon --export=csv --output=metrics.csv
// New tab: "History" showing past data
```

### 12. **Configurable Refresh Intervals**
**Current**: Fixed refresh rate from config.

**Needed**:
- Per-view refresh rates (slow down process collection)
- Pause refresh (`Space` key)
- Manual refresh (`R` key exists, good!)

### 13. **Better Error Handling & Logging**
**Current**: Some errors are silently ignored or just logged.

**Needed**:
- Error dialog/notification in UI
- Error log file (`~/.config/sysmon/error.log`)
- Graceful degradation (show partial data if some metrics fail)
- Warning badges in UI when metrics unavailable

### 14. **Mouse Support Implementation**
**Current**: Mouse is enabled (`tea.WithMouseCellMotion()`) but not used.

**Features**:
- Click tabs to switch
- Click column headers to sort
- Scroll process table with mouse wheel
- Click process row for details

---

## 🔵 Nice-to-Have Features

### 14. **Advanced Process View**
- Tree view (parent-child relationships)
- Group by process name
- Aggregate CPU/Memory by program
- Process timelines (when started/ended)

### 15. **Network Connection Details**
**Current**: Shows connection count only.

**Needed**:
- List active connections (like `netstat`)
- Filter by port/process
- Connection state details
- Bandwidth per connection

### 16. **Disk Details**
- I/O wait time
- Queue depth
- Per-partition I/O stats
- SMART health status (where available)

### 17. **System Load Breakdown**
- Per-core CPU graphs
- CPU frequency/voltage (where available)
- Top wait I/O processes

### 18. **Alert Improvements**
- Configurable alert actions (log, notify, command)
- Alert history with graphs
- Alert suppression rules
- Custom alert thresholds per resource

### 19. **Theme Customization**
**Current**: Hard-coded dark theme.

**Features**:
- Multiple built-in themes (dark, light, solarized)
- Custom theme configuration
- Color blindness mode
- Export/import themes

### 20. **Multi-Host Monitoring**
- Connect to remote hosts via SSH
- Monitor multiple systems in tabs
- Aggregate metrics across hosts

### 21. **Performance Optimizations**
**Issues**:
- Collecting all processes every second is expensive
- Consider caching process info
- Lazy-load detailed process data
- Use goroutines for parallel metric collection

### 22. **Better Terminal Responsiveness**
- Handle terminal resize better
- Graceful exit on SIGTERM/SIGINT
- Background mode option
- Auto-reconnect if terminal restored

---

## 🛠️ Code Quality Improvements

### 23. **Testing**
**Current**: No tests found.

**Needed**:
- Unit tests for metrics collection
- Mock tests for UI rendering
- Integration tests for full flow
- Benchmark tests for performance

### 24. **Documentation**
**Needed**:
- API documentation (godoc)
- Architecture diagrams
- Contributing guide
- Code comments explaining complex logic

### 25. **Code Organization**
**Issues**:
- Some functions are very long (e.g., `collectProcessInfo`)
- Duplicate code between views
- Build tags on some files (`//go:build ignore`) - why?

**Suggestions**:
- Break large functions into smaller ones
- Extract common UI patterns
- Document build tag usage

### 26. **Error Messages**
**Current**: Generic error messages.

**Needed**:
- User-friendly error messages
- Suggestions for fixing errors
- Help text for common issues

### 27. **Platform-Specific Features**
**Current**: Cross-platform, but some features limited.

**Enhancements**:
- Linux: cgroup support, systemd status
- macOS: Spotlight indexer stats, Time Machine status
- Windows: Service status, Windows Update status

---

## 📋 Summary Priority List

### Must Fix (Before Next Release)
1. ✅ Add History fields to CPUInfo/MemoryInfo
2. ✅ Add cardConfig field to Dashboard struct
3. ✅ Remove/consolidate main_update.go
4. ✅ Implement history collection in Collect()

### Should Have (Next Version)
5. Network/disk I/O rates
6. Process scrolling
7. Process details panel
8. Thread count implementation
9. Mouse support implementation

### Nice to Have (Future Versions)
10. Process filtering/search
11. Export functionality
12. Process management (kill/renice)
13. Theme customization
14. Multi-host monitoring

---

## 🎯 Quick Wins (Easy Improvements)

1. **Add thread count** - Just call `p.NumThreads()` in collectProcessInfo
2. **Better formatting** - Add thousands separators to numbers
3. **Keyboard shortcuts** - Add `?` for help (already exists, just document)
4. **Status indicators** - Add visual indicators for metric health
5. **Process count display** - Show "15 of 234 processes" in header

---

## 📝 Implementation Notes

### For Time Series History:
```go
// In Collector struct, add:
cpuHistory    []system.TimeSeriesPoint
memoryHistory []system.TimeSeriesPoint
maxHistoryPoints int

// In Collect():
// Add current point to history
// Trim if exceeds maxHistoryPoints
// Assign to CPUInfo.History and MemoryInfo.History
```

### For I/O Rates:
```go
// Store previous counters
var prevNetCounters map[string]net.IOCountersStat
var prevDiskCounters map[string]disk.IOCountersStat

// Calculate delta and divide by time interval
```

### For Process Scrolling:
```go
// Add to ProcessTable:
offset int  // current scroll position

// Handle ↑↓ keys in main Update()
// Update offset based on keypress
```

---

Generated from codebase analysis on 2025-01-XX

