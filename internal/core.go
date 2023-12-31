package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/pages"
	taskdetailed "github.com/josiahdenton/recall/internal/pages/tasks/detailed"
	tasklist "github.com/josiahdenton/recall/internal/pages/tasks/list"
)

var (
	windowStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func New() Model {
	return Model{
		taskList:     &tasklist.Model{},
		taskDetailed: taskdetailed.New(),
		page:         pages.TaskList,
	}
}

type Model struct {
	taskList     tea.Model
	taskDetailed tea.Model
	page         pages.Page
	width        int
	height       int
}

func (m Model) Init() tea.Cmd {
	// this call the Init from the active child model
	// I will need to use either tea.Batch or tea.Sequence
	// TODO add tick back for timer
	return tea.Batch(m.taskList.Init(), tea.EnterAltScreen)
}

func (m Model) View() string {
	var pageModel tea.Model
	switch m.page {
	case pages.TaskList:
		pageModel = m.taskList
	case pages.TaskDetailed:
		pageModel = m.taskDetailed
	}
	return windowStyle.Width(m.width).Height(m.height).Render(pageModel.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// TODO before quitting repository will need to save all changes
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tasklist.ShowDetailedMsg:
		m.page = pages.TaskDetailed
	case pages.GotoPageMsg:
		m.page = msg.Page
	}

	switch m.page {
	case pages.TaskList:
		m.taskList, cmd = m.taskList.Update(msg)
		cmds = append(cmds, cmd)
	case pages.TaskDetailed:
		m.taskDetailed, cmd = m.taskDetailed.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// TODO I should have a tick for auto saving my changes
