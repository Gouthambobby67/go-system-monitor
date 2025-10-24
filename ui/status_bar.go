package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// StatusBar represents the top status bar of the application
type StatusBar struct {
	width     int
	height    int
	metrics   *system.Collector
	startTime time.Time
}

// NewStatusBar creates a new status bar
func NewStatusBar() *StatusBar {
	return &StatusBar{
		startTime: time.Now(),
	}
}

// SetSize updates the status bar dimensions
func (s *StatusBar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Update updates the metrics reference
func (s *StatusBar) Update(metrics *system.Collector) {
	s.metrics = metrics
}

// Render draws the status bar
func (s *StatusBar) Render() string {
	if s.metrics == nil {
		return ""
	}

	// Left section: hostname and uptime
	hostname := s.metrics.System.Hostname
	uptime := time.Since(s.startTime).Round(time.Second)
	left := fmt.Sprintf("%s | Up %s", hostname, formatDuration(uptime))

	// Center section: key metrics
	cpu := fmt.Sprintf("CPU: %s", formatValue(s.metrics.CPU.Usage))
	mem := fmt.Sprintf("MEM: %s", formatValue(s.metrics.Memory.UsedPercent))
	load := "LOAD: N/A"
	if s.metrics.CPU.LoadAvg != nil {
		load = fmt.Sprintf("LOAD: %.2f", s.metrics.CPU.LoadAvg.Load1)
	}
	center := fmt.Sprintf("%s | %s | %s", cpu, mem, load)

	// Right section: time and active alerts
	activeAlerts := 0
	for _, alert := range s.metrics.AlertManager.Alerts {
		if !alert.Resolved {
			activeAlerts++
		}
	}
	alerts := ""
	if activeAlerts > 0 {
		alerts = CriticalStyle.Render(fmt.Sprintf("⚠ %d", activeAlerts))
	}
	clock := time.Now().Format("15:04:05")
	right := fmt.Sprintf("%s %s", alerts, clock)

	// Calculate spacing
	totalLen := len(stripAnsi(left)) + len(stripAnsi(center)) + len(stripAnsi(right))
	availSpace := s.width - totalLen
	if availSpace < 0 {
		// Fallback for very narrow terminals
		return StatusBarStyle.Width(s.width).Render(left)
	}

	leftPad := availSpace / 3
	rightPad := availSpace - leftPad

	return StatusBarStyle.Width(s.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			left+strings.Repeat(" ", leftPad),
			center,
			strings.Repeat(" ", rightPad)+right,
		),
	)
}

// stripAnsi removes ANSI escape sequences for length calculations
func stripAnsi(str string) string {
	var result strings.Builder
	inEscape := false
	for _, c := range str {
		if c == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(c)
	}
	return result.String()
}

// Helper function to format duration
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// Helper function to format percentage values with color
func formatValue(value float64) string {
	style := StyleValue(value)
	return style.Render(fmt.Sprintf("%.1f%%", value))
}