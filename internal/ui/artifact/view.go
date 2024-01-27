package artifact

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
)

func New() Model {
	return Model{}
}

type Model struct {
	artifact  *domain.Artifact
	releases  list.Model
	resources list.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		m.artifact = msg.State.(*domain.Artifact)
		// releases

		// resources
	}

	return m, nil
}

func releasesToItemList(releases []domain.Release) []list.Item {
}

func resourcesToItemList(resources []domain.Resource) []list.Item {
	return nil
}
