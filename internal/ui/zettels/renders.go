package zettels

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedZettelStyle = styles.SecondaryColor.Copy()
	defaultZettelStyle  = styles.PrimaryGray.Copy()
	cursorStyle         = styles.PrimaryColor.Copy().PaddingRight(1)
)

func renderZettel(zettel *domain.Zettel, selected bool) string {
	cursor := " "
	style := defaultZettelStyle
	if selected {
		style = selectedZettelStyle
		cursor = ">"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), style.Render(zettel.Name))
}
