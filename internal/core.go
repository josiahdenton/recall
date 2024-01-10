package internal

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/adapters/repository"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/menu"
	"github.com/josiahdenton/recall/internal/ui/performance/accomplishment"
	"github.com/josiahdenton/recall/internal/ui/performance/accomplishments"
	"github.com/josiahdenton/recall/internal/ui/performance/cycles"
	taskdetailed "github.com/josiahdenton/recall/internal/ui/projects/task"
	tasklist "github.com/josiahdenton/recall/internal/ui/projects/tasks"
	"github.com/josiahdenton/recall/internal/ui/resources"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"log"
	"os"
)

var (
	windowStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func New() Model {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get home dir: %v", err)
		os.Exit(1)
	}

	return Model{
		taskList:        tasklist.New(),
		taskDetailed:    taskdetailed.New(),
		menu:            menu.New(),
		cycles:          cycles.New(),
		accomplishments: accomplishments.Model{},
		resources:       resources.New(),
		accomplishment:  accomplishment.Model{},
		page:            domain.MenuPage,
		repository:      repository.NewFileStorage(fmt.Sprintf("%s/%s", home, "recall-notes")),
	}
}

type Model struct {
	taskList        tea.Model
	taskDetailed    tea.Model
	cycles          tea.Model
	menu            tea.Model
	accomplishments tea.Model
	accomplishment  tea.Model
	resources       tea.Model
	repository      repository.Repository
	page            domain.Page
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, loadRepository())
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
	case domain.ResourcesPage:
		pageModel = m.resources
	case domain.AccomplishmentPage:
		pageModel = m.accomplishment
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
	case shared.LoadRepositoryMsg:
		err := m.repository.LoadRepository()
		// TODO - have err message go to global status message handler
		if err != nil {
			log.Printf("failed to LoadRepository: %v", err)
		}
	case shared.SaveStateMsg:
		m.updateState(msg)
		err := m.repository.SaveChanges()
		if err != nil {
			log.Printf("failed saving changes: %v", err)
		}
	case router.GotoPageMsg:
		cmds = append(cmds, m.loadPage(msg))
	case router.LoadPageMsg:
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
	case domain.AccomplishmentPage:
		m.accomplishment, cmd = m.accomplishment.Update(msg)
		cmds = append(cmds, cmd)
	case domain.ResourcesPage:
		m.resources, cmd = m.resources.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) loadPage(msg router.GotoPageMsg) tea.Cmd {
	return func() tea.Msg {
		var state any
		switch msg.Page {
		case domain.TaskListPage:
			state = m.repository.AllTasks()
		// TODO add a page for task archive
		case domain.CyclesPage:
			state = m.repository.AllCycles()
		case domain.TaskDetailedPage:
			state = m.repository.Task(msg.RequestedItemId)
		case domain.AccomplishmentsPage:
			c := m.repository.Cycle(msg.RequestedItemId)
			c.AttachAccomplishments(m.repository.LinkedAccomplishments(c.AccomplishmentIds))
			state = c
		case domain.AccomplishmentPage:
			a := m.repository.Accomplishment(msg.RequestedItemId)
			a.AttachAssociatedTasks(m.repository.LinkedTasks(a.AssociatedTaskIds))
			state = a
		case domain.MenuPage:
			// no state attached...
			// TODO - have the repository read the settings file to determine this
		case domain.ResourcesPage:
			state = m.repository.AllResources()
		}
		return router.LoadPageMsg{
			Page:  msg.Page,
			State: state,
		}
	}
}

// updateState should only worry about updating the repository
func (m Model) updateState(msg shared.SaveStateMsg) {
	switch msg.Type {
	case shared.CycleUpdate:
		update := msg.Update.(domain.Cycle)
		m.repository.SaveCycle(update)
	case shared.SettingsUpdate:
		update := msg.Update.(domain.Settings)
		m.repository.SaveSettings(update)
	case shared.TaskUpdate:
		update := msg.Update.(domain.Task)
		m.repository.SaveTask(update)
	case shared.AccomplishmentUpdate:
		// TODO - this does not handle the replace case
		allCycles := m.repository.AllCycles()
		update := msg.Update.(domain.Accomplishment)
		log.Printf("accomplishment added %+v", update)
		for _, cycle := range allCycles {
			if cycle.Id == msg.ParentId || (msg.ParentId == "" && cycle.Active) {
				cycle.AccomplishmentIds = append(cycle.AccomplishmentIds, update.Id)
				log.Printf("appended stuff, haha: %+v", update)
				cycle.AttachAccomplishments(append(cycle.Accomplishments(), update)) // TODO is this necessary?
				m.repository.SaveCycle(cycle)
				m.repository.SaveAccomplishment(update)
				break
			}
		}
	case shared.StepUpdate:
	case shared.ResourceUpdate:
		update := msg.Update.(domain.Resource)
		log.Printf("saving resource: %+v", update)
		m.repository.SaveResource(update)
	case shared.StatusUpdate:
	}
}

func loadRepository() tea.Cmd {
	return func() tea.Msg {
		return shared.LoadRepositoryMsg{}
	}
}
