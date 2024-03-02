package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	taskTitle = iota
	taskTags
	taskDesc
	taskDue
	taskInputCount
)

type editTaskMsg struct {
	task *domain.Task
}

func EditTask(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return editTaskMsg{task: task}
	}
}

func NewTaskForm() *TaskFormModel {
	inputTitle := textinput.New()
	inputTitle.Focus()
	inputTitle.Width = 60
	inputTitle.CharLimit = 60
	inputTitle.Prompt = "Title: "
	inputTitle.PromptStyle = styles.FormLabelStyle
	inputTitle.Placeholder = "..."

	inputTitle.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("task title missing")
		}
		return nil
	}

	inputTags := textinput.New()
	inputTags.Width = 60
	inputTags.CharLimit = 60
	inputTags.Prompt = "Tags: "
	inputTags.PromptStyle = styles.FormLabelStyle
	inputTags.Placeholder = "(comma seperated list - tags improve search!)"

	inputTags.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("task tags missing")
		}
		return nil
	}

	inputDesc := textinput.New()
	inputDesc.Width = 60
	inputDesc.CharLimit = 60
	inputDesc.Prompt = "Description: "
	inputDesc.PromptStyle = styles.FormLabelStyle
	inputDesc.Placeholder = "..."

	inputDesc.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("task description missing")
		}
		return nil
	}

	inputDue := textinput.New()
	inputDue.Width = 60
	inputDue.CharLimit = 60
	inputDue.Prompt = "Due: "
	inputDue.PromptStyle = styles.FormLabelStyle
	inputDue.Placeholder = "Jan 5, 2013 (optional)"

	inputDue.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("task due date missing")
		}
		return nil
	}

	inputs := make([]textinput.Model, taskInputCount)
	inputs[taskTitle] = inputTitle
	inputs[taskTags] = inputTags
	inputs[taskDesc] = inputDesc
	inputs[taskDue] = inputDue

	return &TaskFormModel{
		inputs: inputs,
		boxStyle: styles.Box(styles.BoxOptions{
			Size:        styles.Single,
			BorderColor: styles.SecondaryGray,
		}),
	}
}

type TaskFormModel struct {
	inputs   []textinput.Model
	boxStyle lipgloss.Style
	task     *domain.Task

	active int
}

func (m *TaskFormModel) Init() tea.Cmd {
	return nil
}

func (m *TaskFormModel) View() string {
	var b strings.Builder
	b.WriteString("Add Task\n\n")
	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	return m.boxStyle.Render(b.String())
}

func (m *TaskFormModel) Reset() {
	m.inputs[taskTitle].Reset()
	m.inputs[taskTags].Reset()
	m.inputs[taskDesc].Reset()
	m.inputs[taskDue].Reset()

	m.inputs[taskTitle].Focus()
	m.inputs[taskTags].Blur()
	m.inputs[taskDesc].Blur()
	m.inputs[taskDue].Blur()
	m.active = taskTitle
	m.task = &domain.Task{}
}

func (m *TaskFormModel) Update(msg tea.Msg) (Form, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// local events
	cmd = m.onLocalEvents(msg)
	cmds = append(cmds, cmd)

	// inputs
	cmd = m.onInput(msg)
	cmds = append(cmds, cmd)

	return m, cmd
}

func (m *TaskFormModel) onLocalEvents(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case editTaskMsg:
		m.task = msg.task
		m.inputs[taskTitle].SetValue(m.task.Title)
		m.inputs[taskTags].SetValue(m.task.Tags)
		m.inputs[taskDesc].SetValue(m.task.Description)
		m.inputs[taskDue].SetValue(domain.FormatDate(m.task.Due))
	}
	return nil
}

func (m *TaskFormModel) onInput(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if cmd := validateFrom(m.inputs[taskTitle].Err); cmd != nil {
				return cmd
			}

			// parse form inputs
			if cmd := m.parseFormInputs(); cmd != nil {
				return cmd
			}

			// submit form
			m.Reset()
			cmds = append(cmds, submitTaskForm(m.task), router.Back())
		case "tab":
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		case "esc":
			cmds = append(cmds, router.Back())
			m.Reset()
		}
	}

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *TaskFormModel) parseFormInputs() tea.Cmd {
	m.task.Title = m.inputs[taskTitle].Value()
	m.task.Tags = m.inputs[taskTags].Value()
	m.task.Description = m.inputs[taskDesc].Value()
	due := m.inputs[taskDue].Value()
	parsedDue, err := domain.ParseDate(due)
	if err != nil {
		return toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn)
	}
	m.task.Due = parsedDue

	return nil
}

type TaskFormMsg struct {
	Task domain.Task
}

// submitTaskForm - should send a state.Save
func submitTaskForm(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return state.Save(state.Request{
			State: *task,
			Type:  state.Task,
		})
	}
}
