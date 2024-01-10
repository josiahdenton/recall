package tasks

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	taskStyle               = styles.PrimaryGray.Copy()
	keyTitleStyle           = styles.SecondaryGray.Copy()
	lowPriorityStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	highPriorityStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	selectedTaskStyle       = styles.PrimaryColor.Copy()
	activeTaskStyle         = styles.SecondaryColor.Copy()
	activeSelectedTaskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#30b3a2"))
)

func renderTask(t *domain.Task, selected bool) string {
	var style lipgloss.Style
	var priorityStyle lipgloss.Style
	var priorityMarker string
	// TODO I need to clean all this logic up
	cursor := ""

	switch {
	case selected && t.Active:
		style = activeSelectedTaskStyle
		cursor = ">"
	case selected:
		cursor = ">"
		style = selectedTaskStyle
	case t.Active:
		style = activeTaskStyle
	default:
		style = taskStyle
	}

	switch {
	case t.Priority == domain.TaskPriorityLow:
		priorityStyle = lowPriorityStyle
		priorityMarker = " *"
	case t.Priority == domain.TaskPriorityHigh:
		priorityStyle = highPriorityStyle
		priorityMarker = " ***"
	}

	title := fmt.Sprintf("%s%s", style.Width(2).Render(cursor), style.Width(30).Render(t.Title))
	date := fmt.Sprintf("%s%s", keyTitleStyle.Render("Due "), style.Width(10).Italic(true).Render(t.Due.Format("2006/01/02")))
	priority := priorityStyle.Width(10).Render(priorityMarker)
	return lipgloss.JoinHorizontal(lipgloss.Top, title, priority, date)
}
