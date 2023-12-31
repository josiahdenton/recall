package detailed

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/pages/styles"
	"github.com/josiahdenton/recall/internal/pages/tasks"
)

var (
	itemStyle        = styles.PrimaryGray.Copy().PaddingLeft(2)
	cursorStyle      = styles.PrimaryColor.Copy()
	statusStyle      = lipgloss.NewStyle().Width(60).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	hiStatusStyle    = lipgloss.NewStyle().Width(60).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF"))
	titleStyle       = styles.SecondaryGray.Copy()
	activeTitleStyle = styles.SecondaryColor.Copy()
	metaTitleStyle   = styles.SecondaryGray.Copy()
)

func renderResource(r *tasks.Resource, selected bool) string {
	return renderItem(r.Name, selected)
}

func renderStatus(s *tasks.Status, selected bool) string {
	style := statusStyle
	if selected {
		style = hiStatusStyle
	}
	return style.Render(s.Description)
}

func renderStep(s *tasks.Step, selected bool) string {
	selectedMarker := " "
	completedMarker := "\uE640"
	if selected {
		selectedMarker = ">"

	}
	if s.Complete {
		completedMarker = "\U000F0856"
	}

	return fmt.Sprintf("%s%s%s", cursorStyle.Render(selectedMarker), itemStyle.Render(completedMarker), itemStyle.Render(s.Description))
}

func renderItem(s string, selected bool) string {
	selectedMarker := " "
	if selected {
		selectedMarker = ">"

	}
	return fmt.Sprintf("%s %s", cursorStyle.Render(selectedMarker), itemStyle.Render(s))
}

func renderHeader(task *tasks.Task, headerActive bool) string {
	style := titleStyle
	if headerActive {
		style = activeTitleStyle
	}

	var b strings.Builder
	b.WriteString(style.Render(task.Title) + "\n")
	b.WriteString(fmt.Sprintf("%s  %s\n\n", metaTitleStyle.Render("Due"), titleStyle.Render(task.Due)))
	return b.String()
}
