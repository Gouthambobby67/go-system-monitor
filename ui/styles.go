package ui

import "github.com/charmbracelet/lipgloss"

// Define styles for the UI components
var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2)
)
