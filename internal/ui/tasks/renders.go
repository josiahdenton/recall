package tasks

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"reflect"
)

var (
	taskStyle               = styles.PrimaryGray.Copy()
	keyTitleStyle           = styles.SecondaryGray.Copy()
	lowPriorityStyle        = styles.AccentColor.Copy().Bold(true)
	medPriorityStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	highPriorityStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	selectedTaskStyle       = styles.SecondaryColor.Copy()
	activeTaskStyle         = styles.SecondaryColor.Copy()
	activeSelectedTaskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#30b3a2"))
	cursorStyle             = styles.PrimaryColor.Copy()
	favoriteStyle           = styles.AccentColor.Copy().PaddingLeft(2)
)

func renderTask(t *domain.Task, selected bool) string {
	var style lipgloss.Style
	var priorityStyle lipgloss.Style
	var priorityMarker string
	// TODO I need to clean all this logic up
	favoriteMarker := ""
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
	case t.Priority == domain.TaskPriorityMedium:
		priorityStyle = medPriorityStyle
		priorityMarker = " **"
	case t.Priority == domain.TaskPriorityHigh:
		priorityStyle = highPriorityStyle
		priorityMarker = " ***"
	}

	if t.Favorite {
		favoriteMarker = "\uF005"
	}

	title := lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Width(2).Render(cursor), style.Width(30).Render(t.Title))
	priority := priorityStyle.Width(10).Render(priorityMarker)
	var date string
	if reflect.ValueOf(t.Due).IsZero() {
		date = lipgloss.JoinHorizontal(lipgloss.Left, keyTitleStyle.Render("Due "), style.Width(10).Italic(true).Render("None"))
	} else {
		date = lipgloss.JoinHorizontal(lipgloss.Left, keyTitleStyle.Render("Due "), style.Width(10).Italic(true).Render(t.Due.Format("2006/01/02")))
	}
	favorite := favoriteStyle.Render(favoriteMarker)
	return lipgloss.JoinHorizontal(lipgloss.Top, title, priority, date, favorite)
}
