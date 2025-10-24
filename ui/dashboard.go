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
			helpText = append(helpText, "", "Process Sorting:",
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
		basicHelp := d.activeTab == 5 ?
			"Tab/←→: Navigate • q: Quit • r: Refresh • 1-4: Sort • ?: Help" :
			"Tab/←→: Navigate • q: Quit • r: Refresh • ?: Help"
		elements = append(elements, helpStyle.Render(basicHelp))
	}

	return lipgloss.JoinVertical(lipgloss.Left, elements...)

// Helper methods for rendering different views

func (d *Dashboard) renderOverview(metrics *system.Collector) string {
	cpuUsage := RenderProgress("CPU Usage", metrics.CPU.TotalUsage, d.width-4)
	memUsage := RenderProgress("Memory Usage", metrics.Memory.UsagePercent, d.width-4)
	diskUsage := RenderProgress("Disk Usage", metrics.Disk.UsagePercent, d.width-4)
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		cpuUsage,
		memUsage,
		diskUsage,
	)
	
	return infoSectionStyle.Width(d.width - 4).Render(content)
}

func (d *Dashboard) renderCPU(metrics *system.Collector) string {
	var content []string
	
	content = append(content, RenderProgress("Total CPU", metrics.CPU.TotalUsage, d.width-4))
	for i, usage := range metrics.CPU.PerCPUUsage {
		content = append(content, RenderProgress(fmt.Sprintf("CPU %d", i), usage, d.width-4))
	}
	
	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderMemory(metrics *system.Collector) string {
	var content []string
	
	content = append(content, RenderProgress("Memory Usage", metrics.Memory.UsagePercent, d.width-4))
	content = append(content, fmt.Sprintf("Total: %s", FormatBytes(metrics.Memory.Total)))
	content = append(content, fmt.Sprintf("Used: %s", FormatBytes(metrics.Memory.Used)))
	content = append(content, fmt.Sprintf("Free: %s", FormatBytes(metrics.Memory.Free)))
	
	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderDisk(metrics *system.Collector) string {
	var content []string
	
	content = append(content, RenderProgress("Disk Usage", metrics.Disk.UsagePercent, d.width-4))
	content = append(content, fmt.Sprintf("Total: %s", FormatBytes(metrics.Disk.Total)))
	content = append(content, fmt.Sprintf("Used: %s", FormatBytes(metrics.Disk.Used)))
	content = append(content, fmt.Sprintf("Free: %s", FormatBytes(metrics.Disk.Free)))
	
	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderNetwork(metrics *system.Collector) string {
	var content []string
	
	content = append(content, fmt.Sprintf("Network In: %s/s", FormatBytes(metrics.Network.BytesRecv)))
	content = append(content, fmt.Sprintf("Network Out: %s/s", FormatBytes(metrics.Network.BytesSent)))
	content = append(content, fmt.Sprintf("Packets In: %d/s", metrics.Network.PacketsRecv))
	content = append(content, fmt.Sprintf("Packets Out: %d/s", metrics.Network.PacketsSent))
	
	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}

func (d *Dashboard) renderAlerts(metrics *system.Collector) string {
	if len(metrics.Alerts) == 0 {
		return infoSectionStyle.Width(d.width - 4).Render("No active alerts")
	}
	
	var content []string
	for _, alert := range metrics.Alerts {
		style := normalValueStyle
		if alert.Level == "warning" {
			style = warnValueStyle
		} else if alert.Level == "critical" {
			style = criticalValueStyle
		}
		
		content = append(content, style.Render(fmt.Sprintf("[%s] %s", 
			strings.ToUpper(alert.Level), 
			alert.Message)),
		)
	}
	
	return infoSectionStyle.Width(d.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left, content...),
	)
}
