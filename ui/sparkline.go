package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// Sparkline characters for different resolutions
var (
	sparkCharsLow    = []rune{'▁', '▄', '█'}
	sparkCharsMedium = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
)

// RenderSparkline creates a sparkline visualization of time series data
func RenderSparkline(ts system.TimeSeries, width int, style lipgloss.Style) string {
	if len(ts.Points) < 2 {
		return style.Render(strings.Repeat("▁", width))
	}

	// Find min/max for scaling
	min, max := ts.Points[0].Value, ts.Points[0].Value
	for _, p := range ts.Points {
		if p.Value < min {
			min = p.Value
		}
		if p.Value > max {
			max = p.Value
		}
	}
	if min == max {
		max = min + 1 // avoid division by zero
	}

	// Use medium resolution for normal width, low for very narrow
	chars := sparkCharsMedium
	if width < 10 {
		chars = sparkCharsLow
	}
	numLevels := len(chars)

	// Generate sparkline
	var spark strings.Builder
	points := ts.Points
	if len(points) > width {
		// If we have more points than width, sample evenly
		stride := len(points) / width
		sampled := make([]system.TimeSeriesPoint, width)
		for i := 0; i < width; i++ {
			sampled[i] = points[i*stride]
		}
		points = sampled
	} else if len(points) < width {
		// If we have fewer points than width, repeat last point
		for len(points) < width {
			points = append(points, points[len(points)-1])
		}
	}

	// Convert values to sparkline characters
	for _, p := range points {
		// Scale value to [0, 1]
		scaled := (p.Value - min) / (max - min)
		// Map to character index
		idx := int(scaled * float64(numLevels-1))
		if idx >= numLevels {
			idx = numLevels - 1
		}
		spark.WriteRune(chars[idx])
	}

	return style.Render(spark.String())
}

// RenderSmallMetricCard creates a compact metric card with current value and sparkline
func RenderSmallMetricCard(label string, current float64, history system.TimeSeries, width int) string {
	if width < 20 {
		return "" // too narrow to render meaningfully
	}

	// Style for the value
	valueStyle := normalValueStyle
	if current >= 85 {
		valueStyle = criticalValueStyle
	} else if current >= 60 {
		valueStyle = warnValueStyle
	}

	// Calculate sparkline width (use remaining space)
	sparkWidth := width - len(label) - 10 // 10 chars for value & padding
	if sparkWidth < 5 {
		sparkWidth = 5
	}

	// Render components
	value := valueStyle.Render(formatFloat(current))
	spark := RenderSparkline(history, sparkWidth, valueStyle)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		label,
		" ",
		value,
		" ",
		spark,
	)
}

// formatFloat formats a float value with appropriate precision
func formatFloat(v float64) string {
	if v >= 100 {
		return "100%"
	}
	return sprintf("%.1f%%", v)
}