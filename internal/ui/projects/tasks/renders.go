package tasks

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
)

var (
	taskStyle                 = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	activeStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("#2dd4bf"))
	lowPriorityStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	highPriorityStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	selectedStyle             = lipgloss.NewStyle().Background(lipgloss.Color("#1f2937"))
	selectedTaskStyle         = selectedStyle.Copy().Foreground(lipgloss.Color("#cbd5e1"))
	selectedActiveStyle       = selectedStyle.Copy().Foreground(lipgloss.Color("#99f6e4"))
	selectedLowPriorityStyle  = selectedStyle.Copy().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	selectedHighPriorityStyle = selectedStyle.Copy().Foreground(lipgloss.Color("#ef4444")).Bold(true)
)

func renderTask(t *domain.Task, selected bool) string {
	var style lipgloss.Style
	var priorityStyle lipgloss.Style
	var priorityMarker string
	// TODO I need to clean all this logic up
	activeMarker := "\uF4C3"

	switch {
	case selected && t.Active:
		style = selectedActiveStyle
		activeMarker = "\uF444"
	case selected:
		style = selectedTaskStyle
	case t.Active:
		style = activeStyle
		activeMarker = "\uF444"
	default:
		style = taskStyle
	}

	switch {
	case t.Priority == domain.TaskPriorityNone && selected:
		priorityStyle = selectedStyle
	case t.Priority == domain.TaskPriorityLow && selected:
		priorityStyle = selectedLowPriorityStyle
		priorityMarker = " *"
	case t.Priority == domain.TaskPriorityHigh && selected:
		priorityStyle = selectedHighPriorityStyle
		priorityMarker = " ***"
	case t.Priority == domain.TaskPriorityLow:
		priorityStyle = lowPriorityStyle
		priorityMarker = " *"
	case t.Priority == domain.TaskPriorityHigh:
		priorityStyle = highPriorityStyle
		priorityMarker = " ***"
	}

	content := style.Width(30).Render(activeMarker, t.Title)
	date := style.Width(10).Italic(true).Render(t.Due.Format("2006/01/02"))
	priority := priorityStyle.Width(5).Render(priorityMarker)
	return lipgloss.JoinHorizontal(lipgloss.Top, content, priority, date)
}
