package menu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedMenuOptionStyle = styles.PrimaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(25)
	defaultMenuOptionStyle  = styles.SecondaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(25)
)

func renderMenuOption(cycle *domain.MenuOption, selected bool) string {
	var s string
	if selected {
		s = selectedMenuOptionStyle.Render(cycle.Title)
	} else {
		s = defaultMenuOptionStyle.Render(cycle.Title)
	}
	return s
}
