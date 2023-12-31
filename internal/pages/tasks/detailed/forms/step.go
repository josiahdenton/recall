package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/pages/tasks"
)

type StepFormMsg struct {
	Step tasks.Step
}

type StepModel struct {
	input  textinput.Model
	status string
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
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m StepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Err != nil {
				m.status = errorStyle.Render(fmt.Sprintf("%v", m.input.Err))
			} else {
				cmds = append(cmds, addStep(m.input.Value()))
				m.input.Reset()
			}
		}
		if len(m.input.Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func addStep(s string) tea.Cmd {
	return func() tea.Msg {
		return StepFormMsg{
			Step: tasks.Step{
				Description: s,
			},
		}
	}
}
