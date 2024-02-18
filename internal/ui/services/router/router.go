package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	TaskPage = iota
	TasksPage
	ZettelPage
	ZettelsPage
	ResourcesPage
	ResourcePage
	CyclesPage
	CyclePage
	AccomplishmentPage
	ContentForm
	TaskForm
	StepForm
	ResourceForm
	PageCount
	// Forms...
	PageLoading
)

type Page = int

type Route struct {
	Page Page
	ID   uint
}

func New() *Router {
	return &Router{}
}

type Router struct {
	History []Route
}

func (r *Router) Page() Route {
	if len(r.History) == 0 {
		return Route{
			Page: PageLoading,
			ID:   0,
		}
	}
	return r.History[len(r.History)-1]
}

func (r *Router) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case backMsg:
		cmd = r.Back()
		cmds = append(cmds, cmd)
	case loadPageMsg:
		cmd = r.Route(msg)
		cmds = append(cmds, cmd)
	}

	cmd = r.onInput(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (r *Router) onInput(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			// go to tasks page
		case "2":
			// go to ...
		}
	}
	return nil
}

func (r *Router) Route(msg loadPageMsg) tea.Cmd {
	r.History = append(r.History, msg.route)
	return initPage(msg.route)
}

func (r *Router) Back() tea.Cmd {
	if len(r.History) < 2 {
		return nil
	}

	route := r.History[len(r.History)-2]
	r.History = r.History[:len(r.History)-2]
	return GotoPage(route)
}

// Messages

type OnInitPageMsg struct {
	Page Page
	ID   uint
}

func initPage(r Route) tea.Cmd {
	return func() tea.Msg {
		return OnInitPageMsg{
			Page: r.Page,
			ID:   r.ID,
		}
	}
}

type loadPageMsg struct {
	route Route
}

func GotoPage(route Route) tea.Cmd {
	return func() tea.Msg {
		return loadPageMsg{route: route}
	}
}

// GotoForm TODO - will need to add context support
func GotoForm(form Page) tea.Cmd {
	return func() tea.Msg {
		// TODO - fill in
		return loadPageMsg{route: Route{
			Page: form,
			ID:   0,
		}}
	}
}

type backMsg struct{}

func Back() tea.Cmd {
	return func() tea.Msg {
		return backMsg{}
	}
}
