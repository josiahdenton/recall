package forms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
)

type editStepMsg struct {
	step *domain.Step
}

func EditStep(step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return editStepMsg{step: step}
	}
}
