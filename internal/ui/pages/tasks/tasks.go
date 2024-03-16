package tasks

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	render "github.com/josiahdenton/recall/internal/ui/renders"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	tasks = iota
	resources
)

func New() *Model {
	tasksStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryColor,
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: styles.BaseHeight*2 - 10,
		},
	})
	inactiveTasksStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryGray,
		TextColor:   styles.SecondaryColor,
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: styles.BaseHeight*2 - 10,
		},
	})
	resourcesStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryColor,
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: 10,
		},
	})
	inactiveResourcesStyle := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryGray,
		TextColor:   styles.SecondaryColor,
		BoxSize: styles.BoxSize{
			Width:  styles.BaseWidth * 2,
			Height: 10,
		},
	})
	return &Model{
		forms:                  []forms.Form{forms.NewTaskForm()},
		tasksStyle:             tasksStyle,
		inactiveTasksStyle:     inactiveTasksStyle,
		resourcesStyle:         resourcesStyle,
		inactiveResourcesStyle: inactiveResourcesStyle,
	}
}

type Model struct {
	tasks                  list.Model
	resources              list.Model
	forms                  []forms.Form
	tasksStyle             lipgloss.Style
	inactiveTasksStyle     lipgloss.Style
	resourcesStyle         lipgloss.Style
	inactiveResourcesStyle lipgloss.Style
	ready                  bool
	resourcesLoaded        bool
	tasksLoaded            bool
	showForm               bool
	active                 int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	if m.ready && m.showForm {
		// note - forms handle their own styling
		b.WriteString(m.forms[m.active].View())
	} else if m.ready {
		if m.active == tasks {
			b.WriteString(m.tasksStyle.Render(m.tasks.View()))
			b.WriteString("\n")
			b.WriteString(m.inactiveResourcesStyle.Render(m.resources.View()))
		} else {
			b.WriteString(m.inactiveTasksStyle.Render(m.tasks.View()))
			b.WriteString("\n")
			b.WriteString(m.resourcesStyle.Render(m.resources.View()))
		}
	}

	return b.String()
}

func (m *Model) Reset() {
	m.ready = false
	m.showForm = false
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// global
	cmd = m.onGlobalEvents(msg)
	cmds = append(cmds, cmd)

	// input
	if m.ready && !m.showForm && m.active == tasks {
		cmd = m.onInputTasks(msg)
		cmds = append(cmds, cmd)
	} else if m.ready && !m.showForm && m.active == resources {
		cmd = m.onInputResources(msg)
		cmds = append(cmds, cmd)
	} else if m.ready && m.showForm {
		cmd = m.onFormInput(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) onGlobalEvents(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case forms.TaskFormMsg:
		if msg.Edit {
			m.showForm = false
			return toast.ShowToast("Modified Task", toast.Info)
		}
		m.tasks.InsertItem(len(m.tasks.Items()), &msg.Task)
		// TODO - should also send a toast
		return state.Save(state.Request{
			State: msg.Task,
			Type:  state.Task,
		})
	case router.OnInitPageMsg:
		if msg.Page == router.TasksPage {
			return tea.Batch(
				state.Load(state.Request{
					Type: state.Tasks,
				}),
				state.Load(state.Request{
					Type: state.Resources,
				}),
			)
		}
	case state.LoadedStateMsg:
		if msg.Type == state.Tasks {
			tasks, ok := msg.State.([]domain.Task)
			if !ok {
				return toast.ShowToast("failed to fetch tasks", toast.Warn)
			}
			m.setTasks(tasks)
			m.tasksLoaded = true
		} else if msg.Type == state.Resources {
			resources, ok := msg.State.([]domain.Resource)
			if !ok {
				return toast.ShowToast("failed to fetch resources", toast.Warn)
			}
			m.setResources(resources)
			m.resourcesLoaded = true
		}
		m.ready = m.resourcesLoaded && m.tasksLoaded
	}

	// no match found
	return nil
}

func (m *Model) onInputResources(msg tea.Msg) tea.Cmd {
	if m.resources.FilterState() != list.Filtering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				m.changeFocus()
			case "enter":
				if selected, ok := m.resources.SelectedItem().(*domain.Resource); ok {
					if selected.Open() {
						return toast.ShowToast("opening web page!", toast.Info)
					} else {
						return toast.ShowToast("failed to open resource", toast.Warn)
					}
				}
			case "e":
				if _, ok := m.resources.SelectedItem().(*domain.Task); ok {
					m.showForm = true
					// TODO forms.EditResource
				}
			case "a":
				// cmd to go to forms "add task"
				m.showForm = true
			case "esc":
				return router.Back()
			}
		}
	}

	var cmd tea.Cmd
	m.resources, cmd = m.resources.Update(msg)
	return cmd
}

func (m *Model) onInputTasks(msg tea.Msg) tea.Cmd {
	if m.tasks.FilterState() != list.Filtering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				m.changeFocus()
			case "enter":
				if selected, ok := m.tasks.SelectedItem().(*domain.Task); ok {
					return router.GotoPage(router.Route{
						Page: router.TaskPage,
						ID:   selected.ID,
					})
				}
			case "e":
				if selected, ok := m.tasks.SelectedItem().(*domain.Task); ok {
					m.showForm = true
					return forms.EditTask(selected)
				}
			case "a":
				// cmd to go to forms "add task"
				m.showForm = true
			case "esc":
				return router.Back()
			}
		}
	}

	var cmd tea.Cmd
	m.tasks, cmd = m.tasks.Update(msg)
	return cmd
}

func (m *Model) onFormInput(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.showForm = false
			m.forms[m.active].Reset()
			return nil
		}
	}
	var cmd tea.Cmd
	m.forms[m.active], cmd = m.forms[m.active].Update(msg)
	return cmd
}

func (m *Model) changeFocus() {
	if m.active == tasks {
		m.active = resources
	} else {
		m.active = tasks
	}
}

func (m *Model) setTasks(tasks []domain.Task) {
	m.tasks = list.New(render.TasksToListItems(tasks), render.TaskDelegate{}, 50, 20)
	m.tasks.Title = "Tasks"
	m.tasks.Styles.PaginationStyle = styles.PaginationStyle
	m.tasks.Styles.Title = styles.PageTitleStyle
	m.tasks.SetShowHelp(false)
	m.tasks.KeyMap.Quit.Unbind()
}

func (m *Model) setResources(resources []domain.Resource) {
	m.resources = list.New(render.ResourcesToListItems(resources), render.ResourceDelegate{}, 50, 10)
	m.resources.Title = "Resources"
	m.resources.Styles.PaginationStyle = styles.PaginationStyle
	m.resources.Styles.Title = styles.PageTitleStyle
	m.resources.SetShowHelp(false)
	m.resources.KeyMap.Quit.Unbind()
}
