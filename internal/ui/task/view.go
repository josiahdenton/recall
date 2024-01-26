package task

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"log"
	"strings"
)

var (
	listTitleStyle     = styles.SecondaryGray.Copy()
	activeListStyle    = styles.SecondaryColor.Copy()
	statusMessageStyle = styles.PrimaryColor.Copy().PaddingLeft(1)
)

// active options
const (
	steps = iota
	resources
	status
	header
	formCount
)

type Model struct {
	ready        bool
	showForm     bool
	headerActive bool
	forms        []tea.Model
	task         *domain.Task
	lists        []list.Model
	active       int
	commands     Commands
}

func New() *Model {
	formList := make([]tea.Model, formCount)
	formList[steps] = forms.NewStepForm()
	formList[resources] = forms.NewResourceForm()
	formList[status] = forms.NewStatusForm()
	formList[header] = forms.NewTaskForm()
	return &Model{
		headerActive: true,
		active:       header,
		forms:        formList,
		commands:     DefaultCommands(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	if m.ready && !m.showForm {
		b.WriteString(renderHeader(m.task, m.active == header))
		b.WriteString(m.lists[steps].View() + "\n")
		b.WriteString(m.lists[resources].View() + "\n")
		b.WriteString(m.lists[status].View() + "\n")
	} else if m.showForm {
		b.WriteString(renderHeader(m.task, m.active == header))
		b.WriteString(m.forms[m.active].View())
	} else {
		b.WriteString("loading...")
	}
	return styles.WindowStyle.Render(b.String())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO switch to use the Keymap instead of hard coding everything!
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.showForm {
		m.forms[m.active], cmd = m.forms[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		m.task = msg.State.(*domain.Task)
		m.lists = setupLists(m.task)
		m.ready = true
	case forms.StepFormMsg:
		m.task.Steps = append(m.task.Steps, msg.Step)
		m.lists[steps].InsertItem(len(m.task.Steps), &m.task.Steps[len(m.task.Steps)-1])
		cmds = append(cmds, updateTask(m.task))
	case forms.ResourceFormMsg:
		m.task.Resources = append(m.task.Resources, msg.Resource)
		m.lists[resources].InsertItem(len(m.task.Resources), &m.task.Resources[len(m.task.Resources)-1])
		cmds = append(cmds, updateTask(m.task))
	case forms.StatusFormMsg:
		m.task.Status = append(m.task.Status, msg.Status)
		m.lists[status].InsertItem(len(m.task.Status), &m.task.Status[len(m.task.Status)-1])
		cmds = append(cmds, updateTask(m.task))
	case state.SaveStateMsg:
		if msg.Type == state.ModifyTask {
			m.showForm = false
		}
	case tea.KeyMsg:
		action := m.commands.HandleInput(msg)
		switch action {
		case Back:
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPreviousPage())
				// TODO - add a Reset method
				m.active = header
			}
		case Interact:
			if m.showForm || (m.active < header && len(m.lists[m.active].Items()) < 1) {
				break
			}

			switch m.active {
			case steps:
				step := m.lists[steps].SelectedItem().(*domain.Step)
				step.ToggleStatus()
				if step.Complete {
					cmds = append(cmds, updateStep(step), toast.ShowToast("completed step!"))
				} else {
					cmds = append(cmds, updateStep(step), toast.ShowToast("reset step!"))
				}
			case resources:
				resource := m.lists[resources].SelectedItem().(*domain.Resource)
				if resource.Type == domain.WebResource {
					resource.Open()
					cmds = append(cmds, toast.ShowToast("opened web link!"))
				} else {
					cmds = append(cmds, toast.ShowToast("unsupported type!"))
				}
			case status:
				status := m.lists[status].SelectedItem().(*domain.Status)
				err := clipboard.WriteAll(status.Description)
				if err != nil {
					log.Printf("failed to copy to clipboard: %v", err)
					cmds = append(cmds, toast.ShowToast("failed to copy to clipboard"))
				}
				cmds = append(cmds, toast.ShowToast("copied to clipboard!"))
			case header:
				// TODO - this should open the task form prefilled
				if !m.showForm {
					m.showForm = true
					cmds = append(cmds, forms.EditTask(m.task))
				}
			}
		case Delete:
			if m.showForm || (m.active < header && m.lists[m.active].FilterState() == list.Filtering) {
				break
			}

			if m.active < header && len(m.lists[m.active].Items()) > 0 {
				index := m.lists[m.active].Index()
				switch m.active {
				case steps:
					item := m.task.Steps[index]
					cmds = append(cmds, deleteStep(m.task, &item))
					m.task.Steps = append(m.task.Steps[:index], m.task.Steps[index+1:]...)
				case resources:
					item := m.task.Resources[index]
					cmds = append(cmds, deleteResource(m.task, &item))
					m.task.Resources = append(m.task.Resources[:index], m.task.Resources[index+1:]...)
				case status:
					item := m.task.Status[index]
					cmds = append(cmds, deleteStatus(m.task, &item))
					m.task.Status = append(m.task.Status[:index], m.task.Status[index+1:]...)
				}
				m.lists[m.active].RemoveItem(index)
				cmds = append(cmds, toast.ShowToast("removed item!"))
			}
		case Add:
			if !m.showForm && m.active < header && m.lists[m.active].FilterState() != list.Filtering {
				m.showForm = true
			}
		case MoveFocus:
			if m.showForm || (m.active < header && m.lists[m.active].FilterState() == list.Filtering) {
				break
			}

			if m.active < header {
				m.lists[m.active].Styles.Title = listTitleStyle
			}
			m.active = nextSection(m.active)
			if m.active < header {
				m.lists[m.active].Styles.Title = activeListStyle
			}
		case Favorite:
			if m.active == header {
				m.task.ToggleFavorite()
				cmds = append(cmds, updateTask(m.task), toast.ShowToast("favorited task!"))
			}
		case None:
		}
	}

	if m.ready && !m.showForm && m.active != header {
		m.lists[m.active], cmd = m.lists[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func updateTask(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *task,
			Type:   state.ModifyTask,
		}
	}
}

func updateStep(step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *step,
			Type:   state.ModifyStep,
		}
	}
}

func deleteStep(task *domain.Task, step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.ModifyStep,
			Parent: task,
			Child:  step,
		}
	}
}

func deleteResource(task *domain.Task, resource *domain.Resource) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.ModifyResource,
			Parent: task,
			Child:  resource,
		}
	}
}

func deleteStatus(task *domain.Task, status *domain.Status) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.ModifyStatus,
			Parent: task,
			Child:  status,
		}
	}
}

func setupLists(task *domain.Task) []list.Model {
	_steps := make([]list.Item, len(task.Steps))
	_resources := make([]list.Item, len(task.Resources))
	_status := make([]list.Item, len(task.Status))
	for i := range task.Steps {
		s := &task.Steps[i]
		_steps[i] = s
	}
	for i := range task.Resources {
		r := &task.Resources[i]
		_resources[i] = r
	}
	for i := range task.Status {
		s := &task.Status[i]
		_status[i] = s
	}
	lists := make([]list.Model, formCount)

	lists[steps] = list.New(_steps, stepDelegate{}, 80, 9)
	lists[steps].Title = "Steps"
	lists[steps].Styles.Title = listTitleStyle
	lists[steps].SetFilteringEnabled(false)
	lists[steps].SetShowHelp(false)
	lists[steps].SetShowStatusBar(false)
	lists[steps].KeyMap.Quit.Unbind()

	lists[resources] = list.New(_resources, resourceDelegate{}, 80, 9)
	lists[resources].Title = "Resources"
	lists[resources].Styles.Title = listTitleStyle
	lists[resources].SetShowHelp(false)
	lists[resources].SetShowStatusBar(false)
	lists[resources].KeyMap.Quit.Unbind()

	lists[status] = list.New(_status, statusDelegate{}, 80, 5)
	lists[status].Title = "Status"
	lists[status].Styles.Title = listTitleStyle
	lists[status].SetFilteringEnabled(false)
	lists[status].SetShowHelp(false)
	lists[status].SetShowStatusBar(false)
	lists[status].KeyMap.Quit.Unbind()

	return lists
}

// TODO - can make this a func
func nextSection(active int) int {
	switch active {
	case header:
		return steps
	case steps:
		return resources
	case resources:
		return status
	case status:
		return header
	}
	return header
}
