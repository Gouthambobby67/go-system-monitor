package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2)

	tabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ADB5BD")).
		PaddingLeft(2).
		PaddingRight(2)

	activeTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#383838")).
		PaddingLeft(2).
		PaddingRight(2)

	infoSectionStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#888888")).
		Padding(1).
		MarginTop(1)

	normalValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#82C91E"))

	warnValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FCC419"))

	criticalValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FA5252"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)
)

// Dashboard represents the main dashboard view
type Dashboard struct {
	tabs      []string
	activeTab int
	width     int
	height    int
	help      help.Model
}

// NewDashboard creates a new dashboard
func NewDashboard() Dashboard {
	return Dashboard{
		tabs:      []string{"System", "CPU", "Memory", "Disk", "Network", "Processes", "Alerts"},
		activeTab: 0,
		help:      help.New(),
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
