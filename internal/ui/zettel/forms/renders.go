package forms

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedOptionStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF")).Width(25)
	defaultOptionStyle  = styles.SecondaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(25)
	selectedZettelStyle = styles.SecondaryColor.Copy()
	defaultZettelStyle  = styles.SecondaryGray.Copy()
	cursorStyle         = styles.PrimaryColor.Copy().Width(2)
)

func renderOption(option *createZettelOption, selected bool) string {
	var s string
	if selected {
		s = selectedOptionStyle.Render(option.DisplayName)
	} else {
		s = defaultOptionStyle.Render(option.DisplayName)
	}
	return s
}

func renderZettel(zettel *domain.Zettel, selected bool) string {
	style := defaultZettelStyle
	cursor := ""
	if selected {
		style = selectedZettelStyle
		cursor = ">"
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), style.Render(zettel.Name))
}
