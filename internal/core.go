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
	"github.com/josiahdenton/recall/internal/ui/resources"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	taskdetailed "github.com/josiahdenton/recall/internal/ui/task"
	tasklist "github.com/josiahdenton/recall/internal/ui/tasks"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"github.com/josiahdenton/recall/internal/ui/zettel"
	"github.com/josiahdenton/recall/internal/ui/zettels"
	"log"
	"os"
	"strings"
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

	// zettel will still be files within recall-notes...
	parentPath := fmt.Sprintf("%s/%s", home, "recall-notes")
	path := fmt.Sprintf("%s/%s", parentPath, ".recall.db")

	instance, err := repository.NewGormInstance(fmt.Sprintf(path))
	if err != nil {
		log.Printf("failed to create sqlite connection")
		os.Exit(1)
	}

	// check to make sure root zettel exists, if not, create the root
	err = instance.LoadRepository()
	if err != nil {
		log.Printf("failed to load repository %v", err)
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
		zettel:          zettel.New(),
		zettels:         zettels.New(),
		toast:           toast.New(),
		page:            domain.MenuPage,
		repository:      instance,
		routeHistory: router.History{
			Pages: []router.GotoPageMsg{{Page: domain.MenuPage, RequestedItemId: 0}},
		},
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
	zettel          tea.Model
	zettels         tea.Model
	toast           tea.Model
	repository      repository.Repository
	routeHistory    router.History
	stateHistory    state.History
	page            domain.Page
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m Model) View() string {
	var page tea.Model
	switch m.page {
	case domain.TaskListPage:
		page = m.taskList
	case domain.TaskDetailedPage:
		page = m.taskDetailed
	case domain.CyclesPage:
		page = m.cycles
	case domain.MenuPage:
		page = m.menu
	case domain.AccomplishmentsPage:
		page = m.accomplishments
	case domain.ResourcesPage:
		page = m.resources
	case domain.AccomplishmentPage:
		page = m.accomplishment
	case domain.ZettelPage:
		page = m.zettel
	case domain.ZettelsPage:
		page = m.zettels
	}
	var b strings.Builder
	b.WriteString(page.View())
	b.WriteString("\n")
	b.WriteString(m.toast.View())
	return windowStyle.Width(m.width).Height(m.height).Render(b.String())
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
	case state.SaveStateMsg:
		m.updateState(msg)
	case state.DeleteStateMsg:
		m.deleteState(msg)
		m.stateHistory.Deletes = append(m.stateHistory.Deletes, msg)
	case state.UndoDeleteStateMsg:
		//TODO: if I undo a in the zettel page after deleting a task, it will bring that task back.
		// I want add the modifyType. Maybe make Deletes a map of types to slices?
		if len(m.stateHistory.Deletes) > 0 {
			previousDelete := m.stateHistory.Deletes[len(m.stateHistory.Deletes)-1]
			m.stateHistory.Deletes = m.stateHistory.Deletes[:len(m.stateHistory.Deletes)-1]
			m.undoDeleteState(previousDelete)
		}
		cmds = append(cmds, router.RefreshPage())
	case state.RequestStateMsg:
		cmds = append(cmds, m.fetchState(msg))
	case router.GotoPageMsg:
		m.routeHistory.Pages = append(m.routeHistory.Pages, msg)
		cmds = append(cmds, m.loadPage(msg))
	case router.PreviousPageMsg:
		previousPage := m.routeHistory.Pages[len(m.routeHistory.Pages)-2]
		m.routeHistory.Pages = m.routeHistory.Pages[:len(m.routeHistory.Pages)-2]
		cmds = append(cmds, router.GotoPage(previousPage.Page, previousPage.RequestedItemId))
	case router.RefreshPageMsg:
		previousPage := m.routeHistory.Pages[len(m.routeHistory.Pages)-1]
		m.routeHistory.Pages = m.routeHistory.Pages[:len(m.routeHistory.Pages)-1]
		cmds = append(cmds, router.GotoPage(previousPage.Page, previousPage.RequestedItemId))
	case router.LoadPageMsg:
		m.page = msg.Page
	}

	// toast is always active
	m.toast, cmd = m.toast.Update(msg)
	cmds = append(cmds, cmd)

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
	case domain.ZettelPage:
		m.zettel, cmd = m.zettel.Update(msg)
		cmds = append(cmds, cmd)
	case domain.ZettelsPage:
		m.zettels, cmd = m.zettels.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) fetchState(msg state.RequestStateMsg) tea.Cmd {
	return func() tea.Msg {
		var s any
		switch msg.Type {
		case state.LoadZettel:
			if msg.ID > 0 {
				s = m.repository.Zettel(msg.ID)
			} else {
				s = m.repository.AllZettels()
			}
		case state.LoadResource:
			s = m.repository.AllResources()
		}

		return state.LoadedStateMsg{State: s}
	}
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
			state = m.repository.Cycle(msg.RequestedItemId)
		case domain.AccomplishmentPage:
			state = m.repository.Accomplishment(msg.RequestedItemId)
		case domain.MenuPage:
			// no state attached...
			// TODO - have the repository read the settings file to determine this
		case domain.ResourcesPage:
			state = m.repository.AllResources()
		case domain.ZettelPage:
			state = m.repository.Zettel(msg.RequestedItemId)
		case domain.ZettelsPage:
			state = m.repository.AllZettels()
		}
		return router.LoadPageMsg{
			Page:  msg.Page,
			State: state,
		}
	}
}

// updateState should only worry about updating the repository
func (m Model) updateState(msg state.SaveStateMsg) {
	switch msg.Type {
	case state.ModifyCycle:
		update := msg.Update.(domain.Cycle)
		m.repository.ModifyCycle(update)
	case state.ModifySettings:
		update := msg.Update.(domain.Settings)
		m.repository.ModifySettings(update)
	case state.ModifyTask:
		update := msg.Update.(domain.Task)
		m.repository.ModifyTask(update)
	case state.ModifyAccomplishment:
		allCycles := m.repository.AllCycles()
		update := msg.Update.(domain.Accomplishment)
		activeSet := false
		for _, cycle := range allCycles {
			if cycle.Active {
				cycle.Accomplishments = append(cycle.Accomplishments, update)
				m.repository.ModifyCycle(cycle)
				activeSet = true
				break
			}
		}
		if !activeSet {
			m.repository.ModifyAccomplishment(update)
		}
	case state.ModifyStep:
		update := msg.Update.(domain.Step)
		m.repository.ModifyStep(update)
	case state.ModifyResource:
	case state.ModifyStatus:
	case state.ModifyZettel:
		update := msg.Update.(domain.Zettel)
		m.repository.ModifyZettel(update)
	}
}

func (m Model) deleteState(msg state.DeleteStateMsg) {
	switch msg.Type {
	case state.ModifyTask:
		m.repository.DeleteTask(msg.ID)
	case state.ModifyStep:
		m.repository.DeleteTaskStep(msg.Parent.(*domain.Task), msg.Child.(*domain.Step))
	case state.ModifyResource:
		m.repository.DeleteTaskResource(msg.Parent.(*domain.Task), msg.Child.(*domain.Resource))
	case state.ModifyStatus:
		m.repository.DeleteTaskStatus(msg.Parent.(*domain.Task), msg.Child.(*domain.Status))
	case state.ModifyCycle:
	case state.ModifyZettel:
		m.repository.DeleteZettel(msg.ID)
	case state.ModifyLink:
		m.repository.UnlinkZettel(msg.Parent.(*domain.Zettel), msg.Child.(*domain.Zettel))
	case state.ModifyAccomplishment:
		m.repository.DeleteAccomplishment(msg.ID)
	case state.ModifySettings:
	}
}

func (m Model) undoDeleteState(msg state.DeleteStateMsg) {
	switch msg.Type {
	case state.ModifyTask:
		m.repository.UndoDeleteTask(msg.ID)
	case state.ModifyAccomplishment:
		m.repository.UndoDeleteAccomplishment(msg.ID)
	case state.ModifyZettel:
		m.repository.UndoDeleteZettel(msg.ID)
	case state.ModifyStep:
	case state.ModifyResource:
	case state.ModifyStatus:
	case state.ModifyCycle:
	case state.ModifyLink:
	case state.ModifySettings:
	}
}
