package zettels

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	//selectedZettelStyle = styles.SecondaryColor.Copy().Width(40)
	//defaultZettelStyle  = styles.PrimaryGray.Copy().Width(40)
	nameStyle = styles.PrimaryGray.Copy().Width(45)
	metaStyle = styles.SecondaryGray.Copy().Italic(true).Width(45).Align(lipgloss.Right)

	selectedZettelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(90)
	defaultZettelStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(90)
)

func renderZettel(zettel *domain.Zettel, selected bool) string {
	style := defaultZettelStyle
	if selected {
		style = selectedZettelStyle
	}

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Left, nameStyle.Render(zettel.Name), metaStyle.Render(zettel.Tags)))
}
