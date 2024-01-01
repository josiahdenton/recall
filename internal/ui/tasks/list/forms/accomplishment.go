package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	description = iota
	impact
	strength
)

type AccomplishmentFormMsg struct {
	Accomplishment domain.Accomplishment
}

type AccomplishmentFormModel struct {
	inputs []textinput.Model
	status string
	active int
}

func NewAccomplishmentForm() AccomplishmentFormModel {
	inputDescription := textinput.New()
	inputDescription.Focus()
	inputDescription.Width = 60
	inputDescription.CharLimit = 60
	inputDescription.Prompt = "What did you accomplish? "
	inputDescription.PromptStyle = styles.FormLabelStyle
	inputDescription.Placeholder = "..."

	inputDescription.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputImpact := textinput.New()
	inputImpact.Width = 60
	inputImpact.CharLimit = 120
	inputImpact.Prompt = "What impact did you have? "
	inputImpact.PromptStyle = styles.FormLabelStyle
	inputImpact.Placeholder = "mm/dd/yyyy"

	inputImpact.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputStrength := textinput.New()
	inputStrength.Width = 60
	inputStrength.CharLimit = 120
	inputStrength.Prompt = "What strength did you demonstrate?"
	inputStrength.PromptStyle = styles.FormLabelStyle
	inputStrength.Placeholder = "mm/dd/yyyy"

	inputStrength.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 3)
	inputs[description] = inputDescription
	inputs[impact] = inputImpact
	inputs[strength] = inputStrength

	return AccomplishmentFormModel{
		inputs: inputs,
	}
}

func (m AccomplishmentFormModel) Init() tea.Cmd {
	return nil
}

func (m AccomplishmentFormModel) View() string {
	var b strings.Builder
	b.WriteString(styles.FormTitleStyle.Render("Complete Task"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[description].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[impact].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[strength].View())
	b.WriteString("\n\n")
	b.WriteString(styles.FormErrorStyle.Render(m.status))
	return b.String()
}

func (m AccomplishmentFormModel) Update(msg tea.Msg) (AccomplishmentFormModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.active < priority {
		m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active < strength {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			// TODO fix the <nil>
			// TODO care about strength
			if m.inputs[description].Err != nil || m.inputs[impact].Err != nil {
				m.status = fmt.Sprintf("%v, %v", m.inputs[title].Err, m.inputs[due].Err)
			} else {
				cmds = append(cmds, addAccomplishment(m.inputs[description].Value(), m.inputs[impact].Value(), m.inputs[strength].Value()))
				m.inputs[title].Reset()
				m.inputs[due].Reset()
				m.active = 0
			}
		case tea.KeyTab:
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		case tea.KeyShiftTab:
			if m.active > 0 {
				m.inputs[m.active].Blur()
				m.active--
				m.inputs[m.active].Focus()
			}
		}
		if len(m.inputs[description].Value()) > 0 || len(m.inputs[impact].Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func addAccomplishment(description, impact, strength string) tea.Cmd {
	return func() tea.Msg {
		return AccomplishmentFormMsg{
			Accomplishment: domain.Accomplishment{
				Description: description,
				Impact:      impact,
				Strength:    strength,
			},
		}
	}
}
