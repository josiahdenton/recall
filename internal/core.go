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
	"github.com/josiahdenton/recall/internal/ui/zettel"
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
	setupRootZettel(instance)

	return Model{
		taskList:        tasklist.New(),
		taskDetailed:    taskdetailed.New(),
		menu:            menu.New(),
		cycles:          cycles.New(),
		accomplishments: accomplishments.Model{},
		resources:       resources.New(),
		accomplishment:  accomplishment.Model{},
		zettel:          zettel.New(),
		page:            domain.MenuPage,
		repository:      instance,
	}
}

func setupRootZettel(repository repository.Repository) {
	if len(repository.AllZettels()) == 0 {
		// save a new zettel
		repository.ModifyZettel(domain.NewZettel("root"))
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
	repository      repository.Repository
	history         router.History
	page            domain.Page
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	// TODO - use tea.ExecProcess() to pause while editing in nvim
	return tea.Batch(tea.EnterAltScreen)
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
	case domain.ZettelPage:
		pageModel = m.zettel
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
			// TODO may have to "pause" when I open another app inside the terminal
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
	case shared.DeleteStateMsg:
		m.deleteState(msg)
	case shared.RequestStateMsg:
		cmds = append(cmds, m.fetchState(msg))
	case router.GotoPageMsg:
		m.history.Pages = append(m.history.Pages, msg)
		cmds = append(cmds, m.loadPage(msg))
	case router.PreviousPageMsg:
		// previous page is 1 page before current
		//if len(m.history.Pages) > 1 {
		//	last := m.history.Pages[len(m.history.Pages)-2]
		//	m.history.Pages = m.history.Pages[:len(m.history.Pages)-1]
		//	cmds = append(cmds, router.GotoPage(last.Page, last.RequestedItemId))
		//} else {
		//	cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		//}
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
	case domain.ZettelPage:
		m.zettel, cmd = m.zettel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) fetchState(msg shared.RequestStateMsg) tea.Cmd {
	return func() tea.Msg {
		var state any
		switch msg.Type {
		case shared.LoadZettel:
			if msg.ID > 0 {
				state = m.repository.Zettel(msg.ID)
			} else {
				state = m.repository.AllZettels()
			}
		}

		return shared.LoadedStateMsg{State: state}
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
	case shared.ModifyCycle:
		update := msg.Update.(domain.Cycle)
		m.repository.ModifyCycle(update)
	case shared.ModifySettings:
		update := msg.Update.(domain.Settings)
		m.repository.ModifySettings(update)
	case shared.ModifyTask:
		update := msg.Update.(domain.Task)
		m.repository.ModifyTask(update)
	case shared.ModifyAccomplishment:
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
	case shared.ModifyStep:
	case shared.ModifyResource:
	case shared.ModifyStatus:
	case shared.ModifyZettel:
		update := msg.Update.(domain.Zettel)
		m.repository.ModifyZettel(update)
	}
}

func (m Model) deleteState(msg shared.DeleteStateMsg) {
	switch msg.Type {
	case shared.ModifyTask:
		m.repository.DeleteTask(msg.ID)
	}
}
