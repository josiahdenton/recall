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
		loading:     true,
		listStyle:   listBoxStyle,
		headerStyle: headerStyle,
	}
}

type Model struct {
	tasks       list.Model
	loading     bool
	mode        state.Mode
	listStyle   lipgloss.Style
	headerStyle lipgloss.Style
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	// TODO - also render the summary box above??
	var b strings.Builder
	if m.loading {
		b.WriteString(m.headerStyle.Render(""))
		b.WriteString("\n")
		b.WriteString(m.listStyle.Render(""))
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
	m.loading = true
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
	if m.mode == state.View && !m.loading {
		cmd = m.onInput(msg)
		cmds = append(cmds, cmd)
	}

	// form input

	return m, tea.Batch(cmds...)
}

func (m *Model) onGlobalEvents(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case forms.TaskFormMsg:
		// todo - add task to list
		// send task off to be saved
	case state.ModeSwitchMsg:
		m.mode = msg.Current
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
		m.loading = false
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
				return nil
			case "a":
				// cmd to go to form "add task"
				return router.GotoForm(router.TaskForm)
			case "esc":
				return router.Back()
			}
		}
	}

	var cmd tea.Cmd
	m.tasks, cmd = m.tasks.Update(msg)
	return cmd
}

//func (m *Model) onLocalEvents()

func (m *Model) setTasks(tasks []domain.Task) {
	m.tasks = list.New(render.TasksToListItems(tasks), render.TaskDelegate{}, 50, 20)
	m.tasks.Title = "Tasks"
	m.tasks.Styles.PaginationStyle = styles.PaginationStyle
	m.tasks.Styles.Title = styles.PageTitleStyle
	m.tasks.SetShowHelp(false)
	m.tasks.KeyMap.Quit.Unbind()
	m.loading = false
}
