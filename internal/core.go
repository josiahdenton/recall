package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	taskdetailed "github.com/josiahdenton/recall/internal/pages/tasks/detailed"
	tasklist "github.com/josiahdenton/recall/internal/pages/tasks/list"
)

const (
	TaskList = iota
	TaskDetailed
)

var (
	windowStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

type Page = int

func New() Model {
	return Model{
		taskList:     &tasklist.Model{},
		taskDetailed: &taskdetailed.Model{},
		page:         TaskList,
	}
}

type Model struct {
	taskList     tea.Model
	taskDetailed tea.Model
	page         Page
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
	var style lipgloss.Style
	switch m.page {
	case TaskList:
		pageModel = m.taskList
		style = windowStyle.Width(m.width).Height(m.height)
	case TaskDetailed:
		pageModel = m.taskDetailed
	}
	return style.Render(pageModel.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tasklist.ShowDetailedMsg:
		m.page = TaskDetailed
	}

	switch m.page {
	case TaskList:
		m.taskList, cmd = m.taskList.Update(msg)
		cmds = append(cmds, cmd)
	case TaskDetailed:
		m.taskDetailed, cmd = m.taskDetailed.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
