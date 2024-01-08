package accomplishment

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	activeCycleStyle  = styles.SecondaryColor.Copy().Width(50)
	defaultCycleStyle = styles.PrimaryGray.Copy().Width(50)
	cursorStyle       = styles.PrimaryColor.Copy()
)

func renderMinimumTask(task *domain.Task, selected bool) string {
	selectedMarker := " "
	style := defaultCycleStyle
	if selected {
		selectedMarker = ">"
		style = activeCycleStyle
	}
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), style.Render(task.Title))
}
