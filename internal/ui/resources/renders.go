package resources

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	activeResourceStyle  = styles.SelectedItemStyle.Copy().Width(40).PaddingRight(2)
	defaultResourceStyle = styles.DefaultItemStyle.Copy().Width(40).PaddingRight(2)
	cursorStyle          = styles.PrimaryColorStyle.Copy()
	metaStyle            = styles.SecondaryGrayStyle.Copy().Width(40).Align(lipgloss.Right)
	selectedMetaStyle    = styles.AccentColorStyle.Copy().Width(40).Align(lipgloss.Right)
	alignStyle           = lipgloss.NewStyle().PaddingLeft(1)
)

func renderResource(resource *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := defaultResourceStyle
	tagsStyle := metaStyle
	if selected {
		selectedMarker = ">"
		style = activeResourceStyle
		tagsStyle = selectedMetaStyle
	}
	name := style.Render(resource.Name)
	tags := tagsStyle.Render(resource.Tags)
	s := lipgloss.JoinHorizontal(lipgloss.Left, name, tags)
	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(selectedMarker), alignStyle.Render(s))
}
