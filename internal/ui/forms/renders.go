package forms

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedBorderOptionStyle = styles.PrimaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(25)
	defaultBorderOptionStyle  = styles.SecondaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(25)
	selectedZettelStyle       = styles.SecondaryColorStyle.Copy()
	defaultZettelStyle        = styles.SecondaryGrayStyle.Copy()
	cursorStyle               = styles.PrimaryColorStyle.Copy().Width(2)
)

func renderResourceOption(cycle *createResourceOption, selected bool) string {
	var s string
	if selected {
		s = selectedOptionStyle.Render(cycle.Title)
	} else {
		s = defaultOptionStyle.Render(cycle.Title)
	}
	return s
}

var (
	selectedResourceStyle = styles.AccentColorStyle.Copy().Width(50)
	defaultResourceStyle  = styles.PrimaryGrayStyle.Copy().Width(50)
	titleKeyStyle         = styles.SecondaryGrayStyle.Copy()
	alignStyle            = lipgloss.NewStyle().PaddingLeft(1)
)

func renderResource(resource *domain.Resource, selected bool) string {
	selectedMarker := " "
	style := defaultResourceStyle
	if selected {
		selectedMarker = ">"
		style = selectedResourceStyle
	}
	name := style.Render(resource.Name)
	tags := style.Render(resource.Tags)
	s := lipgloss.JoinHorizontal(lipgloss.Left, name, tags)
	return fmt.Sprintf("%s%s", cursorStyle.Render(selectedMarker), alignStyle.Render(s))
}
