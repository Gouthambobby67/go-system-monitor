package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/system"
)

// ProcessTable renders an enhanced process list
type ProcessTable struct {
	width      int
	height     int
	sortBy     system.SortType
	scrollPos  int // Current scroll position
	filterText string
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
	{"THREADS", 8, lipgloss.Right, func(p system.ProcessDetail) string {
		style := BaseStyle
		if p.NumThreads > 100 {
			style = WarningStyle
		} else if p.NumThreads > 50 {
			style = NormalStyle.Copy().Foreground(Theme.Warning)
		}
		return style.Render(fmt.Sprintf("%7d", p.NumThreads))
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

// ScrollUp moves the view up
func (pt *ProcessTable) ScrollUp() {
	if pt.scrollPos > 0 {
		pt.scrollPos--
	}
}

// ScrollDown moves the view down
func (pt *ProcessTable) ScrollDown(maxRows int) {
	if pt.scrollPos < maxRows-1 {
		pt.scrollPos++
	}
}

// PageUp moves the view up by a page
func (pt *ProcessTable) PageUp() {
	pageSize := (pt.height - 4)
	if pageSize < 1 {
		pageSize = 10
	}
	pt.scrollPos -= pageSize
	if pt.scrollPos < 0 {
		pt.scrollPos = 0
	}
}

// PageDown moves the view down by a page
func (pt *ProcessTable) PageDown(maxRows int) {
	pageSize := (pt.height - 4)
	if pageSize < 1 {
		pageSize = 10
	}
	pt.scrollPos += pageSize
	if pt.scrollPos >= maxRows {
		pt.scrollPos = maxRows - 1
	}
	if pt.scrollPos < 0 {
		pt.scrollPos = 0
	}
}

// Home moves to the top
func (pt *ProcessTable) Home() {
	pt.scrollPos = 0
}

// End moves to the bottom
func (pt *ProcessTable) End(maxRows int) {
	pt.scrollPos = maxRows - 1
	if pt.scrollPos < 0 {
		pt.scrollPos = 0
	}
}

// SetFilterText updates the filter text
func (pt *ProcessTable) SetFilterText(text string) {
	pt.filterText = text
	pt.scrollPos = 0 // Reset scroll when filtering
}

// filterProcesses filters processes based on filter text
func (pt *ProcessTable) filterProcesses(processes []system.ProcessDetail) []system.ProcessDetail {
	if pt.filterText == "" {
		return processes
	}
	
	filtered := make([]system.ProcessDetail, 0)
	filterLower := strings.ToLower(pt.filterText)
	
	for _, proc := range processes {
		// Check if filter matches name, username, or command line
		if strings.Contains(strings.ToLower(proc.Name), filterLower) ||
			strings.Contains(strings.ToLower(proc.Username), filterLower) ||
			strings.Contains(strings.ToLower(proc.CmdLine), filterLower) {
			filtered = append(filtered, proc)
		}
	}
	
	return filtered
}

// Render draws the process table
func (pt *ProcessTable) Render(processes []system.ProcessDetail) string {
	// Apply filtering
	processes = pt.filterProcesses(processes)
	
	if len(processes) == 0 {
		msg := "No processes found"
		if pt.filterText != "" {
			msg = fmt.Sprintf("No processes match filter: %s", pt.filterText)
		}
		return CardStyle.Render(msg)
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

	// Build rows with scrolling support
	var rows []string
	displayRows := (pt.height - 4) // account for header and margins
	if displayRows < 5 {
		displayRows = 5
	}
	
	// Make sure scroll position is valid
	if pt.scrollPos < 0 {
		pt.scrollPos = 0
	}
	if pt.scrollPos >= len(processes) {
		pt.scrollPos = len(processes) - 1
	}
	if pt.scrollPos < 0 {
		pt.scrollPos = 0
	}

	// Calculate visible range
	endPos := pt.scrollPos + displayRows
	if endPos > len(processes) {
		endPos = len(processes)
	}

	// Add a separator line after header
	separator := strings.Repeat("─", pt.width-4)
	rows = append(rows, BaseStyle.Foreground(Theme.Border).Render(separator))

	for i := pt.scrollPos; i < endPos; i++ {
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

	// Show sort method, process count, and scroll position in the header
	title := fmt.Sprintf("Processes (%s) - Sorted by %s - Showing %s-%s",
		FormatNumber(len(processes)),
		getSortMethodName(pt.sortBy),
		FormatNumber(pt.scrollPos+1),
		FormatNumber(endPos))
	
	// Show filter if active
	if pt.filterText != "" {
		title += fmt.Sprintf(" [Filter: %s]", pt.filterText)
	}

	// Create scrollbar if needed
	var scrollbar string
	if len(processes) > displayRows {
		scrollPercent := float64(displayRows) / float64(len(processes))
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
		if len(processes) > displayRows {
			scrollStartPos = int(float64(displayRows) / float64(len(processes)) * float64(pt.height-4))
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
