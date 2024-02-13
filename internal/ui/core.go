package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

// TODO - should this go in the router?
type page interface {
	tea.Model
	// OnDestroy is called when a page loses focus from the main viewing window
	// this should perform any cleanup tasks so the page can be "revisited" after
	// a successful route change
	OnDestroy()
}

type service interface {
	Update(msg tea.Msg) tea.Cmd
}

func New(path string) Model {
	return Model{
		state:  state.New(path),
		router: router.New(),
		toast:  toast.New(),
	}
}

type Model struct {
	toast  tea.Model
	router service
	state  service
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return styles.WindowStyle.Render()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// toast is always running
	m.toast, cmd = m.toast.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
