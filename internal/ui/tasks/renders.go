package tasks

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"reflect"
)

var (
	keyTitleStyle        = styles.SecondaryGray.Copy()
	lowPriorityStyle     = styles.AccentColor.Copy().Bold(true)
	medPriorityStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Bold(true)
	highPriorityStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true)
	favoriteStyle        = styles.AccentColor.Copy().PaddingRight(1)
	activeSpinnerStyle   = styles.SecondaryColor.Copy()
	inactiveSpinnerStyle = styles.SecondaryGray.Copy()
	selectedTaskStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(90)
	defaultTaskStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(90)
	metaKeyStyle         = styles.SecondaryGray.Copy()
	activeValueStyle     = styles.SecondaryColor.Copy()
	inactiveValueStyle   = styles.PrimaryGray.Copy()
)

func renderTask(t *domain.Task, selected bool) string {
	boxStyle := defaultTaskStyle
	if selected {
		boxStyle = selectedTaskStyle
	}

	style := inactiveValueStyle
	if t.Active {
		style = activeValueStyle
	}

	//var priorityStyle lipgloss.Style
	//var priorityMarker string
	//switch {
	//case t.Priority == domain.TaskPriorityLow:
	//	priorityStyle = lowPriorityStyle
	//	priorityMarker = " *"
	//case t.Priority == domain.TaskPriorityMedium:
	//	priorityStyle = medPriorityStyle
	//	priorityMarker = " **"
	//case t.Priority == domain.TaskPriorityHigh:
	//	priorityStyle = highPriorityStyle
	//	priorityMarker = " ***"
	//}

	//activeTime := t.ActiveDuration()
	//var duration string
	//if reflect.ValueOf(activeTime).IsZero() {
	//	duration = "(no active time)"
	//} else if activeTime.Minutes()/60 > 1 {
	//	duration = fmt.Sprintf("%d hr, %d min", activeTime/60, activeTime%60)
	//} else {
	//	duration = fmt.Sprintf("%d hr, %d min", activeTime/60, activeTime%60)
	//}

	title := style.Width(30).Render(t.Title)
	//priority := priorityStyle.Width(10).Render(priorityMarker)
	var date string
	if reflect.ValueOf(t.Due).IsZero() {
		date = lipgloss.JoinHorizontal(lipgloss.Left, keyTitleStyle.Render("Due "), style.Width(15).Italic(true).Render("None"))
	} else {
		date = lipgloss.JoinHorizontal(lipgloss.Left, keyTitleStyle.Render("Due "), style.Width(15).Italic(true).Render(t.Due.Format("2006/01/02")))
	}

	return boxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, title, date, metaKeyStyle.Render(t.Tags)))
}
