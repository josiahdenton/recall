package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/views/tasks"
)

var (
	activeWindowStyle = lipgloss.NewStyle().PaddingLeft(4).PaddingTop(2)
)

func New() Model {
	return Model{
		TaskModel: tasks.New(),
	}
}

type Model struct {
	TaskModel tea.Model
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
	return activeWindowStyle.Render(m.TaskModel.View())
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
	}

	// need another switch for "active page"
	m.TaskModel, cmd = m.TaskModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
