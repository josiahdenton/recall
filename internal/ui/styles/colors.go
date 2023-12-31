package styles

import "github.com/charmbracelet/lipgloss"

// I should just make these colors...
var (
	PrimaryColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#D120AF"))
	SecondaryColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#2dd4bf"))
	PrimaryGray    = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	SecondaryGray  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)
