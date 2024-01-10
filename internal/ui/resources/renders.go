package resources

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	activeCycleStyle  = styles.SecondaryColor.Copy().Width(50)
	defaultCycleStyle = styles.PrimaryGray.Copy().Width(50)
	cursorStyle       = styles.PrimaryColor.Copy()
	titleKeyStyle     = styles.SecondaryGray.Copy()
	alignStyle        = lipgloss.NewStyle().PaddingLeft(1)
)

func renderResource(resource *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := defaultCycleStyle
	if selected {
		selectedMarker = ">"
		style = activeCycleStyle
	}
	name := style.Render(resource.Name)
	resourceType := fmt.Sprintf(" %s %s", titleKeyStyle.Render("Type"), style.Render(resource.StringType()))
	s := lipgloss.JoinHorizontal(lipgloss.Left, name, resourceType)
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), alignStyle.Render(s))
}
