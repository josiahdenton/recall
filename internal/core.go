package internal

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/adapters/repository"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/artifact"
	"github.com/josiahdenton/recall/internal/ui/artifacts"
	"github.com/josiahdenton/recall/internal/ui/concept"
	"github.com/josiahdenton/recall/internal/ui/keybindings"
	"github.com/josiahdenton/recall/internal/ui/menu"
	"github.com/josiahdenton/recall/internal/ui/performance/accomplishment"
	"github.com/josiahdenton/recall/internal/ui/performance/accomplishments"
	"github.com/josiahdenton/recall/internal/ui/performance/cycles"
	"github.com/josiahdenton/recall/internal/ui/resources"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
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
	windowStyle       = lipgloss.NewStyle().Align(lipgloss.Center)
	shortcutKeyStyle  = styles.SecondaryGray.Copy().Bold(true)
	shortcutDescStyle = styles.SecondaryGray.Copy()
)

func New() Model {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get home dir: %v", err)
		os.Exit(1)
	}

	parentPath := fmt.Sprintf("%s/%s", home, ".recall")
	// check if Dir exists, if not, create it
	fi, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(parentPath, 0775)
		if err != nil {
			log.Printf("failed to mkdir for recall: %v", err)
			os.Exit(1)
		}
		return New()
	} else if err != nil {
		log.Printf("failed to stat recall dir: %v", err)
		os.Exit(1)
	}

	if !fi.IsDir() {
		log.Printf(".recall exists and is not a dir")
		os.Exit(1)
	}

	path := fmt.Sprintf("%s/%s", parentPath, "recall.db")

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

	// TODO - integrate with keyset to manage key binds
	//keys := keyset.New(parentPath)
	//keyBinds, err := keys.Load()
	//keybindings.New()
	keyBinds := domain.DefaultKeybindings()

	return Model{
		taskList:        tasklist.New(keyBinds),
		taskDetailed:    taskdetailed.New(keyBinds),
		menu:            menu.New(keyBinds),
		cycles:          cycles.New(keyBinds),
		accomplishments: accomplishments.New(keyBinds),
		resources:       resources.New(keyBinds),
		accomplishment:  accomplishment.New(keyBinds),
		zettel:          zettel.New(keyBinds),
		zettels:         zettels.New(keyBinds),
		artifact:        artifact.New(keyBinds),
		artifacts:       artifacts.New(keyBinds),
		keybindings:     keybindings.New(keyBinds),
		content:         concept.New(),
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
	artifact        tea.Model
	artifacts       tea.Model
	toast           tea.Model
	keybindings     tea.Model
	content         tea.Model
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
	case domain.ArtifactsPage:
		page = m.artifacts
	case domain.ArtifactPage:
		page = m.artifact
	case domain.KeyBindsPage:
		page = m.keybindings
	case domain.ConceptPage:
		page = m.content
	}
	var b strings.Builder
	b.WriteString(page.View())
	b.WriteString("\n")
	notification := m.toast.View()
	if len(notification) > 0 {
		b.WriteString(notification)
	} else {
		b.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Left,
			shortcutKeyStyle.Render("?"),
			shortcutDescStyle.Render(" - keybindings")))
	}
	return windowStyle.Width(m.width).Height(m.height).Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.String() == "?" && m.page != domain.KeyBindsPage {
			cmds = append(cmds, router.GotoPage(domain.KeyBindsPage, 0))
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

	// TODO - add another always active page for keybindings

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
	case domain.ArtifactsPage:
		m.artifacts, cmd = m.artifacts.Update(msg)
		cmds = append(cmds, cmd)
	case domain.ArtifactPage:
		m.artifact, cmd = m.artifact.Update(msg)
		cmds = append(cmds, cmd)
	case domain.KeyBindsPage:
		m.keybindings, cmd = m.keybindings.Update(msg)
		cmds = append(cmds, cmd)
	case domain.ConceptPage:
		m.content, cmd = m.content.Update(msg)
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
		case state.LoadCycle:
			s = m.repository.AllCycles()
		}

		return state.LoadedStateMsg{State: s}
	}
}

func (m Model) loadPage(msg router.GotoPageMsg) tea.Cmd {
	return func() tea.Msg {
		var s any
		switch msg.Page {
		case domain.TaskListPage:
			s = m.repository.AllTasks()
		// TODO add a page for task archive
		case domain.CyclesPage:
			s = m.repository.AllCycles()
		case domain.TaskDetailedPage:
			s = m.repository.Task(msg.RequestedItemId)
		case domain.AccomplishmentsPage:
			s = m.repository.Cycle(msg.RequestedItemId)
		case domain.AccomplishmentPage:
			s = m.repository.Accomplishment(msg.RequestedItemId)
		case domain.MenuPage:
			// no s attached...
			// TODO - have the repository read the settings file to determine this
		case domain.ResourcesPage:
			s = m.repository.AllResources()
		case domain.ZettelPage:
			s = m.repository.Zettel(msg.RequestedItemId)
		case domain.ZettelsPage:
			s = m.repository.AllZettels()
		case domain.ArtifactsPage:
			s = m.repository.AllArtifacts()
		case domain.ArtifactPage:
			s = m.repository.Artifact(msg.RequestedItemId)
		case domain.KeyBindsPage:
		case domain.ConceptPage:
			s = m.repository.Zettel(msg.RequestedItemId)
		}
		return router.LoadPageMsg{
			Page:  msg.Page,
			State: s,
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
		if update.ID == 0 {
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
		} else {
			m.repository.ModifyAccomplishment(update)
		}
	case state.ModifyStep:
		update := msg.Update.(domain.Step)
		m.repository.ModifyStep(update)
	case state.ModifyResource:
		update := msg.Update.(domain.Resource)
		m.repository.ModifyResource(update)
	case state.ModifyStatus:
		update := msg.Update.(domain.Status)
		m.repository.ModifyStatus(update)
	case state.ModifyZettel:
		update := msg.Update.(domain.Zettel)
		m.repository.ModifyZettel(update)
	case state.ModifyLink:
	case state.ModifyArtifact:
		update := msg.Update.(domain.Artifact)
		m.repository.ModifyArtifact(update)
	case state.ModifyRelease:
		update := msg.Update.(domain.Release)
		m.repository.ModifyRelease(update)
	}
}

func (m Model) deleteState(msg state.DeleteStateMsg) {
	switch msg.Type {
	case state.ModifyTask:
		m.repository.DeleteTask(msg.ID)
	case state.ModifyStep:
	case state.UnlinkTaskStep:
		m.repository.UnlinkTaskStep(msg.Parent.(*domain.Task), msg.Child.(*domain.Step))
	case state.ModifyResource:
	case state.UnlinkTaskResource:
		m.repository.UnlinkTaskResource(msg.Parent.(*domain.Task), msg.Child.(*domain.Resource))
	case state.UnlinkTaskStatus:
		m.repository.UnlinkTaskStatus(msg.Parent.(*domain.Task), msg.Child.(*domain.Status))
	case state.ModifyStatus:
	case state.ModifyCycle:
	case state.ModifyZettel:
		m.repository.DeleteZettel(msg.ID)
	case state.UnlinkZettelResource:
		m.repository.UnlinkZettelResource(msg.Parent.(*domain.Zettel), msg.Child.(*domain.Resource))
	case state.ModifyLink:
		m.repository.UnlinkZettel(msg.Parent.(*domain.Zettel), msg.Child.(*domain.Zettel))
	case state.ModifyAccomplishment:
		m.repository.DeleteAccomplishment(msg.ID)
	case state.ModifySettings:
	case state.ModifyArtifact:
		m.repository.DeleteArtifact(msg.ID)
	case state.ModifyRelease:
		m.repository.DeleteArtifactRelease(msg.Parent.(*domain.Artifact), msg.Child.(*domain.Release))
		m.repository.DeleteRelease(msg.ID)
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
	case state.ModifyArtifact:
		m.repository.UndoDeleteArtifact(msg.ID)
	case state.ModifyRelease:
		// TODO - fill in
	}
}
