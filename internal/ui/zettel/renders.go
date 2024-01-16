package zettel

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	defaultZettelStyle  = styles.PrimaryGray.Copy()
	selectedZettelStyle = styles.SecondaryColor.Copy()
	cursorStyle         = styles.PrimaryColor.Copy().Width(2)
)

func renderZettel(z *domain.Zettel, selected bool) string {
	cursor := ""
	style := defaultZettelStyle
	if selected {
		cursor = ">"
		style = selectedZettelStyle
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(cursor), style.Render(z.Name))
}
