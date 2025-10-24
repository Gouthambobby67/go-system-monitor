package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme contains all the colors used in the UI
var Theme = struct {
	// Base colors
	Background lipgloss.Color
	Surface    lipgloss.Color
	Border     lipgloss.Color
	Text       lipgloss.Color
	Muted      lipgloss.Color
	Accent     lipgloss.Color

	// Status indicators
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Critical   lipgloss.Color
	Info       lipgloss.Color

	// Special highlights
	Purple     lipgloss.Color
	Blue       lipgloss.Color
	Cyan       lipgloss.Color
	Green      lipgloss.Color
	Yellow     lipgloss.Color
	Orange     lipgloss.Color
	Red        lipgloss.Color
}{
	// Modern dark theme inspired by Nord and Tokyo Night
	Background: lipgloss.Color("#1a1b26"),
	Surface:    lipgloss.Color("#24283b"),
	Border:     lipgloss.Color("#414868"),
	Text:       lipgloss.Color("#c0caf5"),
	Muted:      lipgloss.Color("#565f89"),
	Accent:     lipgloss.Color("#7aa2f7"),

	Success:  lipgloss.Color("#9ece6a"),
	Warning:  lipgloss.Color("#e0af68"),
	Error:    lipgloss.Color("#f7768e"),
	Critical: lipgloss.Color("#db4b4b"),
	Info:     lipgloss.Color("#7dcfff"),

	Purple: lipgloss.Color("#bb9af7"),
	Blue:   lipgloss.Color("#7aa2f7"),
	Cyan:   lipgloss.Color("#7dcfff"),
	Green:  lipgloss.Color("#9ece6a"),
	Yellow: lipgloss.Color("#e0af68"),
	Orange: lipgloss.Color("#ff9e64"),
	Red:    lipgloss.Color("#f7768e"),
}

// Common style mixins
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
		Background(Theme.Background).
		Foreground(Theme.Text)

	// Card styles with modern borders
	CardStyle = BaseStyle.Copy().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Border).
		Padding(0, 1).
		Margin(0, 1)

	// Status bar style
	StatusBarStyle = BaseStyle.Copy().
		Background(Theme.Surface).
		Padding(0, 1)

	// Section header style
	HeaderStyle = BaseStyle.Copy().
		Bold(true).
		Foreground(Theme.Accent)

	// Value styles for different states
	NormalStyle = BaseStyle.Copy().
		Foreground(Theme.Success)

	WarningStyle = BaseStyle.Copy().
		Foreground(Theme.Warning)

	CriticalStyle = BaseStyle.Copy().
		Bold(true).
		Foreground(Theme.Critical)

	// Process table styles
	TableHeaderStyle = BaseStyle.Copy().
		Bold(true).
		Foreground(Theme.Text).
		Background(Theme.Surface)

	TableRowStyle = BaseStyle.Copy().
		Foreground(Theme.Text)

	TableAltRowStyle = BaseStyle.Copy().
		Foreground(Theme.Text).
		Background(Theme.Surface)

	// Keyboard shortcut style
	KeyStyle = BaseStyle.Copy().
		Background(Theme.Surface).
		Foreground(Theme.Text).
		Padding(0, 1).
		Margin(0, 1)
)

// Utility functions for consistent styling
func StyleValue(value float64) lipgloss.Style {
	switch {
	case value >= 90:
		return CriticalStyle
	case value >= 75:
		return WarningStyle
	default:
		return NormalStyle
	}
}

// Section creates a titled card with content
func Section(title string, content string) string {
	return CardStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			HeaderStyle.Render(title),
			content,
		),
	)
}