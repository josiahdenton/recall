package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/adapters/repository"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/menu"
	"github.com/josiahdenton/recall/internal/ui/performance/accomplishments"
	"github.com/josiahdenton/recall/internal/ui/performance/cycles"
	taskdetailed "github.com/josiahdenton/recall/internal/ui/projects/task"
	tasklist "github.com/josiahdenton/recall/internal/ui/projects/tasks"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"log"
)

var (
	windowStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func New() Model {
	return Model{
		taskList:        tasklist.New(),
		taskDetailed:    taskdetailed.New(),
		menu:            menu.New(),
		cycles:          cycles.New(),
		accomplishments: accomplishments.Model{},
		page:            domain.MenuPage,
		repository:      repository.NewLocalStorage(),
	}
}

type Model struct {
	taskList        tea.Model
	taskDetailed    tea.Model
	cycles          tea.Model
	menu            tea.Model
	accomplishments tea.Model
	repository      repository.Repository
	page            domain.Page
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.taskList.Init(), tea.EnterAltScreen)
}

func (m Model) View() string {
	var pageModel tea.Model
	switch m.page {
	case domain.TaskListPage:
		pageModel = m.taskList
	case domain.TaskDetailedPage:
		pageModel = m.taskDetailed
	case domain.CyclesPage:
		pageModel = m.cycles
	case domain.MenuPage:
		pageModel = m.menu
	case domain.AccomplishmentsPage:
		pageModel = m.accomplishments
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
	case shared.SaveStateMsg:
		m.updateState(msg)
	case router.GotoPageMsg:
		cmd = m.refreshState(msg)
		cmds = append(cmds, cmd)
		m.page = msg.Page
	}

	// only push events to "in focus" pages
	switch m.page {
	case domain.TaskListPage:
		m.taskList, cmd = m.taskList.Update(msg)
		cmds = append(cmds, cmd)
	case domain.TaskDetailedPage:
		m.taskDetailed, cmd = m.taskDetailed.Update(msg)
		cmds = append(cmds, cmd)
	case domain.CyclesPage:
		m.cycles, cmd = m.cycles.Update(msg)
		cmds = append(cmds, cmd)
	case domain.MenuPage:
		m.menu, cmd = m.menu.Update(msg)
		cmds = append(cmds, cmd)
	case domain.AccomplishmentsPage:
		m.accomplishments, cmd = m.accomplishments.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) refreshState(msg router.GotoPageMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.Page {
	case domain.TaskListPage:
		cmd = tasklist.LoadTasks(m.repository.AllTasks())
	case domain.TaskDetailedPage:
	case domain.CyclesPage:
		cmd = cycles.LoadCycles(m.repository.AllCycles())
	case domain.AccomplishmentsPage:
		cycle := m.repository.Cycle(msg.Id)
		cmd = accomplishments.LoadAccomplishments(msg.Id, cycle.Accomplishments)
	case domain.AccomplishmentPage:
		accomplishment := m.repository.Cycle(msg.Id)
	//case domain.MenuPage:
	default:

	}
	//if msg.Page == domain.TaskListPage {
	//	cmd = tasklist.LoadTasks(m.repository.AllTasks())
	//} else if msg.Page == domain.CyclesPage {
	//	cmd = cycles.LoadCycles(m.repository.AllCycles())
	} else if msg.Page == domain.AccomplishmentsPage {
		cycle := msg.Parameter.(*domain.Cycle)
		cmd = accomplishments.LoadAccomplishments(msg.Id, cycle.Accomplishments)
	}
	return cmd
}

// updateState should only worry about updating the repository
func (m Model) updateState(msg shared.SaveStateMsg) {
	switch msg.Type {
	case shared.TaskUpdate:
		update := msg.Update.(domain.Task)
		m.repository.SaveTask(update)
	case shared.AccomplishmentUpdate:
		allCycles := m.repository.AllCycles()
		update := msg.Update.(domain.Accomplishment)
		log.Printf("accomplishment added %+v", update)
		for _, cycle := range allCycles {
			if cycle.Id == msg.ParentId || (msg.ParentId == "" && cycle.Active) {
				cycle.Accomplishments = append(cycle.Accomplishments, update)
				m.repository.SaveCycle(cycle)
				break
			}
		}
	}
}

// TODO I should have a tick for auto saving my changes
// run after every 5 minutes or something like that
