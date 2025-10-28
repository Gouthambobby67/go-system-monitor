package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// ProcessTable renders an enhanced process list
type ProcessTable struct {
	width  int
	height int
	sortBy system.SortType
}

// Column definitions for the process table
var columns = []struct {
	title  string
	width  int
	align  lipgloss.Position
	format func(p system.ProcessDetail) string
}{
	{"PID", 8, lipgloss.Right, func(p system.ProcessDetail) string {
		pidStyle := BaseStyle
		// ProcessDetail doesn't include priority in this collector; highlight by CPU instead
		if p.CPUPercent > 80 {
			pidStyle = NormalStyle.Copy().Bold(true)
		}
		return pidStyle.Render(fmt.Sprintf("%7d", p.PID))
	}},
	{"CPU%", 7, lipgloss.Right, func(p system.ProcessDetail) string {
		style := StyleValue(p.CPUPercent)
		if p.CPUPercent >= 1.0 {
			return style.Bold(true).Render(fmt.Sprintf("%6.1f", p.CPUPercent))
		}
		return style.Render(fmt.Sprintf("%6.1f", p.CPUPercent))
	}},
	{"MEM%", 7, lipgloss.Right, func(p system.ProcessDetail) string {
		style := StyleValue(float64(p.MemPercent))
		if p.MemPercent >= 1.0 {
			return style.Bold(true).Render(fmt.Sprintf("%6.1f", p.MemPercent))
		}
		return style.Render(fmt.Sprintf("%6.1f", p.MemPercent))
	}},
	{"STATE", 8, lipgloss.Left, func(p system.ProcessDetail) string {
		if len(p.Status) == 0 {
			return WarningStyle.Render("? ???")
		}
		state := p.Status[0]
		switch state {
		case "D": // Uninterruptible sleep
			return WarningStyle.Render("⌛ WAIT")
		case "R": // Running
			return NormalStyle.Render("▶ RUN")
		case "S": // Interruptible sleep
			return BaseStyle.Render("💤 slp")
		case "T": // Stopped
			return WarningStyle.Render("⏸ STOP")
		case "Z": // Zombie
			return CriticalStyle.Render("💀 DEAD")
		default:
			return WarningStyle.Render(fmt.Sprintf("? %s", state))
		}
	}},
	{"USER", 12, lipgloss.Left, func(p system.ProcessDetail) string {
		style := BaseStyle
		if p.Username == "root" {
			style = WarningStyle.Copy().Bold(true)
		}
		name := p.Username
		if len(name) > 12 {
			name = name[:9] + "..."
		}
		return style.Render(name)
	}},
	{"THREADS", 5, lipgloss.Right, func(p system.ProcessDetail) string {
		style := BaseStyle
		// ProcessDetail does not provide thread count; show placeholder
		if p.CPUPercent > 50 {
			style = WarningStyle
		}
		return style.Render(fmt.Sprintf("%4s", "-"))
	}},
	{"NAME", 30, lipgloss.Left, func(p system.ProcessDetail) string {
		style := BaseStyle
		// Highlight system services and daemons
		if strings.HasSuffix(p.Name, "d") || strings.HasSuffix(p.Name, "daemon") {
			style = style.Foreground(Theme.Purple)
		}
		name := p.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		return style.Render(name)
	}},
}

// NewProcessTable creates a new process table
func NewProcessTable() *ProcessTable {
	return &ProcessTable{
		sortBy: system.SortByCPU,
	}
}

// SetSize updates the table dimensions
func (pt *ProcessTable) SetSize(width, height int) {
	pt.width = width
	pt.height = height
}

// SetSortBy updates the sort method
func (pt *ProcessTable) SetSortBy(sortBy system.SortType) {
	pt.sortBy = sortBy
}

// Render draws the process table
func (pt *ProcessTable) Render(processes []system.ProcessDetail) string {
	if len(processes) == 0 {
		return CardStyle.Render("No processes found")
	}

	// Calculate available width for columns
	fixedWidth := 0
	for _, col := range columns[:len(columns)-1] { // exclude NAME column
		fixedWidth += col.width + 2 // +2 for better spacing
	}

	// Create header
	var headers []string
	for _, col := range columns {
		style := TableHeaderStyle.Width(col.width).Align(col.align)
		if col.title == "STATE" {
			headers = append(headers, style.Render("STATUS"))
		} else {
			headers = append(headers, style.Render(col.title))
		}
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top, headers...)

	// Adjust name column width based on available space
	nameColWidth := pt.width - fixedWidth - 4 // -4 for margins
	if nameColWidth < 20 {
		nameColWidth = 20
	}
	columns[len(columns)-1].width = nameColWidth

	// Build rows
	var rows []string
	maxRows := (pt.height - 4) // account for header and margins
	if maxRows > len(processes) {
		maxRows = len(processes)
	}

	// Add a separator line after header
	separator := strings.Repeat("─", pt.width-4)
	rows = append(rows, BaseStyle.Foreground(Theme.Border).Render(separator))

	for i := 0; i < maxRows; i++ {
		proc := processes[i]
		var cells []string

		// Alternate row backgrounds for better readability
		rowStyle := TableRowStyle
		if i%2 == 1 {
			rowStyle = rowStyle.Background(Theme.Surface)
		}

		// Highlight resource-intensive processes
		if proc.CPUPercent > 50 || float64(proc.MemPercent) > 50 {
			rowStyle = rowStyle.Bold(true)
		}

		for _, col := range columns {
			cellContent := col.format(proc)
			cell := rowStyle.
				Width(col.width).
				Align(col.align).
				Render(cellContent)
			cells = append(cells, cell)
		}

		// Join cells with proper spacing
		row := lipgloss.JoinHorizontal(lipgloss.Top, cells...)

		rows = append(rows, row)
	}

	// Show sort method and process count in the header
	title := fmt.Sprintf("Processes (%d) - Sorted by %s",
		len(processes),
		getSortMethodName(pt.sortBy))

	// Create scrollbar if needed
	var scrollbar string
	if len(processes) > maxRows {
		scrollPercent := float64(maxRows) / float64(len(processes))
		scrollbarHeight := int(float64(pt.height-4) * scrollPercent)
		if scrollbarHeight < 1 {
			scrollbarHeight = 1
		}

		scrollbar = lipgloss.NewStyle().
			Foreground(Theme.Border).
			Render(strings.Repeat("│", pt.height-4))

		scrollHandle := lipgloss.NewStyle().
			Foreground(Theme.Accent).
			Render(strings.Repeat("┃", scrollbarHeight))

		// Position the scroll handle
		scrollStartPos := 0
		if len(processes) > maxRows {
			scrollStartPos = int(float64(maxRows) / float64(len(processes)) * float64(pt.height-4))
		}
		if scrollStartPos+scrollbarHeight > pt.height-4 {
			scrollStartPos = pt.height - 4 - scrollbarHeight
		}

		// Insert the handle into the scrollbar
		scrollRunes := []rune(scrollbar)
		handleRunes := []rune(scrollHandle)
		for i := 0; i < len(handleRunes) && scrollStartPos+i < len(scrollRunes); i++ {
			scrollRunes[scrollStartPos+i] = handleRunes[i]
		}
		scrollbar = string(scrollRunes)
	}

	// Join the header, separator, and rows
	table := lipgloss.JoinVertical(lipgloss.Left,
		append([]string{header}, rows...)...)

	// Add scrollbar if it exists
	if scrollbar != "" {
		table = lipgloss.JoinHorizontal(lipgloss.Top,
			table,
			" ",
			scrollbar,
		)
	}

	return CardStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			table,
		),
	)
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
