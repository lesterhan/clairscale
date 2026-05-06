package ui

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette
var (
	colorGreen   = lipgloss.Color("#A6E3A1")
	colorBlue    = lipgloss.Color("#89B4FA")
	colorMauve   = lipgloss.Color("#CBA6F7")
	colorRed     = lipgloss.Color("#F38BA8")
	colorYellow  = lipgloss.Color("#F9E2AF")
	colorOverlay = lipgloss.Color("#6C7086")
	colorSubtext = lipgloss.Color("#A6ADC8")
	colorText    = lipgloss.Color("#CDD6F4")
)

var (
	styleBold      = lipgloss.NewStyle().Bold(true)
	styleDim       = lipgloss.NewStyle().Foreground(colorSubtext)
	styleError     = lipgloss.NewStyle().Foreground(colorRed)
	styleOnline    = lipgloss.NewStyle().Foreground(colorGreen)
	styleOffline   = lipgloss.NewStyle().Foreground(colorOverlay)
	styleCursor    = lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
	styleActive    = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	styleFilter    = lipgloss.NewStyle().Foreground(colorMauve).Bold(true)
	styleConnected = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	styleWarning   = lipgloss.NewStyle().Foreground(colorYellow)
	styleLabel     = lipgloss.NewStyle().Foreground(colorSubtext)
	styleTitle     = lipgloss.NewStyle().Foreground(colorText).Bold(true)
)

func onlineDot(online bool) string {
	if online {
		return styleOnline.Render("●")
	}
	return styleOffline.Render("○")
}

func cursor(selected bool) string {
	if selected {
		return styleCursor.Render("▶ ")
	}
	return "  "
}
