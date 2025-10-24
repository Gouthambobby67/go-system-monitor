package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// SystemView renders the System tab
func (d *Dashboard) SystemView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	// Build system information section
	sysInfo := fmt.Sprintf("Hostname: %s\n", metrics.System.Hostname)
	sysInfo += fmt.Sprintf("OS: %s %s\n", metrics.System.Platform, metrics.System.OS)
	sysInfo += fmt.Sprintf("Kernel: %s\n", metrics.System.KernelVer)
	sysInfo += fmt.Sprintf("Uptime: %s\n", formatDuration(metrics.System.Uptime))
	
	// Add summary of key metrics
	sysInfo += "\n--- System Summary ---\n"
	sysInfo += fmt.Sprintf("CPU Usage: %s\n", formatValue(metrics.CPU.Usage))
	sysInfo += fmt.Sprintf("Memory Usage: %s\n", formatValue(metrics.Memory.UsedPercent))
	if len(metrics.Disk.Partitions) > 0 {
		for mp, usage := range metrics.Disk.UsageStats {
			if mp == "/" || mp == "/home" {
				sysInfo += fmt.Sprintf("Disk Usage (%s): %s\n", mp, formatValue(usage.UsedPercent))
			}
		}
	}
	
	// Check if there are any active alerts
	var activeAlertCount int
	for _, alert := range metrics.AlertManager.Alerts {
		if !alert.Resolved {
			activeAlertCount++
		}
	}
	
	// Show alert status
	if activeAlertCount > 0 {
		sysInfo += "\n"
		alertMsg := fmt.Sprintf("⚠️  %d active alerts - check Alerts tab", activeAlertCount)
		sysInfo += criticalValueStyle.Render(alertMsg) + "\n"
	}
	
	// Last updated timestamp
	sysInfo += fmt.Sprintf("\nLast Updated: %s", metrics.System.LastUpdated.Format("15:04:05"))
	
	return infoSectionStyle.Width(width).Render(sysInfo)
}

// CPUView renders the CPU tab
func (d *Dashboard) CPUView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	// Overall CPU usage
	content := RenderProgress("CPU Usage", metrics.CPU.Usage, width)
	
	// Per-core CPU usage
	content += "\n\n--- Per CPU Core Usage ---\n"
	for i, usage := range metrics.CPU.UsagePerCPU {
		content += RenderProgress(fmt.Sprintf("Core %d", i), usage, width) + "\n"
	}
	
	// Load average if available
	if metrics.CPU.LoadAvg != nil {
		content += fmt.Sprintf("\n--- Load Average ---\n")
		content += fmt.Sprintf("1 min: %.2f  5 min: %.2f  15 min: %.2f\n", 
			metrics.CPU.LoadAvg.Load1,
			metrics.CPU.LoadAvg.Load5,
			metrics.CPU.LoadAvg.Load15)
	}

	// CPU temperature section (always shown)
	content += fmt.Sprintf("\n--- CPU Temperature ---\n")
	if metrics.CPU.Temperature <= 0 {
		content += "Temperature: N/A\n"
	} else {
		// Real-time temperature gauge
		content += RenderProgress("Temp", metrics.CPU.Temperature, width) + "\n"
		tempStyle := normalValueStyle
		if metrics.CPU.Temperature > 70 {
			tempStyle = warnValueStyle
		}
		if metrics.CPU.Temperature > 85 {
			tempStyle = criticalValueStyle
		}
		content += fmt.Sprintf("Temperature: %s\n", tempStyle.Render(fmt.Sprintf("%.1f°C", metrics.CPU.Temperature)))
	}

	return infoSectionStyle.Width(width).Render(content)
}

// MemoryView renders the Memory tab
func (d *Dashboard) MemoryView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	// RAM usage
	content := "--- Physical Memory ---\n"
	content += RenderProgress("Memory", metrics.Memory.UsedPercent, width) + "\n"
	content += fmt.Sprintf(
		"Total: %s  Used: %s  Free: %s\n",
		FormatBytes(metrics.Memory.Total),
		FormatBytes(metrics.Memory.Used),
		FormatBytes(metrics.Memory.Free))
	
	// Swap usage
	content += "\n--- Swap ---\n"
	content += RenderProgress("Swap", metrics.Memory.SwapPercent, width) + "\n"
	content += fmt.Sprintf(
		"Total: %s  Used: %s  Free: %s\n",
		FormatBytes(metrics.Memory.SwapTotal),
		FormatBytes(metrics.Memory.SwapUsed),
		FormatBytes(metrics.Memory.SwapFree))
	
	return infoSectionStyle.Width(width).Render(content)
}

// DiskView renders the Disk tab
func (d *Dashboard) DiskView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	content := "--- Storage Usage ---\n"
	
	// Sort partitions by mount point for consistent display
	partitions := metrics.Disk.Partitions
	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].Mountpoint < partitions[j].Mountpoint
	})
	
	// Show usage for each partition
	for _, partition := range partitions {
		mp := partition.Mountpoint
		if usage, ok := metrics.Disk.UsageStats[mp]; ok {
			// Skip some virtual filesystems
			if strings.HasPrefix(mp, "/sys") || strings.HasPrefix(mp, "/proc") ||
			   strings.HasPrefix(mp, "/dev") || strings.HasPrefix(mp, "/run") {
				continue
			}
			
			label := fmt.Sprintf("%-15s", mp)
			content += RenderProgress(label, usage.UsedPercent, width) + "\n"
			content += fmt.Sprintf("  %s / %s (%s)\n\n",
				FormatBytes(usage.Used),
				FormatBytes(usage.Total),
				partition.Fstype)
		}
	}
	
	// IO statistics
	if len(metrics.Disk.IOCounters) > 0 {
		content += "\n--- Disk I/O ---\n"
		for name, io := range metrics.Disk.IOCounters {
			content += fmt.Sprintf("%s:\n", name)
			content += fmt.Sprintf("  Read:  %s  Write: %s\n",
				FormatBytes(io.ReadBytes),
				FormatBytes(io.WriteBytes))
		}
	}
	
	return infoSectionStyle.Width(width).Render(content)
}

// NetworkView renders the Network tab
func (d *Dashboard) NetworkView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	content := "--- Network Interfaces ---\n"
	
	for _, iface := range metrics.Network.Interfaces {
		// Skip loopback interface
		if iface.Name == "lo" {
			continue
		}
		
		content += fmt.Sprintf("%s (%s):\n", iface.Name, iface.HardwareAddr)
		
		// Show IP addresses
		for _, addr := range iface.Addrs {
			content += fmt.Sprintf("  Address: %s\n", addr.Addr)
		}
		
		// Show IO stats if available
		if io, ok := metrics.Network.IOCounters[iface.Name]; ok {
			content += fmt.Sprintf("  Received: %s  Sent: %s\n\n",
				FormatBytes(io.BytesRecv),
				FormatBytes(io.BytesSent))
		} else {
			content += "\n"
		}
	}
	
	// Connection count
	if len(metrics.Network.Connections) > 0 {
		// Count connections by status
		statusCount := make(map[string]int)
		for _, conn := range metrics.Network.Connections {
			statusCount[conn.Status]++
		}
		
		content += "--- Network Connections ---\n"
		for status, count := range statusCount {
			content += fmt.Sprintf("%s: %d\n", status, count)
		}
		content += fmt.Sprintf("Total: %d\n", len(metrics.Network.Connections))
	}
	
	return infoSectionStyle.Width(width).Render(content)
}

// ProcessView renders the Processes tab
func (d *Dashboard) ProcessView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	// Show sort method in the header
	sortMethodText := fmt.Sprintf("Sorted by: %s", getSortMethodName(metrics.Process.SortBy))
	sortMethodStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F9FA")).
		Background(lipgloss.Color("#4361EE")).
		Padding(0, 1).
		Bold(true)
	
	headerText := fmt.Sprintf("--- Processes (Total: %d) ---", metrics.Process.Total)
	content := lipgloss.JoinHorizontal(lipgloss.Center, 
		headerText, 
		"   ", // spacing
		sortMethodStyle.Render(sortMethodText),
	) + "\n\n"
	
	// Add keyboard shortcuts for sorting
	content += "Sort: [1] CPU  [2] Memory  [3] PID  [4] Name\n"
	
	// Display only the top processes to avoid overwhelming the UI
	maxProcesses := metrics.MaxProcesses
	if maxProcesses <= 0 {
		maxProcesses = 15
	}
	if len(metrics.Process.Processes) < maxProcesses {
		maxProcesses = len(metrics.Process.Processes)
	}
	
	// Build the header with highlighted columns based on sort
	headerFmt := "%-7s %-6s %-20s %-10s %-8s %-8s\n"
	headerPID := "PID"
	headerCPU := "CPU%"
	headerName := "NAME"
	headerUser := "USER"
	headerMem := "MEM%"
	headerStatus := "STATUS"
	
	// Highlight the sorted column
	switch metrics.Process.SortBy {
	case system.SortByCPU:
		headerCPU = normalValueStyle.Render(headerCPU)
	case system.SortByMemory:
		headerMem = normalValueStyle.Render(headerMem)
	case system.SortByPID:
		headerPID = normalValueStyle.Render(headerPID)
	case system.SortByName:
		headerName = normalValueStyle.Render(headerName)
	}
	
	content += fmt.Sprintf(headerFmt, headerPID, headerCPU, headerName, headerUser, headerMem, headerStatus)
	content += strings.Repeat("-", width) + "\n"
	
	// Build table of processes
	rowFmt := "%-7d %-6.1f %-20s %-10s %-8.1f %-8s\n"
	for i := 0; i < maxProcesses; i++ {
		p := metrics.Process.Processes[i]
		name := p.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}
		
		username := p.Username
		if len(username) > 10 {
			username = username[:7] + "..."
		}
		
		// Format status from slice to string
		statusStr := "unknown"
		if len(p.Status) > 0 {
			statusStr = p.Status[0]
		}
		
		// Use lipgloss to style the row based on CPU usage
		cpuStyle := lipgloss.NewStyle()
		if p.CPUPercent > 50 {
			cpuStyle = cpuStyle.Foreground(lipgloss.Color("#FCC419")) // Yellow for high usage
		} else if p.CPUPercent > 80 {
			cpuStyle = cpuStyle.Foreground(lipgloss.Color("#FA5252")) // Red for very high usage
		}
		
		content += fmt.Sprintf(rowFmt,
			p.PID,
			p.CPUPercent,
			name,
			username,
			p.MemPercent,
			statusStr)
	}
	
	return infoSectionStyle.Width(width).Render(content)
}

// RenderHelp shows the help section
func (d *Dashboard) RenderHelp() string {
	help := "Tab/← →: Navigate tabs • 1-4: Sort processes • q: Quit • r: Refresh"
	return helpStyle.Render(help)
}

// getSortMethodName returns a user-friendly name for the sort method
func getSortMethodName(sortBy system.SortType) string {
	switch sortBy {
	case system.SortByCPU:
		return "CPU Usage"
	case system.SortByMemory:
		return "Memory Usage"
	case system.SortByPID:
		return "Process ID"
	case system.SortByName:
		return "Name"
	default:
		return "CPU Usage"
	}
}

// AlertsView renders the Alerts tab
func (d *Dashboard) AlertsView(metrics *system.Collector) string {
	width := d.width - 4 // Account for padding and borders
	
	// Add header and alert summary
	alertSummary := fmt.Sprintf("--- System Alerts ---")
	content := alertSummary + "\n\n"
	
	// Add current thresholds section
	content += "Current Alert Thresholds:\n"
	content += fmt.Sprintf("CPU Usage: %.1f%%  Memory Usage: %.1f%%  Disk Usage: %.1f%%  Swap Usage: %.1f%%\n\n", 
		metrics.AlertManager.CPUThreshold,
		metrics.AlertManager.MemThreshold,
		metrics.AlertManager.DiskThreshold,
		metrics.AlertManager.SwapThreshold)
	
	// If no alerts, show a message
	if len(metrics.AlertManager.Alerts) == 0 {
		content += "\nNo alerts at this time. System is running normally.\n"
		return infoSectionStyle.Width(width).Render(content)
	}
	
	// Group alerts by status (active vs resolved)
	var activeAlerts, resolvedAlerts []system.Alert
	for _, alert := range metrics.AlertManager.Alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		} else {
			resolvedAlerts = append(resolvedAlerts, alert)
		}
	}
	
	// Show active alerts
	if len(activeAlerts) > 0 {
		content += "Active Alerts:\n"
		content += strings.Repeat("-", width) + "\n"
		
		for _, alert := range activeAlerts {
			alertStyle := normalValueStyle
			if alert.Level == system.WarningLevel {
				alertStyle = warnValueStyle
			} else if alert.Level == system.CriticalLevel {
				alertStyle = criticalValueStyle
			}
			
			content += fmt.Sprintf("%s - %s\n", 
				alert.Timestamp.Format("15:04:05"),
				alertStyle.Render(alert.Message))
		}
		content += "\n"
	}
	
	// Show resolved alerts (with limit)
	if len(resolvedAlerts) > 0 {
		maxResolvedToShow := 5 // Limit the number of resolved alerts shown
		if len(resolvedAlerts) < maxResolvedToShow {
			maxResolvedToShow = len(resolvedAlerts)
		}
		
		content += "Recently Resolved:\n"
		content += strings.Repeat("-", width) + "\n"
		
		for i := 0; i < maxResolvedToShow; i++ {
			alert := resolvedAlerts[i]
			content += fmt.Sprintf("%s - %s\n", 
				alert.Timestamp.Format("15:04:05"),
				lipgloss.NewStyle().Faint(true).Render(alert.Message))
		}
		
		if len(resolvedAlerts) > maxResolvedToShow {
			content += fmt.Sprintf("...and %d more resolved alerts\n", 
				len(resolvedAlerts) - maxResolvedToShow)
		}
	}
	
	return infoSectionStyle.Width(width).Render(content)
}

var (
	systemCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(PaletteAccent).Padding(1).Margin(1).Background(PaletteSurface)
	cpuCardStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(PaletteSuccess).Padding(1).Margin(1).Background(PaletteSurface)
	memoryCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#60A5FA")).Padding(1).Margin(1).Background(PaletteSurface)
	diskCardStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(PaletteWarning).Padding(1).Margin(1).Background(PaletteSurface)
	networkCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#06B6D4")).Padding(1).Margin(1).Background(PaletteSurface)
	processCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FB7185")).Padding(1).Margin(1).Background(PaletteSurface)
	alertsCardStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(PaletteCritical).Padding(1).Margin(1).Background(PaletteSurface)
)

// CombinedView renders all sections in one scrollable view with colored headers
func (d *Dashboard) CombinedView(metrics *system.Collector) string {
	// Enhanced combined view (overview)
	// Layout: top row: system | cpu | memory
	// middle row: disk | network | alerts
	// bottom: full-width processes table

	width := d.width - 6 // account for some padding
	if width < 40 {
		// fallback to simple combined rendering when terminal is very narrow
		return infoSectionStyle.Width(width).Render(d.SystemView(metrics) + "\n" + d.CPUView(metrics) + "\n" + d.MemoryView(metrics) + "\n" + d.ProcessView(metrics))
	}

	// compute column widths
	colWidth := (width - 6) / 3

	sysCard := systemCardStyle.Width(colWidth).Render(d.SystemView(metrics))
	cpuCard := cpuCardStyle.Width(colWidth).Render(d.CPUView(metrics))
	memCard := memoryCardStyle.Width(colWidth).Render(d.MemoryView(metrics))

	diskCard := diskCardStyle.Width(colWidth).Render(d.DiskView(metrics))
	netCard := networkCardStyle.Width(colWidth).Render(d.NetworkView(metrics))
	alertsCard := alertsCardStyle.Width(colWidth).Render(d.AlertsView(metrics))

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, sysCard, cpuCard, memCard)
	midRow := lipgloss.JoinHorizontal(lipgloss.Top, diskCard, netCard, alertsCard)

	// Processes: render a full-width processes section beneath
	procSection := processCardStyle.Width(width).Render(d.ProcessView(metrics))

	content := lipgloss.JoinVertical(lipgloss.Left, topRow, midRow, procSection)
	return content
}

// ActiveTabContent returns the content for the currently active tab
func (d *Dashboard) ActiveTabContent(metrics *system.Collector) string {
	switch d.activeTab {
	case 0:
		return d.SystemView(metrics)
	case 1:
		return d.CPUView(metrics)
	case 2:
		return d.MemoryView(metrics)
	case 3:
		return d.DiskView(metrics)
	case 4:
		return d.NetworkView(metrics)
	case 5:
		return d.ProcessView(metrics)
	case 6:
		return d.AlertsView(metrics)
	case 7:
		return d.CombinedView(metrics)
	default:
		return "Unknown tab"
	}
}

// formatDuration formats a time.Duration to a human-readable string
func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// formatValue formats a percentage value with color coding
func formatValue(value float64) string {
	if value < 60 {
		return normalValueStyle.Render(fmt.Sprintf("%.1f%%", value))
	} else if value < 85 {
		return warnValueStyle.Render(fmt.Sprintf("%.1f%%", value))
	}
	return criticalValueStyle.Render(fmt.Sprintf("%.1f%%", value))
}
