package zettel

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	defaultZettelStyle  = styles.PrimaryGray.Copy()
	selectedZettelStyle = styles.AccentColor.Copy()
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

var (
	selectedResourceStyle = styles.AccentColor.Copy().Width(50)
	defaultResourceStyle  = styles.PrimaryGray.Copy().Width(50)
	titleKeyStyle         = styles.SecondaryGray.Copy()
	alignStyle            = lipgloss.NewStyle().PaddingLeft(1)
)

func renderResource(resource *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := defaultResourceStyle
	if selected {
		selectedMarker = ">"
		style = selectedResourceStyle
	}
	name := style.Render(resource.Name)
	tags := style.Render(resource.Tags)
	s := lipgloss.JoinHorizontal(lipgloss.Left, name, tags)
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), alignStyle.Render(s))
}
