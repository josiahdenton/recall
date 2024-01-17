package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"strings"
)

type ZettelFormMsg struct {
	Zettel domain.Zettel
}

func NewZettelForm() ZettelFormModel {
	name := textinput.New()
	name.Focus()
	name.Width = 60
	name.CharLimit = 60
	name.Prompt = "Name: "
	name.Placeholder = "..."
	name.PromptStyle = formLabelStyle

	name.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("missing name field")
		}
		return nil
	}
	return ZettelFormModel{
		nameInput: name,
	}
}

type ZettelFormModel struct {
	nameInput textinput.Model
	status    string
}

func (m ZettelFormModel) Init() tea.Cmd {
	return nil
}

func (m ZettelFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Zettel"))
	b.WriteString("\n\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m ZettelFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.nameInput.Reset()
		case tea.KeyEnter:
			if m.nameInput.Err != nil {
				m.status = "missing name for zettel"
				break
			}
			cmds = append(cmds, addZettel(domain.Zettel{Name: m.nameInput.Value()}))
		}
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func addZettel(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Update: zettel,
			Type:   shared.ModifyZettel,
		}
	}
}
