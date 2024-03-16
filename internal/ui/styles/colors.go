package styles

import "github.com/charmbracelet/lipgloss"

// I should just make these colors...
var (
	PrimaryColorStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#D120AF"))
	SecondaryColorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#2dd4bf"))
	AccentColorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcd34d"))
	SecondaryAccentColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF21"))
	PrimaryGrayStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	SecondaryGrayStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

var (
	PrimaryColor       = lipgloss.Color("#D120AF")
	PrimaryGrayColor   = lipgloss.Color("#767676")
	SecondaryGrayColor = lipgloss.Color("#3a3b5b")
	AccentColor        = lipgloss.Color("#fcd34d")
)
