package status

import "github.com/charmbracelet/lipgloss"

var (
	updateStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1)
)

type Update struct {
	Description string
}

func (u *Update) Render() string {
	return updateStyle.Render(u.Description)
}
