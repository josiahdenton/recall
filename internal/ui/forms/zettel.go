package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

const (
	zName = iota
	zTags
)

type editZettelMsg struct {
	zettel *domain.Zettel
}

func EditZettel(zettel *domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return editZettelMsg{zettel: zettel}
	}
}

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

	tags := textinput.New()
	tags.Width = 60
	tags.CharLimit = 60
	tags.Prompt = "Tags: "
	tags.Placeholder = "(comma seperated list - tags improve search)"
	tags.PromptStyle = formLabelStyle

	inputs := make([]textinput.Model, 2)
	inputs[zName] = name
	inputs[zTags] = tags

	return ZettelFormModel{
		inputs: inputs,
		active: zName,
		zettel: &domain.Zettel{},
	}
}

type ZettelFormModel struct {
	active int
	inputs []textinput.Model
	zettel *domain.Zettel
}

func (m ZettelFormModel) Init() tea.Cmd {
	return nil
}

func (m ZettelFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Zettel"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[zName].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[zTags].View())
	return b.String()
}

func (m ZettelFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case editZettelMsg:
		m.zettel = msg.zettel
		m.inputs[zName].SetValue(m.zettel.Name)
		m.inputs[zTags].SetValue(m.zettel.Tags)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.inputs[zName].Reset()
			m.inputs[zName].Focus()
			m.inputs[zTags].Reset()
			m.inputs[zTags].Blur()
		case tea.KeyEnter:
			if m.active < zTags {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if err := m.inputs[zName].Err; err != nil {
				cmds = append(cmds, toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn))
				return m, tea.Batch(cmds...)
			}
			m.zettel.Name = m.inputs[zName].Value()
			m.zettel.Tags = m.inputs[zTags].Value()
			cmds = append(cmds, addZettel(*m.zettel))
			m.inputs[zName].Reset()
			m.inputs[zName].Focus()
			m.inputs[zTags].Reset()
			m.inputs[zTags].Blur()
		case tea.KeyTab:
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		}
	}

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func addZettel(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: zettel,
			Type:   state.ModifyZettel,
		}
	}
}
