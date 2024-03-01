package task

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	render "github.com/josiahdenton/recall/internal/ui/renders"
	"github.com/josiahdenton/recall/internal/ui/services/clipboard"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"log"
	"strings"
)

const (
	steps = iota
	resources
	status
	header
)

func New() *Model {
	boxStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryGray,
		BoxSize: styles.BoxSize{
			Width:  150,
			Height: 35,
		},
	})

	return &Model{
		box:   boxStyle,
		lists: make([]list.Model, 3),
	}
}

type Model struct {
	task     *domain.Task
	lists    []list.Model
	forms    []tea.Model
	box      lipgloss.Style
	active   int
	ready    bool
	showForm bool
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
		b.WriteString(render.TaskHeader(m.task))
		b.WriteString("\n")
		b.WriteString("\n")
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, m.lists[steps].View(), m.lists[resources].View()))
		b.WriteString("\n")
		b.WriteString(m.lists[status].View())
		b.WriteString("\n")
	}
	return m.box.Render(b.String())
}

func (m *Model) Reset() {
	m.active = 0
	m.ready = false
	m.showForm = false
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case state.DeletedStateMsg:
		// reload the task to pull the latest data
		return state.Load(state.Request{
			Type: state.Task,
			ID:   m.task.ID,
		})
	case router.OnInitPageMsg:
		log.Printf("init task page")
		if msg.Page == router.TaskPage {
			return state.Load(state.Request{
				Type: state.Task,
				ID:   msg.ID,
			})
		}
	case state.LoadedStateMsg:
		log.Printf("loaded %+v", msg)
		if task, ok := msg.State.(*domain.Task); ok {
			m.setSteps(task.Steps)
			m.setResources(task.Resources)
			m.setStatus(task.Status)
			m.task = task
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
		case "d":
			cmd = m.delete()
			cmds = append(cmds, cmd)
		case "enter":
			cmd = m.interact()
			cmds = append(cmds, cmd)
		case "space":
			cmd = m.copy()
			cmds = append(cmds, cmd)
		case "esc":
			cmds = append(cmds, router.Back())
		case "tab":
			m.changeFocus()
		}
	}

	if m.active < header {
		m.lists[m.active], cmd = m.lists[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *Model) changeFocus() {
	if m.active == header {
		m.active = steps
		m.lists[m.active].Styles.Title = styles.ActivePageTitleStyle
	} else if m.active == steps {
		m.lists[m.active].Styles.Title = styles.PageTitleStyle
		m.active = resources
		m.lists[m.active].Styles.Title = styles.ActivePageTitleStyle
	} else if m.active == resources {
		m.lists[m.active].Styles.Title = styles.PageTitleStyle
		m.active = status
		m.lists[m.active].Styles.Title = styles.ActivePageTitleStyle
	} else if m.active == status {
		m.lists[m.active].Styles.Title = styles.PageTitleStyle
		m.active = header
	}
}

func (m *Model) openAddForm() tea.Cmd {
	if m.active == header {
		toast.ShowToast("unsupported!", toast.Warn)
	} else {
		m.showForm = true
	}
	return nil
}

func (m *Model) openEditForm() tea.Cmd {
	switch m.active {
	case status:
		// TODO - add support
		return toast.ShowToast("unsupported!", toast.Warn)
	case steps:
		if selected, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			m.showForm = true
			return forms.EditStep(selected)
		}
	case resources:
		return toast.ShowToast("unsupported!", toast.Warn)
	case header:
		m.showForm = true
		return forms.EditTask(m.task)
	}

	return toast.ShowToast("failed to open edit form", toast.Warn)
}

func (m *Model) delete() tea.Cmd {
	switch m.active {
	case status:
		if selected, ok := m.lists[status].SelectedItem().(*domain.Status); ok {
			m.lists[status].RemoveItem(m.lists[status].Index())
			return state.Delete(state.Request{
				ID:         selected.ID,
				Type:       state.Status,
				ParentType: state.Task,
				Parent:     m.task,
			})
		}
	case steps:
		if selected, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			m.lists[steps].RemoveItem(m.lists[status].Index())
			return state.Delete(state.Request{
				ID:         selected.ID,
				Type:       state.Step,
				ParentType: state.Task,
				Parent:     m.task,
			})
		}
	case resources:
		if selected, ok := m.lists[resources].SelectedItem().(*domain.Resource); ok {
			m.lists[resources].RemoveItem(m.lists[status].Index())
			return state.Delete(state.Request{
				ID:         selected.ID,
				Type:       state.Resource,
				ParentType: state.Task,
				Parent:     m.task,
			})
		}
	case header:
		return toast.ShowToast("unsupported!", toast.Warn)
	}

	return nil
}

func (m *Model) interact() tea.Cmd {
	switch m.active {
	case status:
		return toast.ShowToast("unsupported", toast.Warn)
	case steps:
		if selected, ok := m.lists[steps].SelectedItem().(*domain.Step); ok {
			selected.ToggleStatus()
			return state.Save(state.Request{
				State: *selected,
				Type:  state.Step,
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
	case header:
		return toast.ShowToast("unsupported!", toast.Warn)
	case status:
		if s, ok := m.lists[status].SelectedItem().(*domain.Step); ok {
			return clipboard.Copy(s.Description)
		}
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

func (m *Model) setStatus(statuses []domain.Status) {
	m.lists[status] = list.New(render.StatusToListItems(statuses), render.StatusDelegate{}, 120, 10)
	m.lists[status].Title = "Status"
	m.lists[status].Styles.PaginationStyle = styles.PaginationStyle
	m.lists[status].Styles.Title = styles.PageTitleStyle
	m.lists[status].SetShowHelp(false)
	m.lists[status].KeyMap.Quit.Unbind()
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
