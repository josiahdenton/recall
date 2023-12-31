package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	name = iota
	source
)

type ResourceFormMsg struct {
	Resource domain.Resource
}

type ResourceModel struct {
	inputs []textinput.Model
	active int

	status string
}

func NewStepResourceForm() ResourceModel {
	inputName := textinput.New()
	inputName.Focus()
	inputName.Width = 60
	inputName.CharLimit = 60
	inputName.Prompt = "Name: "
	inputName.PromptStyle = formLabelStyle
	inputName.Placeholder = "..."

	inputName.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputSource := textinput.New()
	inputSource.Width = 60
	inputSource.CharLimit = 120
	inputSource.Prompt = "Source: "
	inputSource.PromptStyle = formLabelStyle
	inputSource.Placeholder = "..."

	inputSource.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 2)
	inputs[name] = inputName
	inputs[source] = inputSource

	return ResourceModel{
		inputs: inputs,
	}
}

func (m ResourceModel) Init() tea.Cmd {
	return nil
}

func (m ResourceModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Resource"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[name].View())
	b.WriteString("\n")
	b.WriteString(m.inputs[source].View())
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m ResourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active == name {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			// TODO fix the <nil>
			if m.inputs[name].Err != nil || m.inputs[source].Err != nil {
				m.status = errorStyle.Render(fmt.Sprintf("%v, %v", m.inputs[name].Err, m.inputs[source].Err))
			} else {
				cmds = append(cmds, addResource(m.inputs[name].Value(), m.inputs[source].Value()))
				m.inputs[name].Reset()
				m.inputs[source].Reset()
			}
		case tea.KeyTab:
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		case tea.KeyShiftTab:
			if m.active > 0 {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active--
				m.inputs[m.active%len(m.inputs)].Focus()
			}
		}
		if len(m.inputs[name].Value()) > 0 || len(m.inputs[source].Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func addResource(name string, source string) tea.Cmd {
	return func() tea.Msg {
		return ResourceFormMsg{
			Resource: domain.Resource{
				Name:   name,
				Source: source,
			},
		}
	}
}
