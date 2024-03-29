package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"strings"
	"time"
)

const (
	cycleTitle = iota
	startDate
)

type CycleFormModel struct {
	inputs []textinput.Model
	active int
	status string
}

func NewCycleForm() CycleFormModel {
	inputTitle := textinput.New()
	inputTitle.Focus()
	inputTitle.Width = 60
	inputTitle.CharLimit = 60
	inputTitle.Prompt = "Title: "
	inputTitle.PromptStyle = formLabelStyle
	inputTitle.Placeholder = "..."

	inputTitle.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputStartDate := textinput.New()
	inputStartDate.Width = 60
	inputStartDate.CharLimit = 120
	inputStartDate.Prompt = "Start Date: "
	inputStartDate.PromptStyle = formLabelStyle
	inputStartDate.Placeholder = "Jan 8, 2023"

	inputStartDate.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 2)
	inputs[cycleTitle] = inputTitle
	inputs[startDate] = inputStartDate

	return CycleFormModel{
		inputs: inputs,
	}
}

func (m CycleFormModel) Init() tea.Cmd {
	return nil
}

func (m CycleFormModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Add Performance Cycle"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[cycleTitle].View())
	b.WriteString("\n")
	b.WriteString(m.inputs[startDate].View())
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.status))
	return b.String()
}

func (m CycleFormModel) Update(msg tea.Msg) (CycleFormModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active == cycleTitle {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			}

			// TODO fix the <nil>
			if m.inputs[cycleTitle].Err != nil || m.inputs[startDate].Err != nil {
				m.status = errorStyle.Render(fmt.Sprintf("%v, %v", m.inputs[cycleTitle].Err, m.inputs[startDate].Err))
			} else {
				cmds = append(cmds, addCycle(m.inputs[cycleTitle].Value(), mustParseDate(m.inputs[startDate].Value())))
				m.inputs[cycleTitle].Reset()
				m.inputs[startDate].Reset()
				cmds = append(cmds, router.GotoPage(domain.CyclesPage, 0))
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
		if len(m.inputs[cycleTitle].Value()) > 0 || len(m.inputs[startDate].Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func mustParseDate(date string) time.Time {
	input := fmt.Sprintf("%s at 7:00am (EST)", date)
	t, err := time.Parse(longDateForm, input)
	if err != nil {
		return time.Time{}
	}
	return t
}

func addCycle(title string, start time.Time) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: domain.Cycle{
				Title:     title,
				StartDate: start,
			},
			Type: state.ModifyCycle,
		}
	}
}
