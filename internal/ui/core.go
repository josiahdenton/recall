package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/pages/tasks"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/services/user"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

type service interface {
	Update(msg tea.Msg) tea.Cmd
}

func New(path string) *Model {
	pages := make([]tea.Model, router.PageCount)
	pages[router.TasksPage] = tasks.New()
	pages[router.TaskForm] = forms.NewTaskForm()

	return &Model{
		state:   state.New(path),
		pages:   pages,
		active:  router.TasksPage,
		router:  router.New(),
		toast:   toast.New(),
		effects: user.New(),
	}
}

type Model struct {
	toast   tea.Model
	pages   []tea.Model
	active  router.Page
	router  *router.Router
	state   *state.State
	effects service

	width  int
	height int
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, router.GotoPage(router.Route{
		Page: router.TasksPage,
	}))
}

func (m *Model) View() string {
	// TODO fixme..
	var b strings.Builder
	if m.active != router.PageLoading {
		b.WriteString(styles.CenterStyle.Width(m.width).Height(m.height).Render(m.pages[m.active].View()))
	} else {
		b.WriteString("...")
	}

	return b.String()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// toast
	m.toast, cmd = m.toast.Update(msg)
	cmds = append(cmds, cmd)

	// services run before pages
	cmd = m.router.Update(msg)
	cmds = append(cmds, cmd)

	cmd = m.state.Update(msg)
	cmds = append(cmds, cmd)

	// set active page
	m.active = m.router.Page().Page
	if m.active != router.PageLoading {
		m.pages[m.active], cmd = m.pages[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
