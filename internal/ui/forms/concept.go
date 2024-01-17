package forms

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type ConceptFormMsg struct {
	Concept string
}

type attachConceptMsg struct {
	concept string
}

func AttachConcept(concept string) tea.Cmd {
	return func() tea.Msg {
		return attachConceptMsg{concept: concept}
	}
}

func NewConceptForm() ConceptFormModel {
	input := textarea.New()
	input.Focus()
	input.MaxWidth = 80
	input.SetWidth(80)
	input.MaxHeight = 20 // TODO - this may change, check zettel/view
	input.SetHeight(20)
	input.CharLimit = 2000

	return ConceptFormModel{
		input: input,
	}
}

type ConceptFormModel struct {
	input  textarea.Model
	status string
	ready  bool
}

func (m ConceptFormModel) Init() tea.Cmd { return nil }

func (m ConceptFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Modify Concept"))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m ConceptFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.ready {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case attachConceptMsg:
		m.input.SetValue(msg.concept)
		m.ready = true
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlS {
			if m.input.Err != nil || len(strings.Trim(m.input.Value(), " \n")) == 0 {
				m.status = errorStyle.Render("missing concept")
			} else {
				cmds = append(cmds, addConcept(m.input.Value()))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func addConcept(concept string) tea.Cmd {
	return func() tea.Msg {
		return ConceptFormMsg{Concept: concept}
	}
}
