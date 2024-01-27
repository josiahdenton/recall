package artifact

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	cursorStyle            = styles.PrimaryColor.Copy()
	resourceStyle          = styles.PrimaryGray.Copy().PaddingLeft(2)
	selectedResourceStyle  = styles.SecondaryColor.Copy().PaddingLeft(2)
	resourceMetaTitleStyle = styles.SecondaryGray.Copy()
)

func renderResource(r *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := resourceStyle
	if selected {
		selectedMarker = ">"
		style = selectedResourceStyle
	}
	name := lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Width(2).Render(selectedMarker), style.Width(50).Render(r.Name))
	resourceType := lipgloss.JoinHorizontal(lipgloss.Left, resourceMetaTitleStyle.Width(5).Render("Type"), style.Width(10).Render(r.StringType()))
	return lipgloss.JoinHorizontal(lipgloss.Left, name, resourceType)
}
