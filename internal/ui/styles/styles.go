package styles

import "github.com/charmbracelet/lipgloss"

var (
	WindowStyle       = lipgloss.NewStyle().Padding(2).Width(100).Height(45).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	FormTitleStyle    = SecondaryColorStyle.Copy()
	FormLabelStyle    = SecondaryGrayStyle.Copy()
	FormErrorStyle    = PrimaryColorStyle.Copy()
	FocusedInputStyle = AccentColorStyle.Copy()
)

var (
	CursorStyle       = lipgloss.NewStyle().Foreground(PrimaryColor)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(AccentColor).Bold(true)
	DefaultItemStyle  = lipgloss.NewStyle().Foreground(SecondaryGrayColor)
)

// TODO - bring over the box styling from the other PR...
