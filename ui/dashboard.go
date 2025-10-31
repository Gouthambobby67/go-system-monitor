package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

var (
	tabStyle = lipgloss.NewStyle().
		Foreground(PaletteMuted).
		Background(PaletteSurface).
		PaddingLeft(2).
		PaddingRight(2)

	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Background(PaletteAccent).
		PaddingLeft(2).
		PaddingRight(2)

	infoSectionStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(PaletteMuted).
		Background(PaletteSurface).
		Padding(1).
		MarginTop(1)

	normalValueStyle = lipgloss.NewStyle().
		Foreground(PaletteSuccess)

	warnValueStyle = lipgloss.NewStyle().
		Foreground(PaletteWarning)

	criticalValueStyle = lipgloss.NewStyle().
		Foreground(PaletteCritical)

	helpStyle = lipgloss.NewStyle().
		Foreground(PaletteMuted).
		MarginTop(1)

	// Sidebar navigation styles
	sidebarStyle = lipgloss.NewStyle().
		Width(22).
		Foreground(PaletteMuted).
		Background(PaletteSurface).
		PaddingLeft(1).
		PaddingRight(1)
	activeSidebarStyle = lipgloss.NewStyle().
		Width(22).
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Background(PaletteAccent).
		PaddingLeft(1).
		PaddingRight(1)
)

// Dashboard represents the main dashboard view
type Dashboard struct {
	tabs          []string
	activeTab     int
	width         int
	height        int
	help          help.Model
	statusBar     *StatusBar
	processTable  *ProcessTable
	compactMode   bool
	showHelp      bool
	fullscreen    bool
	showStatusBar bool
	cardConfig    CardConfig
}

// NewDashboard creates a new dashboard
func NewDashboard() Dashboard {
	return Dashboard{
		tabs:          []string{"Overview", "CPU", "Memory", "Disk", "Network", "Processes", "Alerts"},
		activeTab:     0,
		help:          help.New(),
		statusBar:     NewStatusBar(),
		processTable:  NewProcessTable(),
		showStatusBar: true,
		compactMode:   false,
		showHelp:      false,
		fullscreen:    false,
		cardConfig:    DefaultCardConfig(),
	}
}

// NextTab switches to the next tab
func (d *Dashboard) NextTab() {
	d.activeTab = (d.activeTab + 1) % len(d.tabs)
}

// PrevTab switches to the previous tab
func (d *Dashboard) PrevTab() {
	d.activeTab = (d.activeTab - 1 + len(d.tabs)) % len(d.tabs)
}

// ActiveTab returns the index of the currently active tab
func (d *Dashboard) ActiveTab() int {
	return d.activeTab
}

// SetSize sets the size for the dashboard
func (d *Dashboard) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.statusBar.SetSize(width, 1)
	d.processTable.SetSize(width-4, height-8) // account for margins and other elements
}

// FormatTabs renders the tab navigation
func (d *Dashboard) FormatTabs() string {
	var renderedTabs []string

	for i, t := range d.tabs {
		if i == d.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// FormatSidebar renders sidebar navigation with the current active section
func (d *Dashboard) FormatSidebar() string {
	var items []string
	for i, t := range d.tabs {
		if i == d.activeTab {
			items = append(items, activeSidebarStyle.Render(t))
		} else {
			items = append(items, sidebarStyle.Render(t))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// RenderProgress creates a visual progress bar
func RenderProgress(label string, value float64, width int) string {
	// Set a reasonable minimum width to ensure the bar is visible
	barWidth := width - len(label) - 10 // Subtract label length and some padding
	if barWidth < 10 {
		barWidth = 10
	}

	// Create style based on value
	var valueStyle lipgloss.Style
	if value < 60 {
		valueStyle = normalValueStyle
	} else if value < 85 {
		valueStyle = warnValueStyle
	} else {
		valueStyle = criticalValueStyle
	}

	// Calculate filled and empty portions
	filled := int((value / 100) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	
	empty := barWidth - filled

	// Create the progress bar
	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)
	
	progressBar := fmt.Sprintf(
		"%s %s%s %s", 
		label,
		valueStyle.Render(filledStr),
		lipgloss.NewStyle().Render(emptyStr),
		valueStyle.Render(fmt.Sprintf("%.1f%%", value)),
	)

	return progressBar
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatNumber formats a number with thousands separators
func FormatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	
	s := fmt.Sprintf("%d", n)
	result := ""
	count := 0
	
	for i := len(s) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result = "," + result
		}
		result = string(s[i]) + result
		count++
	}
	
	return result
}

// RenderMainContent returns the content for the currently active tab
func (d *Dashboard) RenderMainContent(metrics *system.Collector) string {
	switch d.activeTab {
	case 0: // Overview
		return d.renderOverview(metrics)
	case 1: // CPU
		return d.renderCPU(metrics)
	case 2: // Memory
		return d.renderMemory(metrics)
	case 3: // Disk
		return d.renderDisk(metrics)
	case 4: // Network
		return d.renderNetwork(metrics)
	case 5: // Processes
		return d.processTable.Render(metrics.Process.Processes)
	case 6: // Alerts
		return d.renderAlerts(metrics)
	default:
		return "Unknown tab"
	}
}

// Toggle methods
func (d *Dashboard) ToggleCompactMode() {
	d.compactMode = !d.compactMode
}

func (d *Dashboard) ToggleHelp() {
	d.showHelp = !d.showHelp
}

func (d *Dashboard) ToggleFullscreen() {
	d.fullscreen = !d.fullscreen
}

func (d *Dashboard) ToggleStatusBar() {
	d.showStatusBar = !d.showStatusBar
}

// UpdateCardConfig updates card visibility based on a key press
func (d *Dashboard) UpdateCardConfig(key string) {
	switch key {
	case "s":
		d.cardConfig.ShowSystem = !d.cardConfig.ShowSystem
	case "c":
		d.cardConfig.ShowCPU = !d.cardConfig.ShowCPU
	case "m":
		d.cardConfig.ShowMemory = !d.cardConfig.ShowMemory
	case "d":
		d.cardConfig.ShowDisk = !d.cardConfig.ShowDisk
	case "n":
		d.cardConfig.ShowNetwork = !d.cardConfig.ShowNetwork
	case "a":
		d.cardConfig.ShowAlerts = !d.cardConfig.ShowAlerts
	case "g":
		d.cardConfig.ShowSparkline = !d.cardConfig.ShowSparkline
	}
}

// Process table scrolling delegation methods
func (d *Dashboard) ScrollProcessUp() {
	d.processTable.ScrollUp()
}

func (d *Dashboard) ScrollProcessDown(maxRows int) {
	d.processTable.ScrollDown(maxRows)
}

func (d *Dashboard) PageUpProcess() {
	d.processTable.PageUp()
}

func (d *Dashboard) PageDownProcess(maxRows int) {
	d.processTable.PageDown(maxRows)
}

func (d *Dashboard) HomeProcess() {
	d.processTable.Home()
}

func (d *Dashboard) EndProcess(maxRows int) {
	d.processTable.End(maxRows)
}

// SetProcessFilter sets the filter text for the process table
func (d *Dashboard) SetProcessFilter(text string) {
	d.processTable.SetFilterText(text)
}

// Render returns the complete dashboard view
func (d *Dashboard) Render(metrics *system.Collector) string {
	if d.fullscreen {
		// In fullscreen mode, only show the current view content
		return d.RenderMainContent(metrics)
	}

	var elements []string

	// Status bar (if enabled)
	if d.showStatusBar {
		d.statusBar.Update(metrics)
		elements = append(elements, d.statusBar.Render())
	}

	// Tabs (unless in compact mode)
	if !d.compactMode {
		elements = append(elements, d.FormatTabs())
	}

	// Main content
	content := d.RenderMainContent(metrics)
	if d.compactMode {
		// In compact mode, reduce padding and borders
		content = lipgloss.NewStyle().
			MaxHeight(d.height - 2). // Leave room for status bar
			MaxWidth(d.width).
			Render(content)
	}
	elements = append(elements, content)

	// Help text or help overlay
	if d.showHelp {
		helpText := []string{
			"Keyboard Shortcuts:",
			"  Tab/←→: Change tab",
			"  q: Quit",
			"  r: Refresh",
			"  c: Toggle compact mode",
			"  f: Toggle fullscreen",
			"  s: Toggle status bar",
			"  ?: Toggle this help",
		}
		if d.activeTab == 5 { // Processes tab
			helpText = append(helpText, "", "Process Navigation:",
				"  ↑/k: Scroll up",
				"  ↓/j: Scroll down",
				"  PgUp/Ctrl+u: Page up",
				"  PgDn/Ctrl+d: Page down",
				"  Home/g: Jump to top",
				"  End/G: Jump to bottom",
				"", "Process Sorting:",
				"  1: Sort by CPU",
				"  2: Sort by Memory",
				"  3: Sort by PID",
				"  4: Sort by Name",
			)
		}
		elements = append(elements, CardStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left, helpText...)),
		)
	} else {
		var basicHelp string
		if d.activeTab == 5 {
			basicHelp = "Tab/←→: Navigate • ↑↓/jk: Scroll • PgUp/PgDn: Page • Home/End: Jump • 1-4: Sort • q: Quit • r: Refresh • ?: Help"
		} else {
			basicHelp = "Tab/←→: Navigate • q: Quit • r: Refresh • ?: Help"
		}
		elements = append(elements, helpStyle.Render(basicHelp))
	}

	return lipgloss.JoinVertical(lipgloss.Left, elements...)
}

// Helper methods for rendering different views

func (d *Dashboard) renderOverview(metrics *system.Collector) string {
	// Use the Collector's current fields
	cpuUsage := RenderProgress("CPU Usage", metrics.CPU.Usage, d.width-4)
	memUsage := RenderProgress("Memory Usage", metrics.Memory.UsedPercent, d.width-4)

	// Pick a representative disk usage (first partition) if available
	var diskUsagePercent float64
	if len(metrics.Disk.UsageStats) > 0 {
		for _, u := range metrics.Disk.UsageStats {
			diskUsagePercent = u.UsedPercent
			break
		}
	}
	diskUsage := RenderProgress("Disk Usage", diskUsagePercent, d.width-4)

	content := lipgloss.JoinVertical(lipgloss.Left,
		cpuUsage,
		memUsage,
		diskUsage,
	)

	return infoSectionStyle.Width(d.width - 4).Render(content)
}

func (d *Dashboard) renderCPU(metrics *system.Collector) string {
	var content []string
	content = append(content, RenderProgress("Total CPU", metrics.CPU.Usage, d.width-4))
	for i, usage := range metrics.CPU.UsagePerCPU {
		content = append(content, RenderProgress(fmt.Sprintf("CPU %d", i), usage, d.width-4))
	}

	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderMemory(metrics *system.Collector) string {
	var content []string
	content = append(content, RenderProgress("Memory Usage", metrics.Memory.UsedPercent, d.width-4))
	content = append(content, fmt.Sprintf("Total: %s", FormatBytes(metrics.Memory.Total)))
	content = append(content, fmt.Sprintf("Used: %s", FormatBytes(metrics.Memory.Used)))
	content = append(content, fmt.Sprintf("Free: %s", FormatBytes(metrics.Memory.Free)))

	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderDisk(metrics *system.Collector) string {
	var content []string
	// Show the first partition's usage as a representative sample
	var total, used, free uint64
	var usedPercent float64
	if len(metrics.Disk.UsageStats) > 0 {
		for _, u := range metrics.Disk.UsageStats {
			total = u.Total
			used = u.Used
			free = u.Free
			usedPercent = u.UsedPercent
			break
		}
	}

	content = append(content, RenderProgress("Disk Usage", usedPercent, d.width-4))
	content = append(content, fmt.Sprintf("Total: %s", FormatBytes(total)))
	content = append(content, fmt.Sprintf("Used: %s", FormatBytes(used)))
	content = append(content, fmt.Sprintf("Free: %s", FormatBytes(free)))

	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderNetwork(metrics *system.Collector) string {
	var content []string
	
	// Show per-interface rates if available
	if len(metrics.Network.RecvRate) > 0 {
		content = append(content, "Network I/O Rates:")
		
		// Aggregate total rates
		var totalRecv, totalSent float64
		for iface, recvRate := range metrics.Network.RecvRate {
			if iface == "lo" {
				continue // Skip loopback
			}
			sentRate := metrics.Network.SentRate[iface]
			totalRecv += recvRate
			totalSent += sentRate
			
			if recvRate > 0 || sentRate > 0 {
				content = append(content, fmt.Sprintf("  %s: ↓ %s/s  ↑ %s/s", 
					iface, 
					FormatBytes(uint64(recvRate)), 
					FormatBytes(uint64(sentRate))))
			}
		}
		
		content = append(content, "")
		content = append(content, fmt.Sprintf("Total: ↓ %s/s  ↑ %s/s", 
			FormatBytes(uint64(totalRecv)), 
			FormatBytes(uint64(totalSent))))
	} else {
		// Fallback to cumulative counters
		var bytesIn, bytesOut uint64
		for _, c := range metrics.Network.IOCounters {
			bytesIn += c.BytesRecv
			bytesOut += c.BytesSent
		}
		content = append(content, fmt.Sprintf("Total Received: %s", FormatBytes(bytesIn)))
		content = append(content, fmt.Sprintf("Total Sent: %s", FormatBytes(bytesOut)))
	}
	
	// Add connection count
	content = append(content, "")
	content = append(content, fmt.Sprintf("Active Connections: %d", len(metrics.Network.Connections)))

	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderAlerts(metrics *system.Collector) string {
	if metrics.AlertManager == nil || len(metrics.AlertManager.Alerts) == 0 {
		return infoSectionStyle.Width(d.width - 4).Render("No active alerts")
	}

	var content []string
	for _, alert := range metrics.AlertManager.Alerts {
		style := normalValueStyle
		if alert.Level == system.WarningLevel {
			style = warnValueStyle
		} else if alert.Level == system.CriticalLevel {
			style = criticalValueStyle
		}

		content = append(content, style.Render(fmt.Sprintf("[%s] %s",
			strings.ToUpper(string(alert.Level)),
			alert.Message)),
		)
	}

	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}
