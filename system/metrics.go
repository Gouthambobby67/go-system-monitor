package system

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemInfo contains all system information
type SystemInfo struct {
	Hostname    string
	Platform    string
	OS          string
	KernelVer   string
	Uptime      time.Duration
	LastUpdated time.Time
}

// CPUInfo contains CPU metrics
type CPUInfo struct {
	Usage       float64
	UsagePerCPU []float64
	Cores       int
	LoadAvg     *load.AvgStat
	Temperature float64 // Note: Might not be available on all systems
}

// MemoryInfo contains memory metrics
type MemoryInfo struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapFree    uint64
	SwapPercent float64
}

// DiskInfo contains disk metrics
type DiskInfo struct {
	Partitions  []disk.PartitionStat
	UsageStats  map[string]*disk.UsageStat
	IOCounters  map[string]disk.IOCountersStat
}

// NetworkInfo contains network metrics
type NetworkInfo struct {
	Interfaces  []net.InterfaceStat
	IOCounters  map[string]net.IOCountersStat
	Connections []net.ConnectionStat
}

// SortType defines process sorting methods
type SortType string

// Process sort types
const (
	SortByCPU     SortType = "cpu"
	SortByMemory  SortType = "memory"
	SortByPID     SortType = "pid"
	SortByName    SortType = "name"
)

// ProcessInfo contains process metrics
type ProcessInfo struct {
	Processes []ProcessDetail
	Total     int
	SortBy    SortType
}

// ProcessDetail contains details for a single process
type ProcessDetail struct {
	PID        int32
	Name       string
	Username   string
	Status     []string
	CPUPercent float64
	MemPercent float32
	CreatedAt  time.Time
}

// Collector handles collecting and storing metrics
type Collector struct {
	System      SystemInfo
	CPU         CPUInfo
	Memory      MemoryInfo
	Disk        DiskInfo
	Network     NetworkInfo
	Process     ProcessInfo
	Interval    time.Duration
	AlertManager *AlertManager
	// MaxProcesses limits how many processes are shown in the UI
	MaxProcesses int
}

// NewCollector creates a new metrics collector with optional configuration
func NewCollector(cpuThreshold, memThreshold, diskThreshold, swapThreshold float64, refreshMs int, sortMode string, maxProcesses int, maxAlerts int) *Collector {
	// Convert milliseconds to time.Duration
	interval := time.Duration(refreshMs) * time.Millisecond
	
	// Convert string sort mode to SortType
	sortBy := SortByCPU
	switch sortMode {
	case "memory":
		sortBy = SortByMemory
	case "pid":
		sortBy = SortByPID
	case "name":
		sortBy = SortByName
	}
	
	return &Collector{
		Interval: interval,
		Process: ProcessInfo{
			SortBy: sortBy,
		},
		AlertManager: NewAlertManager(cpuThreshold, memThreshold, diskThreshold, swapThreshold, maxAlerts),
		MaxProcesses: maxProcesses,
	}
}

// SortProcesses sorts the processes according to the specified sort type
func (c *Collector) SortProcesses(processes []ProcessDetail) {
	switch c.Process.SortBy {
	case SortByCPU:
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
		
	case SortByMemory:
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].MemPercent > processes[j].MemPercent
		})
		
	case SortByPID:
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].PID < processes[j].PID
		})
		
	case SortByName:
		sort.Slice(processes, func(i, j int) bool {
			return strings.ToLower(processes[i].Name) < strings.ToLower(processes[j].Name)
		})
		
	default:
		// Default to CPU sorting
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
	}
}

// Collect gathers all system metrics
func (c *Collector) Collect() error {
	var err error

	// Update timestamp
	c.System.LastUpdated = time.Now()

	// Collect system info
	if err = c.collectSystemInfo(); err != nil {
		return fmt.Errorf("failed to collect system info: %v", err)
	}

	// Collect CPU info
	if err = c.collectCPUInfo(); err != nil {
		return fmt.Errorf("failed to collect CPU info: %v", err)
	}

	// Collect memory info
	if err = c.collectMemoryInfo(); err != nil {
		return fmt.Errorf("failed to collect memory info: %v", err)
	}

	// Collect disk info
	if err = c.collectDiskInfo(); err != nil {
		return fmt.Errorf("failed to collect disk info: %v", err)
	}

	// Collect network info
	if err = c.collectNetworkInfo(); err != nil {
		return fmt.Errorf("failed to collect network info: %v", err)
	}

	// Collect process info
	if err = c.collectProcessInfo(); err != nil {
		return fmt.Errorf("failed to collect process info: %v", err)
	}

	// Check for any alerts based on collected metrics
	if c.AlertManager != nil {
		c.AlertManager.CheckResourceAlerts(c)
	}

	return nil
}

// collectSystemInfo gathers system information
func (c *Collector) collectSystemInfo() error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	c.System.Hostname = info.Hostname
	c.System.Platform = info.Platform
	c.System.OS = info.OS
	c.System.KernelVer = info.KernelVersion
	c.System.Uptime = time.Duration(info.Uptime) * time.Second

	return nil
}

// collectCPUInfo gathers CPU metrics
func (c *Collector) collectCPUInfo() error {
	// Get CPU usage (overall)
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	if len(percentages) > 0 {
		c.CPU.Usage = percentages[0]
	}

	// Get per-CPU usage
	perCPU, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}
	c.CPU.UsagePerCPU = perCPU
	c.CPU.Cores = len(perCPU)

	// Get load average
	loadAvg, err := load.Avg()
	if err != nil {
		// Not critical, just log and continue
		log.Printf("Warning: Could not get load average: %v", err)
	} else {
		c.CPU.LoadAvg = loadAvg
	}

	// Try to get temperature (might not work on all systems)
	// This is a simplified approach - real implementation might need to be platform-specific
	temps, err := host.SensorsTemperatures()
	if err == nil {
		for _, temp := range temps {
			if temp.SensorKey == "coretemp_packageid0_input" || 
			   temp.SensorKey == "k10temp_tdie" ||
			   temp.SensorKey == "cpu_thermal_input" {
				c.CPU.Temperature = temp.Temperature
				break
			}
		}
	}

	return nil
}

// collectMemoryInfo gathers memory metrics
func (c *Collector) collectMemoryInfo() error {
	// Get virtual memory stats
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	c.Memory.Total = vmem.Total
	c.Memory.Used = vmem.Used
	c.Memory.Free = vmem.Free
	c.Memory.UsedPercent = vmem.UsedPercent

	// Get swap memory stats
	swap, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	c.Memory.SwapTotal = swap.Total
	c.Memory.SwapUsed = swap.Used
	c.Memory.SwapFree = swap.Free
	c.Memory.SwapPercent = swap.UsedPercent

	return nil
}

// collectDiskInfo gathers disk metrics
func (c *Collector) collectDiskInfo() error {
	// Get partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return err
	}
	c.Disk.Partitions = partitions

	// Get usage for each partition
	usageStats := make(map[string]*disk.UsageStat)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Printf("Warning: Could not get usage for %s: %v", partition.Mountpoint, err)
			continue
		}
		usageStats[partition.Mountpoint] = usage
	}
	c.Disk.UsageStats = usageStats

	// Get IO counters
	ioCounters, err := disk.IOCounters()
	if err != nil {
		log.Printf("Warning: Could not get disk IO counters: %v", err)
	} else {
		c.Disk.IOCounters = ioCounters
	}

	return nil
}

// collectNetworkInfo gathers network metrics
func (c *Collector) collectNetworkInfo() error {
	// Get network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	c.Network.Interfaces = interfaces

	// Get network IO counters
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		log.Printf("Warning: Could not get network IO counters: %v", err)
	} else {
		countersMap := make(map[string]net.IOCountersStat)
		for _, ioc := range ioCounters {
			countersMap[ioc.Name] = ioc
		}
		c.Network.IOCounters = countersMap
	}

	// Get network connections (might require elevated privileges)
	connections, err := net.Connections("all")
	if err != nil {
		log.Printf("Warning: Could not get network connections: %v", err)
	} else {
		c.Network.Connections = connections
	}

	return nil
}

// collectProcessInfo gathers process metrics
func (c *Collector) collectProcessInfo() error {
	// Get all processes
	pids, err := process.Pids()
	if err != nil {
		return err
	}

	processes := make([]ProcessDetail, 0, len(pids))
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue // Skip this process
		}

		// Get process name
		name, err := p.Name()
		if err != nil {
			name = "unknown"
		}

		// Get process username
		username, err := p.Username()
		if err != nil {
			username = "unknown"
		}

		// Get process status
		status, err := p.Status()
		if err != nil {
			status = []string{"unknown"}
		}

		// Get CPU usage
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		// Get memory usage
		memPercent, err := p.MemoryPercent()
		if err != nil {
			memPercent = 0
		}

		// Get creation time
		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		// CreateTime returns milliseconds since epoch; use UnixMilli for clarity
		createdAt := time.UnixMilli(createTime)

		processes = append(processes, ProcessDetail{
			PID:        pid,
			Name:       name,
			Username:   username,
			Status:     status,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			CreatedAt:  createdAt,
		})
	}

	// Sort the processes according to the sort type
	c.SortProcesses(processes)

	c.Process.Total = len(processes)

	// Store processes; keep full list but UI will respect MaxProcesses when rendering
	c.Process.Processes = processes

	return nil
}
