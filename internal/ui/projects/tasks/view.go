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
	"log"
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

type GotoDetailedPageMsg struct {
	Task *domain.Task
}

type loadTasks struct {
	tasks []list.Item
}

func LoadTasks(tasks []domain.Task) tea.Cmd {
	log.Printf("load tasks cmd!")
	return func() tea.Msg {
		log.Printf("load tasks cmd running!!!")
		items := make([]list.Item, len(tasks))
		for i := range tasks {
			item := &tasks[i]
			items[i] = item
		}
		return loadTasks{tasks: items}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var s string
	if m.showForm {
		s = styles.WindowStyle.Render(m.forms[m.activeForm].View())
	} else {
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
	case loadTasks:
		// TODO split this setup into its own func
		m.tasks = list.New(msg.tasks, taskDelegate{}, 50, 20)
		m.tasks.SetShowStatusBar(false)
		m.tasks.SetFilteringEnabled(false)
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
				item := m.tasks.SelectedItem().(*domain.Task)
				cmd = router.GotoPage(domain.TaskDetailedPage, item, item.Id)
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
				cmds = append(cmds, forms.AttachTask(selected.Id))
			}
			// I want the accomplishment to refer to the task completed
			// so I need to send the task down to the accomplishment form
		case "esc":
			if m.showForm {
				m.showForm = false
			} else {
				cmd = router.GotoPage(domain.MenuPage, nil, "")
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
