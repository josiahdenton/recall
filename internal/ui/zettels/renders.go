package zettels

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedZettelStyle = styles.SecondaryColor.Copy().PaddingLeft(2)
	defaultZettelStyle  = styles.PrimaryGray.Copy().PaddingLeft(2)
	cursorStyle         = styles.PrimaryColor.Copy().PaddingLeft(2)
)

func renderZettel(zettel *domain.Zettel, selected bool) string {
	cursor := ""
	style := defaultZettelStyle
	if selected {
		style = selectedZettelStyle
		cursor = ">"
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), style.Render(zettel.Name))
}
