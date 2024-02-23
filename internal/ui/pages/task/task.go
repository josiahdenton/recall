package task

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	render "github.com/josiahdenton/recall/internal/ui/renders"
	"github.com/josiahdenton/recall/internal/ui/services/clipboard"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	zettels = iota
	steps
	resources
	header
)

func New() *Model {
	boxStyle := styles.Box(styles.BoxOptions{
		Size:        styles.Full,
		BorderColor: styles.SecondaryGray,
	})

	return &Model{
		box: boxStyle,
	}
}

type Model struct {
	task  *domain.Task
	lists []list.Model
	//steps     list.Model
	//resources list.Model
	//zettels   list.Model
	box lipgloss.Style

	active int
	ready  bool
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var b strings.Builder

	if !m.ready {
		b.WriteString(m.box.Render(""))
	} else {
		// if header active - switch "active/inactive" box styles
		b.WriteString(m.box.Render(render.RenderTaskHeader(m.task)))
		b.WriteString("\n")
		for i, l := range m.lists {
			if m.active != i {
				b.WriteString(styles.InactiveStyle.Render(l.View()))
			} else {
				b.WriteString(l.View())
			}
		}
	}
	return b.String()
}

func (m *Model) Reset() {
	m.ready = false
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global events
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	cmd = m.onGlobalEvents(msg)
	cmds = append(cmds, cmd)

	if m.ready {
		cmd = m.onInput(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) onGlobalEvents(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case router.OnInitPageMsg:
		if msg.Page == router.TasksPage {
			return state.Load(state.Request{
				Type: state.Task,
				ID:   msg.ID,
			})
		}
	case state.LoadedStateMsg:
		// setup lists and
		task, ok := msg.State.(*domain.Task)
		if ok {
			m.setZettels(task.Zettels)
			m.setSteps(task.Steps)
			m.setResources(task.Resources)
			m.ready = true
		}
	}
	return nil
}

func (m *Model) onInput(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			cmd = m.openAddForm()
			cmds = append(cmds, cmd)
		case "e":
			cmd = m.openEditForm()
			cmds = append(cmds, cmd)
		case "enter":
			cmd = m.interact()
			cmds = append(cmds, cmd)
		case "space":
			cmd = m.copy()
			cmds = append(cmds, cmd)
		}
	}

	if m.active < header {
		m.lists[m.active%len(m.lists)], cmd = m.lists[m.active%len(m.lists)].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *Model) openAddForm() tea.Cmd {
	switch m.active {
	case zettels:
		// create new / add from existing
		return toast.ShowToast("unsupported!", toast.Warn)
	case steps:
		// only new
		return router.GotoForm(router.Route{
			Page: router.StepForm,
		})
	case resources:
		// create new / add from existing
		return toast.ShowToast("unsupported!", toast.Warn)
	case header:
		return router.GotoForm(router.Route{
			Page: router.TaskPage,
			ID:   m.task.ID,
		})
	}

	return toast.ShowToast("failed to open add form", toast.Warn)
}

func (m *Model) openEditForm() tea.Cmd {
	switch m.active {
	case zettels:
		return toast.ShowToast("unsupported!", toast.Warn)
	case steps:
		if selected, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			return router.GotoForm(router.Route{
				Page: router.StepForm,
				ID:   selected.ID,
			})
		}
	case resources:
		return toast.ShowToast("unsupported!", toast.Warn)
	case header:
		return router.GotoForm(router.Route{
			Page: router.TaskPage,
			ID:   m.task.ID,
		})
	}

	return toast.ShowToast("failed to open edit form", toast.Warn)
}

func (m *Model) interact() tea.Cmd {
	switch m.active {
	case zettels:
		if selected, ok := m.lists[zettels].SelectedItem().(*domain.Zettel); ok {
			return router.GotoPage(router.Route{
				Page: router.ZettelPage,
				ID:   selected.ID,
			})
		}
	case steps:
		// this should only modify the step itself, not the task
		if selected, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			selected.ToggleStatus()
			return state.Save(state.Request{
				State: *selected,
				Type:  state.Task,
			})
		}
	case resources:
		// open the resource
		if selected, ok := m.lists[resources].SelectedItem().(*domain.Resource); ok && selected.Open() {
			return toast.ShowToast("opening in browser", toast.Info)
		} else {
			return toast.ShowToast("failed to open", toast.Warn)
		}
	case header:
		m.task.ToggleActive()
		return state.Save(state.Request{
			State: *m.task,
			Type:  state.Task,
		})
	}
	return nil
}

func (m *Model) copy() tea.Cmd {
	switch m.active {
	case zettels:
		return toast.ShowToast("unsupported!", toast.Warn)
	case steps:
		if step, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			return clipboard.Copy(step.Description)
		}
	case resources:
		if resource, ok := m.lists[resources].SelectedItem().(*domain.Resource); ok {
			return clipboard.Copy(resource.Source)
		}
	}
	return nil
}

func (m *Model) setZettels(zls []domain.Zettel) {
	m.lists[zettels] = list.New(render.ZettelsToListItems(zls), render.ZettelDelegate{}, 120, 10)
	m.lists[zettels].Title = "Zettels"
	m.lists[zettels].Styles.PaginationStyle = styles.PaginationStyle
	m.lists[zettels].Styles.Title = styles.PageTitleStyle
	m.lists[zettels].SetShowHelp(false)
	m.lists[zettels].KeyMap.Quit.Unbind()
}

func (m *Model) setSteps(sps []domain.Step) {
	m.lists[steps] = list.New(render.StepsToListItems(sps), render.StepDelegate{}, 120, 10)
	m.lists[steps].Title = "Steps"
	m.lists[steps].Styles.PaginationStyle = styles.PaginationStyle
	m.lists[steps].Styles.Title = styles.PageTitleStyle
	m.lists[steps].SetShowHelp(false)
	m.lists[steps].KeyMap.Quit.Unbind()
}

func (m *Model) setResources(rss []domain.Resource) {
	m.lists[resources] = list.New(render.ResourcesToListItems(rss), render.ResourceDelegate{}, 120, 10)
	m.lists[resources].Title = "Resources"
	m.lists[resources].Styles.PaginationStyle = styles.PaginationStyle
	m.lists[resources].Styles.Title = styles.PageTitleStyle
	m.lists[resources].SetShowHelp(false)
	m.lists[resources].KeyMap.Quit.Unbind()
}
