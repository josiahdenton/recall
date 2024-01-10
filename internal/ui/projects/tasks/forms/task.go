package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	title = iota
	due
	priority
)

const longDateForm = "Jan 2, 2006 at 3:04pm (MST)"

var (
	priorityKeys = []string{"None", "Low", "High"}
)

type TaskFormModel struct {
	inputs         []textinput.Model
	priorityMap    map[string]domain.Priority
	priorityCursor int
	status         string
	active         int
}

func NewTaskForm() TaskFormModel {
	inputTitle := textinput.New()
	inputTitle.Focus()
	inputTitle.Width = 60
	inputTitle.CharLimit = 60
	inputTitle.Prompt = "Title: "
	inputTitle.PromptStyle = styles.FormLabelStyle
	inputTitle.Placeholder = "..."

	inputTitle.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputDue := textinput.New()
	inputDue.Width = 60
	inputDue.CharLimit = 120
	inputDue.Prompt = "Due: "
	inputDue.PromptStyle = styles.FormLabelStyle
	inputDue.Placeholder = "Jan 5, 2013"

	inputDue.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("step description missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, 2)
	inputs[title] = inputTitle
	inputs[due] = inputDue

	priority := make(map[string]domain.Priority, 3)
	priority[priorityKeys[domain.TaskPriorityNone]] = domain.TaskPriorityNone
	priority[priorityKeys[domain.TaskPriorityLow]] = domain.TaskPriorityLow
	priority[priorityKeys[domain.TaskPriorityHigh]] = domain.TaskPriorityHigh

	return TaskFormModel{
		inputs:         inputs,
		priorityMap:    priority,
		priorityCursor: 0,
	}
}

func (m TaskFormModel) Init() tea.Cmd {
	return nil
}

func (m TaskFormModel) View() string {
	var b strings.Builder
	b.WriteString(styles.FormTitleStyle.Render("Add Task"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[title].View())
	b.WriteString("\n\n")
	b.WriteString(m.inputs[due].View())
	b.WriteString("\n\n")
	b.WriteString(styles.FormLabelStyle.Render("Priority: "))
	b.WriteString(shared.VerticalOptions(priorityKeys, m.priorityCursor))
	b.WriteString("\n\n")
	b.WriteString(styles.FormErrorStyle.Render(m.status))
	return b.String()
}

func (m TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.active < priority {
		m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			if m.active == priority && m.priorityCursor < len(priorityKeys)-1 {
				m.priorityCursor++
			}
		case "h":
			if m.active == priority && m.priorityCursor > 0 {
				m.priorityCursor--
			}
		}

		switch msg.Type {
		case tea.KeyEnter:
			if m.active == title {
				m.inputs[title].Blur()
				m.active++
				m.inputs[due].Focus()
				break
			} else if m.active == due {
				m.inputs[due].Blur()
				m.active++
				break
			}

			// TODO fix the <nil>
			if m.inputs[title].Err != nil || m.inputs[due].Err != nil {
				m.status = fmt.Sprintf("%v, %v", m.inputs[title].Err, m.inputs[due].Err)
			} else {
				date := m.inputs[due].Value()
				t := mustParseDate(date)
				cmds = append(cmds, addTask(m.inputs[title].Value(), t, m.priorityMap[priorityKeys[m.priorityCursor]]))
				m.inputs[title].Reset()
				m.inputs[due].Reset()
				m.priorityCursor = -1
				m.active = 0
				m.inputs[m.active].Focus()
			}
		case tea.KeyTab:
			if m.active < priority {
				m.inputs[m.active].Blur()
				m.active = m.nextInput(m.active)
			}
			if m.active < priority {
				m.inputs[m.active].Focus()
			}
		case tea.KeyShiftTab:
			if m.active > 0 {
				if m.active < priority {
					m.inputs[m.active].Blur()
				}
				m.active--
				if m.active < priority {
					m.inputs[m.active].Focus()
				}
			}
		}
		if len(m.inputs[title].Value()) > 0 || len(m.inputs[due].Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func (m TaskFormModel) nextInput(current int) int {
	switch current {
	case title:
		return due
	case due:
		return priority
	case priority:
		return title
	}
	return title

}

func addTask(title string, due time.Time, priority domain.Priority) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Type:   shared.TaskUpdate,
			Update: domain.NewTask(title, due, priority),
		}
	}
}

// TODO - no longer have any "musts", have it go to a global notify service

func mustParseDate(date string) time.Time {
	input := fmt.Sprintf("%s at 10:00pm (EST)", date)
	t, err := time.Parse(longDateForm, input)
	if err != nil {
		return time.Time{}
	}
	return t
}
