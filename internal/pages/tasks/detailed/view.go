package detailed

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/pages/tasks"
	tasklist "github.com/josiahdenton/recall/internal/pages/tasks/list"
	"strings"
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
)

type Model struct {
	ready     bool
	task      *tasks.Task
	steps     list.Model
	resources list.Model
	status    list.Model
}

func (m *Model) Init() tea.Cmd {
	// TODO - this may have to load more?
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	b.WriteString(m.task.Title + "\n")
	b.WriteString(fmt.Sprintf("%s  %s\n\n", m.task.Due, m.task.Runtime()))
	b.WriteString(m.steps.View() + "\n")
	b.WriteString(m.resources.View() + "\n")
	b.WriteString(m.status.View() + "\n")
	return b.String()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tasklist.ShowDetailedMsg:
		m.task = msg.Task
		m.setupLists(msg.Task)
	}
	return m, nil
}

func (m *Model) setupLists(task *tasks.Task) {
	steps := make([]list.Item, len(task.Steps))
	resources := make([]list.Item, len(task.Resources))
	status := make([]list.Item, len(task.Status))
	for i := range task.Steps {
		s := task.Steps[i]
		steps[i] = &s
	}
	for i := range task.Resources {
		r := task.Resources[i]
		resources[i] = &r
	}
	for i := range task.Status {
		s := task.Status[i]
		status[i] = &s
	}

	m.steps = list.New(steps, stepDelegate{}, 50, 5)
	m.steps.Title = "Steps"
	m.steps.SetFilteringEnabled(false)

	m.resources = list.New(resources, resourceDelegate{}, 50, 5)
	m.resources.Title = "Resources"
	m.resources.SetFilteringEnabled(false)

	m.status = list.New(status, statusDelegate{}, 50, 5)
	m.status.Title = "Status"
	m.status.SetFilteringEnabled(false)
}
