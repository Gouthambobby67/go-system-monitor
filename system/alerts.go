package system

import (
	"fmt"
	"time"
)

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	// Alert levels
	InfoLevel    AlertLevel = "info"
	WarningLevel AlertLevel = "warning"
	CriticalLevel AlertLevel = "critical"
)

// Alert represents a system alert for high resource usage
type Alert struct {
	Timestamp time.Time
	Message   string
	Level     AlertLevel
	Source    string
	Resolved  bool
}

// AlertManager handles system alerts
type AlertManager struct {
	Alerts       []Alert
	MaxAlerts    int
	CPUThreshold float64
	MemThreshold float64
	DiskThreshold float64
	SwapThreshold float64
}

// NewAlertManager creates a new alert manager with configurable thresholds
func NewAlertManager(cpuThreshold, memThreshold, diskThreshold, swapThreshold float64, maxAlerts int) *AlertManager {
	return &AlertManager{
		Alerts:       []Alert{},
		MaxAlerts:    maxAlerts,
		CPUThreshold: cpuThreshold,
		MemThreshold: memThreshold,
		DiskThreshold: diskThreshold,
		SwapThreshold: swapThreshold,
	}
}

// AddAlert adds a new alert
func (am *AlertManager) AddAlert(message string, level AlertLevel, source string) {
	// Check if a similar unresolved alert already exists
	for i, alert := range am.Alerts {
		if alert.Source == source && alert.Level == level && !alert.Resolved {
			// Similar alert exists, don't create a duplicate
			return
		}
		
		// If we find a resolved alert with the same source and it's a different level,
		// we can mark it as resolved since the condition has changed
		if alert.Source == source && alert.Level != level && !alert.Resolved {
			am.Alerts[i].Resolved = true
		}
	}

	// Add new alert to the beginning of the slice (most recent first)
	am.Alerts = append([]Alert{{
		Timestamp: time.Now(),
		Message:   message,
		Level:     level,
		Source:    source,
		Resolved:  false,
	}}, am.Alerts...)

	// Trim old alerts if we exceed the maximum
	if len(am.Alerts) > am.MaxAlerts {
		am.Alerts = am.Alerts[:am.MaxAlerts]
	}
}

// ResolveAlert marks an alert as resolved
func (am *AlertManager) ResolveAlert(source string) {
	for i, alert := range am.Alerts {
		if alert.Source == source && !alert.Resolved {
			am.Alerts[i].Resolved = true
			// Add a resolution notice
			am.AddAlert(
				fmt.Sprintf("%s has returned to normal levels", source),
				InfoLevel,
				fmt.Sprintf("%s_resolved", source),
			)
			break
		}
	}
}

// CheckResourceAlerts generates alerts based on resource usage
func (am *AlertManager) CheckResourceAlerts(metrics *Collector) {
	// Check CPU usage
	if metrics.CPU.Usage >= am.CPUThreshold {
		am.AddAlert(
			fmt.Sprintf("CPU usage is high (%.1f%%)", metrics.CPU.Usage),
			CriticalLevel,
			"cpu_usage",
		)
	} else if metrics.CPU.Usage < am.CPUThreshold-10 { // 10% hysteresis
		am.ResolveAlert("cpu_usage")
	}

	// Check memory usage
	if metrics.Memory.UsedPercent >= am.MemThreshold {
		am.AddAlert(
			fmt.Sprintf("Memory usage is high (%.1f%%)", metrics.Memory.UsedPercent),
			CriticalLevel,
			"memory_usage",
		)
	} else if metrics.Memory.UsedPercent < am.MemThreshold-10 { // 10% hysteresis
		am.ResolveAlert("memory_usage")
	}

	// Check swap usage if swap is enabled
	if metrics.Memory.SwapTotal > 0 && metrics.Memory.SwapPercent >= am.SwapThreshold {
		am.AddAlert(
			fmt.Sprintf("Swap usage is high (%.1f%%)", metrics.Memory.SwapPercent),
			WarningLevel,
			"swap_usage",
		)
	} else if metrics.Memory.SwapTotal > 0 && metrics.Memory.SwapPercent < am.SwapThreshold-10 {
		am.ResolveAlert("swap_usage")
	}

	// Check disk usage
	for mountpoint, usage := range metrics.Disk.UsageStats {
		if usage.UsedPercent >= am.DiskThreshold {
			am.AddAlert(
				fmt.Sprintf("Disk usage on %s is high (%.1f%%)", mountpoint, usage.UsedPercent),
				WarningLevel,
				fmt.Sprintf("disk_usage_%s", mountpoint),
			)
		} else if usage.UsedPercent < am.DiskThreshold-5 { // 5% hysteresis for disk
			am.ResolveAlert(fmt.Sprintf("disk_usage_%s", mountpoint))
		}
	}
}
