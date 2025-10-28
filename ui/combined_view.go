//go:build ignore
// +build ignore

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// Enhanced CombinedView that supports sparklines and configurable cards
func (d *Dashboard) CombinedView(metrics *system.Collector) string {
	width := d.width - 6 // account for some padding
	if width < 40 {
		// fallback to simple combined rendering when terminal is very narrow
		return infoSectionStyle.Width(width).Render(
			d.SystemView(metrics) + "\n" +
				d.CPUView(metrics) + "\n" +
				d.MemoryView(metrics) + "\n" +
				d.ProcessView(metrics))
	}

	// Metrics badges at the very top
	compactMetrics := CompactSystemMetrics(metrics)
	var badges []string
	badgeWidth := (width - len(compactMetrics)*2) / len(compactMetrics)
	for _, m := range compactMetrics {
		style := normalValueStyle
		if m.Critical {
			style = criticalValueStyle
		} else if m.Alert {
			style = warnValueStyle
		}
		badges = append(badges,
			lipgloss.NewStyle().
				Width(badgeWidth).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder()).
				BorderForeground(style.GetForeground()).
				Render(fmt.Sprintf("%s: %s", m.Label,
					style.Render(formatFloat(m.Value)+m.Unit))))
	}
	metricsBar := lipgloss.JoinHorizontal(lipgloss.Top, badges...)

	// Only show enabled cards based on config
	var topCards, midCards []string
	colWidth := (width - 6) / 3

	if d.cardConfig.ShowSystem {
		topCards = append(topCards,
			systemCardStyle.Width(colWidth).Render(d.SystemView(metrics)))
	}
	if d.cardConfig.ShowCPU {
		cpuContent := d.CPUView(metrics)
		if d.cardConfig.ShowSparkline && len(metrics.CPU.History.Points) > 0 {
			sparkline := RenderSparkline(metrics.CPU.History, colWidth-4, normalValueStyle)
			cpuContent = cpuContent + "\n" + sparkline
		}
		topCards = append(topCards,
			cpuCardStyle.Width(colWidth).Render(cpuContent))
	}
	if d.cardConfig.ShowMemory {
		memContent := d.MemoryView(metrics)
		if d.cardConfig.ShowSparkline && len(metrics.Memory.History.Points) > 0 {
			sparkline := RenderSparkline(metrics.Memory.History, colWidth-4, normalValueStyle)
			memContent = memContent + "\n" + sparkline
		}
		topCards = append(topCards,
			memoryCardStyle.Width(colWidth).Render(memContent))
	}

	if d.cardConfig.ShowDisk {
		midCards = append(midCards,
			diskCardStyle.Width(colWidth).Render(d.DiskView(metrics)))
	}
	if d.cardConfig.ShowNetwork {
		midCards = append(midCards,
			networkCardStyle.Width(colWidth).Render(d.NetworkView(metrics)))
	}
	if d.cardConfig.ShowAlerts {
		midCards = append(midCards,
			alertsCardStyle.Width(colWidth).Render(d.AlertsView(metrics)))
	}

	// Join rows
	var content string
	if len(badges) > 0 {
		content = metricsBar + "\n\n"
	}
	if len(topCards) > 0 {
		content += lipgloss.JoinHorizontal(lipgloss.Top, topCards...) + "\n"
	}
	if len(midCards) > 0 {
		content += lipgloss.JoinHorizontal(lipgloss.Top, midCards...) + "\n"
	}

	// Always show processes at the bottom
	procSection := processCardStyle.Width(width).Render(d.ProcessView(metrics))
	content += procSection

	return content
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

// Enhanced help text including card toggles
func (d *Dashboard) RenderHelp() string {
	basic := "Tab/← →: Navigate • 1-4: Sort processes • q: Quit • r: Refresh"
	cards := "s/c/m/d/n/a: Toggle cards • g: Toggle graphs"
	return helpStyle.Render(basic + " • " + cards)
}