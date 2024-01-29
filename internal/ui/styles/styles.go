package styles

import "github.com/charmbracelet/lipgloss"

var (
	WindowStyle       = lipgloss.NewStyle().Padding(2).Width(100).Height(45).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	FormTitleStyle    = SecondaryColor.Copy()
	FormLabelStyle    = SecondaryGray.Copy()
	FormErrorStyle    = PrimaryColor.Copy()
	FocusedInputStyle = AccentColor.Copy()
)
