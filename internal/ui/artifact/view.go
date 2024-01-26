package artifact

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
)

func New() Model {
	return Model{}
}

type Model struct {
	artifact  domain.Artifact
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
	return m, nil
}

func releasesToItemList() []list.Item {
	return nil
}

func resourcesToItemList() []list.Item {
	return nil
}
