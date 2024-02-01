package forms

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	title = iota
	due
	tTags
)

var (
	leftPad = lipgloss.NewStyle().PaddingLeft(2)
)

func EditTask(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return editTaskMsg{
			Task: task,
		}
	}
}

type editTaskMsg struct {
	Task *domain.Task
}

type TaskFormModel struct {
	title  string
	inputs []textinput.Model
	active int
	task   *domain.Task
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
			return fmt.Errorf("task title missing")
		}
		return nil
	}

	inputDue := textinput.New()
	inputDue.Width = 60
	inputDue.CharLimit = 120
	inputDue.Prompt = "Due: "
	inputDue.PromptStyle = styles.FormLabelStyle
	inputDue.Placeholder = "Jan 5, 2013 (optional)"

	inputDue.Validate = func(s string) error {
		return nil
	}

	inputTags := textinput.New()
	inputTags.Width = 60
	inputTags.CharLimit = 120
	inputTags.Prompt = "Tags: "
	inputTags.PromptStyle = styles.FormLabelStyle
	inputTags.Placeholder = "(comma seperated list - tags improve search)"

	inputs := make([]textinput.Model, 3)
	inputs[title] = inputTitle
	inputs[due] = inputDue
	inputs[tTags] = inputTags

	return TaskFormModel{
		title:  "Add Task",
		inputs: inputs,
		task:   &domain.Task{},
	}
}

func (m TaskFormModel) Init() tea.Cmd {
	return nil
}

func (m TaskFormModel) View() string {
	var b strings.Builder
	b.WriteString(styles.FormTitleStyle.Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(leftPad.Render(m.inputs[title].View()))
	b.WriteString("\n\n")
	b.WriteString(leftPad.Render(m.inputs[due].View()))
	b.WriteString("\n\n")
	b.WriteString(leftPad.Render(m.inputs[tTags].View()))
	return b.String()
}

func (m TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case editTaskMsg:
		m.task = msg.Task
		m.inputs[title].SetValue(m.task.Title)
		m.inputs[due].SetValue(formatDate(m.task.Due))
		m.inputs[tTags].SetValue(m.task.Tags)
		m.title = "Edit Task"
	case tea.KeyMsg:
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

			if err := m.inputs[title].Err; err != nil {
				cmds = append(cmds, toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn))
				return m, tea.Batch(cmds...)
			}

			date := m.inputs[due].Value()
			dueDate, cmd := parseDate(date)
			if cmd != nil {
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			m.task.Title = m.inputs[title].Value()
			m.task.Due = dueDate
			m.task.Tags = m.inputs[tTags].Value()

			cmds = append(cmds, addTask(m.task), router.RefreshPage())
			// Reset form to default state
			m.inputs[title].Reset()
			m.inputs[due].Reset()
			m.active = 0
			m.inputs[m.active].Focus()
		case tea.KeyTab:
			m.inputs[m.active].Blur()
			m.active = m.nextInput(m.active)
			m.inputs[m.active].Focus()
		}
	}

	m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TaskFormModel) nextInput(current int) int {
	switch current {
	case title:
		return due
	case due:
		return tTags
	case tTags:
		return title
	}
	return title

}

func addTask(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Type:   state.ModifyTask,
			Update: *task,
		}
	}
}

func formatDate(due time.Time) string {
	if reflect.ValueOf(due).IsZero() {
		return ""
	}

	s := due.Format(longDateForm)
	value := strings.Split(s, "at")[0]
	return value
}

func parseDate(date string) (time.Time, tea.Cmd) {
	if len(strings.Trim(date, " \n")) < 1 {
		return time.Time{}, nil
	}

	input := fmt.Sprintf("%s at 10:00pm (EST)", date)
	t, err := time.Parse(longDateForm, input)
	if err != nil {
		return time.Time{}, toast.ShowToast("failed to parse date", toast.Warn)
	}
	return t, nil
}
