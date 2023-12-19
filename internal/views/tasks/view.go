package tasks

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ------------------------------- //
//             styles              //
// ------------------------------- //

var (
	//windowStyle = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder())
	highPriority      = lipgloss.Color("#FF3131")
	mediumPriority    = lipgloss.Color("#FFF01F")
	regularPriority   = lipgloss.Color("#39FF14")
	unStartedPriority = lipgloss.Color("#1F51FF")
)

func New() Model {
	return Model{}
}

type Model struct {
	Tasks []string // for now, just a string
}

type LoadTasks struct {
	Tasks []string
}

func loadTasks() tea.Msg {
	return LoadTasks{[]string{"update pm for eoy", "clean dishes", "update docs"}}
}

func (m Model) Init() tea.Cmd {
	return loadTasks
}

func (m Model) View() string {
	var s string
	for _, task := range m.Tasks {
		s += fmt.Sprintf("[ ] %s\n", task)
	}
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadTasks:
		m.Tasks = msg.Tasks
	}
	return nil, nil
}
