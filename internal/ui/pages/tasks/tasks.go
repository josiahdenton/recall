package tasks

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	render "github.com/josiahdenton/recall/internal/ui/renders"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"log"
	"strings"
)

func New() *Model {
	listBoxStyle := styles.Box(styles.BoxOptions{
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: styles.BaseHeight*2 - 10,
		},
		BorderColor: styles.SecondaryGray,
	})
	headerStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryGray,
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: 10,
		},
	})
	return &Model{
		form:        forms.NewTaskForm(),
		listStyle:   listBoxStyle,
		headerStyle: headerStyle,
	}
}

type Model struct {
	tasks       list.Model
	form        forms.Form
	listStyle   lipgloss.Style
	headerStyle lipgloss.Style
	ready       bool
	showForm    bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	// TODO - also render the summary box above??
	var b strings.Builder
	if !m.ready {
		b.WriteString(m.headerStyle.Render(""))
		b.WriteString("\n")
		b.WriteString(m.listStyle.Render(""))
	} else if m.ready && m.showForm {
		// note - forms handle their own styling
		b.WriteString(m.form.View())
	} else {
		selected, ok := m.tasks.SelectedItem().(*domain.Task)
		if ok {
			b.WriteString(m.headerStyle.Render(render.RenderTaskHeader(selected)))
		} else {
			b.WriteString(m.headerStyle.Render("Try adding a task..."))
		}

		b.WriteString("\n")
		b.WriteString(m.listStyle.Render(m.tasks.View()))
	}

	return b.String()
}

func (m *Model) Reset() {
	m.ready = false
	m.showForm = false
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// global
	cmd = m.onGlobalEvents(msg)
	cmds = append(cmds, cmd)

	// input
	if m.ready && !m.showForm {
		cmd = m.onInput(msg)
		cmds = append(cmds, cmd)
	} else if m.ready && m.showForm {
		cmd = m.onFormInput(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) onGlobalEvents(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case forms.TaskFormMsg:
		// todo - add task to list
		// send task off to be saved
	case router.OnInitPageMsg:
		// request to load tasks from DB
		log.Printf("did we get an on init")
		if msg.Page == router.TasksPage {
			return state.Load(state.Request{
				Type: state.Tasks,
			})
		}
	case state.LoadedStateMsg:
		log.Printf("did we get a loaded?")
		// get tasks and create lists from them
		tasks, ok := msg.State.([]domain.Task)
		if !ok {
			return toast.ShowToast("failed to fetch tasks", toast.Warn)
		}
		m.setTasks(tasks)
		m.ready = true
	}

	// no match found
	return nil
}

func (m *Model) onInput(msg tea.Msg) tea.Cmd {
	if m.tasks.FilterState() != list.Filtering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "e":
				if selected, ok := m.tasks.SelectedItem().(*domain.Task); ok {
					m.showForm = true
					return forms.EditTask(selected)
				}
			case "a":
				// cmd to go to form "add task"
				m.showForm = true
			case "esc":
				return router.Back()
			}
		}
	}

	var cmd tea.Cmd
	m.tasks, cmd = m.tasks.Update(msg)
	return cmd
}

func (m *Model) onFormInput(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.showForm = false
			m.form.Reset()
			return nil
		}
	}
	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)
	return cmd
}

func (m *Model) setTasks(tasks []domain.Task) {
	m.tasks = list.New(render.TasksToListItems(tasks), render.TaskDelegate{}, 50, 20)
	m.tasks.Title = "Tasks"
	m.tasks.Styles.PaginationStyle = styles.PaginationStyle
	m.tasks.Styles.Title = styles.PageTitleStyle
	m.tasks.SetShowHelp(false)
	m.tasks.KeyMap.Quit.Unbind()
}
