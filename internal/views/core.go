package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/views/tasks"
)

// height for each window should be 25
var (
	activeWindowStyle = lipgloss.NewStyle().Width(80).Padding(1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#292524"))
)

func New() Model {
	return Model{
		TaskModel: tasks.New(),
	}
}

type Model struct {
	TaskModel tea.Model
	//windowStyle lipgloss.Style

	width  int
	height int
	// Projects
	// ^ Projects will have Categories
	// tasks should have a "child" model that handles the logic here...
	// child models will use the same methods (
}

func (m Model) Init() tea.Cmd {
	// this call the Init from the active child model
	// I will need to use either tea.Batch or tea.Sequence
	return tea.Batch(m.TaskModel.Init(), tasks.Tick(), tea.EnterAltScreen)
}

func (m Model) View() string {
	// create a header...
	// it should have task info...
	// this might create a lot more garbage collection... may need to fix this
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).Render(activeWindowStyle.Render(m.TaskModel.View()))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		//m.windowStyle = m.windowStyle.Width(msg.Width).Height(msg.Height).Align(lipgloss.Center)
		m.height = msg.Height
		m.width = msg.Width
	}

	// need another switch for "active page"
	m.TaskModel, cmd = m.TaskModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
