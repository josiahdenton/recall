package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type editStatusMsg struct {
	status *domain.Status
}

func EditStatus(status *domain.Status) tea.Cmd {
	return func() tea.Msg {
		return editStatusMsg{status: status}
	}
}

type StatusFormMsg struct {
	Status domain.Status
	Edit   bool
}

type StatusModel struct {
	input  textarea.Model
	status *domain.Status
}

func NewStatusForm() StatusModel {
	input := textarea.New()
	input.Focus()
	input.MaxWidth = 80
	input.SetWidth(80)

	return StatusModel{
		input:  input,
		status: &domain.Status{},
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
	return b.String()
}

func (m StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editStatusMsg:
		m.status = msg.status
		m.input.SetValue(msg.status.Description)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.input.Err != nil {
				cmds = append(cmds, toast.ShowToast(fmt.Sprintf("%v", m.input.Err), toast.Warn))
			} else {
				m.status.Description = m.input.Value()
				cmds = append(cmds, addStatus(*m.status))
				m.input.Reset()
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func addStatus(status domain.Status) tea.Cmd {
	return func() tea.Msg {
		return StatusFormMsg{
			Status: status,
			Edit:   status.ID != 0,
		}
	}
}
