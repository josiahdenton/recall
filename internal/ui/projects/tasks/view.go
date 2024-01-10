package tasks

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/projects/tasks/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

const (
	addForm = iota
	completeForm
)

type activeForm = int

type Model struct {
	ready      bool
	tasks      list.Model
	forms      []tea.Model
	activeForm activeForm
	showForm   bool
}

func New() Model {
	return Model{
		forms: []tea.Model{
			forms.NewTaskForm(),
			forms.NewAccomplishmentForm(),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var s string
	if m.showForm && m.ready {
		s = styles.WindowStyle.Render(m.forms[m.activeForm].View())
	} else if m.ready {
		s = styles.WindowStyle.Render(m.tasks.View())
	}
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.showForm && m.ready {
		m.forms[m.activeForm], cmd = m.forms[m.activeForm].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		tasks := msg.State.([]domain.Task)
		m.tasks = list.New(toItemList(tasks), taskDelegate{}, 50, 20)
		// TODO - cool, but I need to make no action input during filter
		//m.tasks.SetShowStatusBar(false)
		//m.tasks.SetFilteringEnabled(false)
		m.tasks.Title = "Tasks"
		m.tasks.Styles.PaginationStyle = paginationStyle
		m.tasks.Styles.Title = titleStyle
		m.tasks.SetShowHelp(false)
		m.tasks.KeyMap.Quit.Unbind()
		m.ready = true
	case shared.SaveStateMsg:
		m.showForm = false
		if msg.Type == shared.TaskUpdate {
			task := msg.Update.(domain.Task)
			m.tasks.InsertItem(len(m.tasks.Items()), &task)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.showForm {
				task := m.tasks.SelectedItem().(*domain.Task)
				cmd = router.GotoPage(domain.TaskDetailedPage, task.Id)
				cmds = append(cmds, cmd)
			}
		case "a":
			if !m.showForm {
				m.showForm = true
				m.activeForm = addForm
			}
		case "c":
			if !m.showForm {
				m.showForm = true
				m.activeForm = completeForm
				selected := m.tasks.SelectedItem().(*domain.Task)
				cmds = append(cmds, forms.AttachTask(*selected))
			}
		case "esc":
			if m.showForm {
				m.showForm = false
			} else {
				cmd = router.GotoPage(domain.MenuPage, "")
				cmds = append(cmds, cmd)
			}
		}
	}

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func toItemList(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i := range tasks {
		item := &tasks[i]
		items[i] = item
	}
	return items
}
