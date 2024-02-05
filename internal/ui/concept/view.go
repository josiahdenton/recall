package concept

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

func New() Model {
	return Model{
		vp: viewport.New(100, 40),
	}
}

type Model struct {
	vp viewport.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return styles.WindowStyle.Render(m.vp.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		zettel := msg.State.(*domain.Zettel)
		m.vp.SetContent(zettel.Concept)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPreviousPage())
		}
	}

	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
