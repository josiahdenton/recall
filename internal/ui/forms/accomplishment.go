package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"log"
	"strings"
)

const (
	description = iota
	impact
	strength
)

type editAccomplishmentMsg struct {
	accomplishment *domain.Accomplishment
}

func EditAccomplishment(accomplishment *domain.Accomplishment) tea.Cmd {
	return func() tea.Msg {
		return editAccomplishmentMsg{accomplishment: accomplishment}
	}
}

type attachTaskMsg struct {
	Task domain.Task
}

func AttachTask(task domain.Task) tea.Cmd {
	return func() tea.Msg {
		return attachTaskMsg{Task: task}
	}
}

type AccomplishmentFormModel struct {
	inputs         []textinput.Model
	accomplishment *domain.Accomplishment
	active         int
	ready          bool
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
			return fmt.Errorf("accomplishment description missing")
		}
		return nil
	}

	inputImpact := textinput.New()
	inputImpact.Width = 60
	inputImpact.CharLimit = 120
	inputImpact.Prompt = "What impact did you have? "
	inputImpact.PromptStyle = styles.FormLabelStyle
	inputImpact.Placeholder = "..."

	inputImpact.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("accomplishment impact missing")
		}
		return nil
	}

	inputStrength := textinput.New()
	inputStrength.Width = 60
	inputStrength.CharLimit = 120
	inputStrength.Prompt = "What strength did you demonstrate? "
	inputStrength.PromptStyle = styles.FormLabelStyle
	inputStrength.Placeholder = "..."

	inputStrength.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("accomplishment strength missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 3)
	inputs[description] = inputDescription
	inputs[impact] = inputImpact
	inputs[strength] = inputStrength

	return AccomplishmentFormModel{
		inputs:         inputs,
		accomplishment: &domain.Accomplishment{},
	}
}

func (m AccomplishmentFormModel) Init() tea.Cmd {
	return nil
}

func (m AccomplishmentFormModel) View() string {
	var b strings.Builder
	b.WriteString(styles.FormLabelStyle.Render("When you completed "))
	if len(m.accomplishment.Tasks) > 0 {
		b.WriteString(styles.FormTitleStyle.Render(m.accomplishment.Tasks[0].Title))
	}
	b.WriteString("\n\n")
	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m AccomplishmentFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editAccomplishmentMsg:
		m.accomplishment = msg.accomplishment
		m.inputs[description].SetValue(m.accomplishment.Description)
		m.inputs[impact].SetValue(m.accomplishment.Impact)
		m.inputs[strength].SetValue(m.accomplishment.Strength)
		m.ready = true
	case attachTaskMsg:
		m.accomplishment.Tasks = append(m.accomplishment.Tasks, msg.Task)
		cmds = append(cmds, state.RequestState(state.LoadCycle, 0))
	case state.LoadedStateMsg:
		cycles := msg.State.([]domain.Cycle)
		active := false
		for _, cylcle := range cycles {
			if cylcle.Active {
				active = true
			}
		}
		if !active {
			cmds = append(cmds, toast.ShowToast("No active cycle found, please activate one", toast.Warn))
		} else {
			m.ready = true
		}
	case tea.KeyMsg:
		if !m.ready {
			break
		}
		switch msg.Type {
		case tea.KeyEnter:
			if m.active < strength {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			if cmd := validateForm(m.inputs[description].Err, m.inputs[impact].Err, m.inputs[strength].Err); cmd != nil {
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			m.accomplishment.Description = m.inputs[description].Value()
			m.accomplishment.Impact = m.inputs[impact].Value()
			m.accomplishment.Strength = m.inputs[strength].Value()

			cmds = append(cmds, saveAccomplishment(*m.accomplishment), archiveTask(m.accomplishment.Tasks[0]))
			m.inputs[description].Reset()
			m.inputs[impact].Reset()
			m.inputs[strength].Reset()
			m.active = 0
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
	}

	if m.ready {
		m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func archiveTask(task domain.Task) tea.Cmd {
	return func() tea.Msg {
		task.Archive = true
		return state.SaveStateMsg{
			Update: task,
			Type:   state.ModifyTask,
		}
	}
}

func saveAccomplishment(accomplishment domain.Accomplishment) tea.Cmd {
	log.Printf("%+v", accomplishment)
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: accomplishment,
			Type:   state.ModifyAccomplishment,
		}
	}
}
