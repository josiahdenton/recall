package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type editStepMsg struct {
	step *domain.Step
}

func EditStep(step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return editStepMsg{step: step}
	}
}

type StepFormMsg struct {
	Step domain.Step
	Edit bool
}

type StepModel struct {
	input textinput.Model
	step  *domain.Step
}

func NewStepForm() StepModel {
	input := textinput.New()
	input.Focus()
	input.Width = 60
	input.CharLimit = 120
	input.Prompt = "Description: "
	input.PromptStyle = formLabelStyle
	input.Placeholder = "..."

	input.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	return StepModel{
		input: input,
		step:  &domain.Step{},
	}
}

func (m StepModel) Init() tea.Cmd {
	return nil
}

func (m StepModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Step"))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	return b.String()
}

func (m StepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editStepMsg:
		m.step = msg.step
		m.input.SetValue(msg.step.Description)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Err != nil {
				cmds = append(cmds, toast.ShowToast(fmt.Sprintf("%v", m.input.Err), toast.Warn))
			} else {
				m.step.Description = m.input.Value()
				cmds = append(cmds, addStep(*m.step))
				m.input.Reset()
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func addStep(step domain.Step) tea.Cmd {
	return func() tea.Msg {
		return StepFormMsg{
			Step: step,
			Edit: step.ID != 0,
		}
	}
}
