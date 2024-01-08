package accomplishments

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

func renderAccomplishment(accomplishment *domain.Accomplishment, selected bool) string {
	selectedMarker := " "
	style := defaultCycleStyle
	if selected {
		selectedMarker = ">"
		style = activeCycleStyle
	}
	description := fmt.Sprintf("%s %s", titleKeyStyle.Render("Accomplishment:"), style.Render(accomplishment.Description))
	impact := fmt.Sprintf(" %s %s", titleKeyStyle.Render("Impact:"), style.Render(accomplishment.Impact))
	s := lipgloss.JoinVertical(lipgloss.Right, description, impact)
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), alignStyle.Render(s))
}
