package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/ui/router"
	taskdetailed "github.com/josiahdenton/recall/internal/ui/tasks/detailed"
	tasklist "github.com/josiahdenton/recall/internal/ui/tasks/list"
)

var (
	windowStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func New() Model {
	return Model{
		taskList:     tasklist.New(),
		taskDetailed: taskdetailed.New(),
		page:         router.TaskListPage,
	}
}

type Model struct {
	taskList     tea.Model
	taskDetailed tea.Model
	page         router.Page
	width        int
	height       int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.taskList.Init(), tea.EnterAltScreen)
}

func (m Model) View() string {
	var pageModel tea.Model
	switch m.page {
	case router.TaskListPage:
		pageModel = m.taskList
	case router.TaskDetailedPage:
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
	// TODO move this to the router
	case tasklist.GotoDetailedPageMsg:
		m.page = router.TaskDetailedPage
	case router.GotoPageMsg:
		m.page = msg.Page
	}

	switch m.page {
	case router.TaskListPage:
		m.taskList, cmd = m.taskList.Update(msg)
		cmds = append(cmds, cmd)
	case router.TaskDetailedPage:
		m.taskDetailed, cmd = m.taskDetailed.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// TODO I should have a tick for auto saving my changes
// run after every 5 minutes or something like that
