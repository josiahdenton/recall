package accomplishment

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedTasksStyle = styles.PrimaryColor.Copy().Width(50)
	defaultTaskStyle   = styles.PrimaryGray.Copy().Width(50)
	cursorStyle        = styles.PrimaryColor.Copy()
)

func renderMinimumTask(task *domain.Task, selected bool) string {
	selectedMarker := " "
	style := defaultTaskStyle
	if selected {
		selectedMarker = ">"
		style = selectedTasksStyle
	}
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), style.Render(task.Title))
}
