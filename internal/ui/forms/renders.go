package forms

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedLinkOptionStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF")).Width(25)
	defaultLinkOptionStyle  = styles.SecondaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(25)
	selectedZettelStyle     = styles.SecondaryColor.Copy()
	defaultZettelStyle      = styles.SecondaryGray.Copy()
	cursorStyle             = styles.PrimaryColor.Copy().Width(2)
)

func renderLinkOption(option *linkZettelOption, selected bool) string {
	var s string
	if selected {
		s = selectedLinkOptionStyle.Render(option.DisplayName)
	} else {
		s = defaultLinkOptionStyle.Render(option.DisplayName)
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
