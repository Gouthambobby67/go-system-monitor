package ui

import "github.com/charmbracelet/lipgloss"

// Define styles for the UI components
var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F8FAFC")).
		Background(lipgloss.Color("#7C3AED")).
		PaddingLeft(2).
		PaddingRight(2)
)

// New palette (used across the UI)
var (
	PaletteBackground = lipgloss.Color("#0F172A") // dark navy
	PaletteSurface    = lipgloss.Color("#0B1220") // slightly lighter
	PaletteAccent     = lipgloss.Color("#7C3AED") // purple
	PaletteMuted      = lipgloss.Color("#94A3B8") // muted text
	PaletteSuccess    = lipgloss.Color("#16A34A") // green
	PaletteWarning    = lipgloss.Color("#F59E0B") // amber
	PaletteCritical   = lipgloss.Color("#EF4444") // red
)

