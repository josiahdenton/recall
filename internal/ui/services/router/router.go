package router

import tea "github.com/charmbracelet/bubbletea"

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
	PageCount
	// Forms...
	// TODO - fill in
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
		return Route{}
	}
	return r.History[len(r.History)-1]
}

func (r *Router) Update(msg tea.Msg) tea.Cmd {
	// TODO - fill in
	return nil
}

// Messages

type loadPageMsg struct {
	route Route
}

func GotoPage(route Route) tea.Cmd {
	return func() tea.Msg {
		return loadPageMsg{route: route}
	}
}

func GotoForm() tea.Cmd {
	return func() tea.Msg {
		// TODO - fill in
		return nil
	}
}
