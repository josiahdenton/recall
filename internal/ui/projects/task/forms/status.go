package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type StatusFormMsg struct {
	Status domain.Status
}

type StatusModel struct {
	input  textarea.Model
	status string
}

func NewStatusForm() StatusModel {
	input := textarea.New()
	input.Focus()
	input.MaxWidth = 60

	return StatusModel{
		input: input,
	}
}

func (m StatusModel) Init() tea.Cmd {
	return nil
}

func (m StatusModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Status"))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.input.Err != nil {
				m.status = errorStyle.Render(fmt.Sprintf("%v", m.input.Err))
			} else {
				cmds = append(cmds, addStatus(m.input.Value()))
				m.input.Reset()
			}
		}
		if len(m.input.Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func addStatus(s string) tea.Cmd {
	return func() tea.Msg {
		return StatusFormMsg{
			Status: domain.Status{
				Description: s,
			},
		}
	}
}
