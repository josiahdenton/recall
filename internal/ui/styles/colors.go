package styles

import "github.com/charmbracelet/lipgloss"

// I should just make these colors...
var (
	PrimaryColor         = lipgloss.NewStyle().Foreground(lipgloss.Color("#D120AF"))
	SecondaryColor       = lipgloss.NewStyle().Foreground(lipgloss.Color("#2dd4bf"))
	AccentColor          = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcd34d"))
	SecondaryAccentColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF21"))
	PrimaryGray          = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	SecondaryGray        = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)
