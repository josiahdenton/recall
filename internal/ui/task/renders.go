package task

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = styles.PrimaryColor.Copy()
	// steps
	stepStyle         = styles.PrimaryGray.Copy().PaddingLeft(2)
	selectedStepStyle = styles.SecondaryColor.Copy().PaddingLeft(2)
	titleStyle        = styles.SecondaryGray.Copy()
	activeTitleStyle  = styles.SecondaryColor.Copy()
	metaTitleStyle    = styles.SecondaryGray.Copy()
	// resources
	resourceStyle          = styles.PrimaryGray.Copy().PaddingLeft(2)
	selectedResourceStyle  = styles.SecondaryColor.Copy().PaddingLeft(2)
	resourceMetaTitleStyle = styles.SecondaryGray.Copy()
	// status
	statusStyle   = lipgloss.NewStyle().Width(60).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	hiStatusStyle = lipgloss.NewStyle().Width(60).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF"))
)

func renderResource(r *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := resourceStyle
	if selected {
		selectedMarker = ">"
		style = selectedResourceStyle
	}
	name := lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Width(2).Render(selectedMarker), style.Width(50).Render(r.Name))
	resourceType := lipgloss.JoinHorizontal(lipgloss.Left, resourceMetaTitleStyle.Width(5).Render("Type"), style.Width(10).Render(r.StringType()))
	return lipgloss.JoinHorizontal(lipgloss.Left, name, resourceType)
}

func renderStatus(s *domain.Status, selected bool) string {
	style := statusStyle
	if selected {
		style = hiStatusStyle
	}
	return style.Render(s.Description)
}

func renderStep(s *domain.Step, selected bool) string {
	selectedMarker := " "
	completedMarker := "\uE640"
	style := stepStyle
	if selected {
		style = selectedStepStyle
		selectedMarker = ">"

	}
	if s.Complete {
		completedMarker = "\U000F0856"
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(selectedMarker), style.Render(completedMarker), style.Render(s.Description))
}

func renderHeader(task *domain.Task, headerActive bool) string {
	style := titleStyle
	if headerActive {
		style = activeTitleStyle
	}

	var b strings.Builder
	b.WriteString(style.Render(task.Title) + "\n")
	if reflect.ValueOf(task.Due).IsZero() {
		b.WriteString(fmt.Sprintf("%s  %s\n\n", metaTitleStyle.Render("Due"), titleStyle.Render("None")))
	} else {
		b.WriteString(fmt.Sprintf("%s  %s\n\n", metaTitleStyle.Render("Due"), titleStyle.Render(task.Due.Format("2006/01/02"))))
	}
	return b.String()
}
