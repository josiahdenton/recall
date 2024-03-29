package accomplishments

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	selectedAccomplishmentStyle = styles.PrimaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d")).Width(70)
	defaultAccomplishmentStyle  = styles.SecondaryGrayStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(70)
	descriptionStyle            = styles.SecondaryColorStyle.Copy()
	impactStyle                 = styles.PrimaryGrayStyle.Copy()
)

func renderAccomplishment(accomplishment *domain.Accomplishment, selected bool) string {
	style := defaultAccomplishmentStyle
	if selected {
		style = selectedAccomplishmentStyle
	}
	return style.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		descriptionStyle.Render(accomplishment.Description),
		impactStyle.Render(accomplishment.Impact),
	))
}
