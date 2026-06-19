package tui

import lipgloss "charm.land/lipgloss/v2"

var (
	// Colors
	activeColor   = lipgloss.Color("62")  // Indigo
	inactiveColor = lipgloss.Color("240") // Dark Grey
	cursorColor   = lipgloss.Color("212") // Pink
	scrollColor   = lipgloss.Color("241") // Subtle grey for scroll indicators
	dimColor      = lipgloss.Color("245") // Dimmed text

	// Text Styles
	titleStyle  = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	cursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Bold(true)
	detailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	scrollStyle = lipgloss.NewStyle().Foreground(scrollColor)
	dimStyle    = lipgloss.NewStyle().Foreground(dimColor)
)

// paneStyle returns an active or inactive pane border style.
// Height and width are set dynamically per render, so no fixed
// dimensions here — just border + padding.
func paneStyle(active bool) lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 1)

	if active {
		return s.BorderForeground(activeColor)
	}
	return s.BorderForeground(inactiveColor)
}
