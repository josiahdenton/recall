package zettels

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedZettelStyle = styles.SecondaryColor.Copy().Width(40)
	defaultZettelStyle  = styles.PrimaryGray.Copy().Width(40)
	cursorStyle         = styles.PrimaryColor.Copy().PaddingRight(1)
	metaStyle           = styles.SecondaryGray.Copy()
)

func renderZettel(zettel *domain.Zettel, selected bool) string {
	cursor := " "
	style := defaultZettelStyle
	if selected {
		style = selectedZettelStyle
		cursor = ">"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), style.Render(zettel.Name), metaStyle.Render(zettel.Tags))
}
