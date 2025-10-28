package ui

import "go_system_monitor/system"

// CardConfig controls which cards are shown in the combined view
type CardConfig struct {
	ShowSystem    bool
	ShowCPU      bool
	ShowMemory   bool
	ShowDisk     bool
	ShowNetwork  bool
	ShowAlerts   bool
	ShowSparkline bool
}

// DefaultCardConfig returns the default card configuration
func DefaultCardConfig() CardConfig {
	return CardConfig{
		ShowSystem:    true,
		ShowCPU:      true,
		ShowMemory:   true,
		ShowDisk:     true,
		ShowNetwork:  true,
		ShowAlerts:   true,
		ShowSparkline: true,
	}
}

// CardMetrics holds compact metric data for badges
type CardMetrics struct {
	Label    string
	Value    float64
	Unit     string
	Alert    bool
	Critical bool
}

// CompactSystemMetrics extracts key metrics for badges
func CompactSystemMetrics(info *system.Collector) []CardMetrics {
	metrics := make([]CardMetrics, 0, 6)

	// CPU usage
	metrics = append(metrics, CardMetrics{
		Label:    "CPU",
		Value:    info.CPU.Usage,
		Unit:     "%",
		Alert:    info.CPU.Usage >= 60,
		Critical: info.CPU.Usage >= 85,
	})

	// Memory usage
	metrics = append(metrics, CardMetrics{
		Label:    "MEM",
		Value:    info.Memory.UsedPercent,
		Unit:     "%",
		Alert:    info.Memory.UsedPercent >= 60,
		Critical: info.Memory.UsedPercent >= 85,
	})

	// Disk usage (root partition)
	if usage, ok := info.Disk.UsageStats["/"]; ok {
		metrics = append(metrics, CardMetrics{
			Label:    "DISK",
			Value:    usage.UsedPercent,
			Unit:     "%",
			Alert:    usage.UsedPercent >= 80,
			Critical: usage.UsedPercent >= 90,
		})
	}

	// Load average (1min)
	if info.CPU.LoadAvg != nil {
		metrics = append(metrics, CardMetrics{
			Label:    "LOAD",
			Value:    info.CPU.LoadAvg.Load1,
			Unit:     "",
			Alert:    info.CPU.LoadAvg.Load1 > float64(info.CPU.Cores),
			Critical: info.CPU.LoadAvg.Load1 > float64(info.CPU.Cores*2),
		})
	}

	return metrics
}