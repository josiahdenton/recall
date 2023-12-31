package forms

import (
	"fmt"
	"github.com/josiahdenton/recall/internal/domain"
	styles2 "github.com/josiahdenton/recall/internal/ui/styles"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	title = iota
	due
	priority
)

var (
	dateRe                = regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{2,4}`)
	priorityKeys          = []string{"None", "Low", "High"}
	selectedPriorityStyle = styles2.PrimaryColor.Copy()
	priorityStyle         = styles2.SecondaryGray.Copy()
)

type TaskFormMsg struct {
	Task domain.Task
}

type TaskModel struct {
	inputs        []textinput.Model
	priorityMap   map[string]domain.Priority
	prioriyCursor int
	status        string
	active        int
}

func NewTaskForm() TaskModel {
	inputTitle := textinput.New()
	inputTitle.Focus()
	inputTitle.Width = 60
	inputTitle.CharLimit = 60
	inputTitle.Prompt = "Title: "
	inputTitle.PromptStyle = styles2.FormLabelStyle
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
	inputDue.PromptStyle = styles2.FormLabelStyle
	inputDue.Placeholder = "mm/dd/yyyy"

	inputDue.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 || !dateRe.Match([]byte(s)) {
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

	return TaskModel{
		inputs:      inputs,
		priorityMap: priority,
	}
}

func (m TaskModel) Init() tea.Cmd {
	return nil
}

func (m TaskModel) View() string {
	var b strings.Builder
	b.WriteString(styles2.FormTitleStyle.Render("Add Status"))
	b.WriteString("\n\n")
	b.WriteString(m.inputs[title].View())
	b.WriteString("\n")
	b.WriteString(m.inputs[due].View())
	b.WriteString("\n")
	b.WriteString("Priority: ")
	b.WriteString(fmt.Sprintf("%s %s %s"))
	b.WriteString("\n\n")
	b.WriteString(styles2.FormErrorStyle.Render(m.status))
	return b.String()
}

func (m TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.active < priority {
		m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.active == priority {
			switch msg.String() {
			case "l":
				m.prioriyCursor++
			case "h":
				if m.prioriyCursor > 0 {
					m.prioriyCursor--
				}
			}
		}

		switch msg.Type {
		case tea.KeyEnter:
			if m.active == title {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
				m.inputs[m.active%len(m.inputs)].Focus()
				break
			} else if m.active == due {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
			}

			// TODO fix the <nil>
			if m.inputs[title].Err != nil || m.inputs[due].Err != nil {
				m.status = fmt.Sprintf("%v, %v", m.inputs[title].Err, m.inputs[due].Err)
			} else {
				cmds = append(cmds, addTask(m.inputs[title].Value(), m.inputs[due].Value()))
				m.inputs[title].Reset()
				m.inputs[due].Reset()
			}
		case tea.KeyTab:
			if m.active < priority {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active++
			}
			if m.active < priority {
				m.inputs[m.active%len(m.inputs)].Focus()
			}
		case tea.KeyShiftTab:
			if m.active > 0 {
				m.inputs[m.active%len(m.inputs)].Blur()
				m.active--
				m.inputs[m.active%len(m.inputs)].Focus()
			}
		}
		if len(m.inputs[title].Value()) > 0 || len(m.inputs[due].Value()) > 0 {
			m.status = ""
		}
	}

	return m, tea.Batch(cmds...)
}

func nextInput(current int) int {
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

func addTask(title, due string) tea.Cmd {
	return func() tea.Msg {
		return TaskFormMsg{
			Task: domain.Task{
				Title: title,
				Due:   due,
			},
		}
	}
}
