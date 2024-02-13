package styles

import "github.com/charmbracelet/lipgloss"

var (
	WindowStyle    = lipgloss.NewStyle().Padding(2).Width(100).Height(45).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	CenterStyle    = lipgloss.NewStyle().Align(lipgloss.Center)
	FormTitleStyle = lipgloss.NewStyle().Foreground(SecondaryColor)
	FormLabelStyle = lipgloss.NewStyle().Foreground(SecondaryGray)
	FormErrorStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	FocusedStyle   = lipgloss.NewStyle().Foreground(AccentColor)

	WarnToastStyle = lipgloss.NewStyle().Foreground(PrimaryGray).Border(lipgloss.RoundedBorder()).BorderForeground(PrimaryColor).Width(25).Align(lipgloss.Center)
	InfoToastStyle = lipgloss.NewStyle().Foreground(PrimaryGray).Border(lipgloss.RoundedBorder()).BorderForeground(SecondaryColor).Width(25).Align(lipgloss.Center)
)
