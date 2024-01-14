package accomplishment

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedTasksStyle = styles.SecondaryColor.Copy().Width(50)
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
	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Width(2).Render(selectedMarker), style.Render(task.Title))
}
