package tasks

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
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

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		tasks := msg.State.([]domain.Task)
		m.tasks = list.New(toItemList(tasks), taskDelegate{}, 50, 30)
		m.tasks.Title = "Tasks"
		m.tasks.Styles.PaginationStyle = paginationStyle
		m.tasks.Styles.Title = titleStyle
		m.tasks.SetShowHelp(false)
		m.tasks.KeyMap.Quit.Unbind()
		m.ready = true
	case shared.SaveStateMsg:
		m.showForm = false
		if msg.Type == shared.ModifyTask {
			task := msg.Update.(domain.Task)
			m.tasks.InsertItem(len(m.tasks.Items()), &task)
		}
	}

	if m.showForm && m.ready {
		m.forms[m.activeForm], cmd = m.forms[m.activeForm].Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.showForm {
			m.showForm = false
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}

	if m.showForm {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if len(m.tasks.Items()) < 1 {
				return m, tea.Batch(cmds...)
			}
			task := m.tasks.SelectedItem().(*domain.Task)
			cmd = router.GotoPage(domain.TaskDetailedPage, task.ID)
			cmds = append(cmds, cmd)
		}
	}

	if m.tasks.FilterState() != list.Unfiltered {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.showForm = true
			m.activeForm = addForm
		case "c":
			m.showForm = true
			m.activeForm = completeForm
			selected := m.tasks.SelectedItem().(*domain.Task)
			cmds = append(cmds, forms.AttachTask(*selected))
		case "d":
			selected := m.tasks.SelectedItem().(*domain.Task)
			m.tasks.RemoveItem(m.tasks.Index())
			cmds = append(cmds, deleteTask(selected.ID))
		}
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

func deleteTask(id uint) tea.Cmd {
	return func() tea.Msg {
		return shared.DeleteStateMsg{
			Type: shared.ModifyTask,
			ID:   id,
		}
	}
}
