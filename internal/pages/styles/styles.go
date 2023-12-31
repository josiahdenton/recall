package styles

import "github.com/charmbracelet/lipgloss"

var (
	WindowStyle    = lipgloss.NewStyle().Padding(2).Width(80).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	FormTitleStyle = PrimaryColor.Copy()
	FormLabelStyle = SecondaryGray.Copy()
	FormErrorStyle = PrimaryColor.Copy()
)
