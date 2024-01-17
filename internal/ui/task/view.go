package task

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"log"
	"strings"
	"time"
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
)

const (
	formCount = 3
)

type Model struct {
	ready         bool
	showForm      bool
	headerActive  bool
	forms         []tea.Model
	statusMessage string
	task          *domain.Task
	lists         []list.Model
	active        int
	commands      Commands
}

func New() *Model {
	formList := make([]tea.Model, formCount)
	formList[steps] = forms.NewStepForm()
	formList[resources] = forms.NewResourceForm()
	formList[status] = forms.NewStatusForm()
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
		b.WriteString(statusMessageStyle.Render(m.statusMessage))
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

	// TODO - have the task form do something like this for due date
	//func mustParseDate(date string) time.Time {
	//	input := fmt.Sprintf("%s at 7:00am (EST)", date)
	//	t, err := time.Parse(longDateForm, input)
	//	if err != nil {
	//	return time.Time{}
	//}
	//	return t
	//}

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		m.task = msg.State.(*domain.Task)
		m.lists = setupLists(m.task)
		m.ready = true
	case clearStatusMessage:
		m.statusMessage = ""
	case forms.StepFormMsg:
		m.task.Steps = append(m.task.Steps, msg.Step)
		m.lists[steps].InsertItem(len(m.task.Steps), &m.task.Steps[len(m.task.Steps)-1])
		m.showForm = false
		cmds = append(cmds, updateTask(m.task))
	case forms.ResourceFormMsg:
		m.task.Resources = append(m.task.Resources, msg.Resource)
		m.lists[resources].InsertItem(len(m.task.Resources), &m.task.Resources[len(m.task.Resources)-1])
		m.showForm = false
		cmds = append(cmds, updateTask(m.task))
	case forms.StatusFormMsg:
		m.task.Status = append(m.task.Status, msg.Status)
		m.lists[status].InsertItem(len(m.task.Status), &m.task.Status[len(m.task.Status)-1])
		m.showForm = false
		cmds = append(cmds, updateTask(m.task))
	case tea.KeyMsg:
		action := m.commands.HandleInput(msg)
		switch action {
		case Back:
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPage(domain.TaskListPage, 0))
				// TODO - add a Reset method
				m.active = header
			}
		case Interact:
			if m.showForm {
				break
			}

			switch m.active {
			case steps:
				step := m.lists[steps].SelectedItem().(*domain.Step)
				step.ToggleStatus()
				if step.Complete {
					m.statusMessage = "completed step!"
				} else {
					m.statusMessage = "reset step!"
				}
				cmds = append(cmds, clearStatus(), updateStep(step))
			case resources:
				resource := m.lists[resources].SelectedItem().(*domain.Resource)
				if resource.Type == domain.WebResource {
					resource.Open()
					m.statusMessage = "opened web link!"
				} else {
					m.statusMessage = "unsupported type!"
				}
				cmds = append(cmds, clearStatus())
			case status:
				status := m.lists[status].SelectedItem().(*domain.Status)
				err := clipboard.WriteAll(status.Description)
				if err != nil {
					log.Printf("failed to copy to clipboard: %v", err)
					m.statusMessage = "failed to copy to clipboard"
				}
				m.statusMessage = "copied to clipboard!"
				cmds = append(cmds, clearStatus())
			case header:
				// nothing for now
			}
		case Delete:
			if m.showForm {
				break
			}

			if m.active < header && len(m.lists[m.active].Items()) > 0 {
				index := m.lists[m.active].Index()
				item := m.lists[m.active].SelectedItem()
				switch m.active {
				case steps:
					m.task.Steps = append(m.task.Steps[:index], m.task.Steps[index+1:]...)
					cmds = append(cmds, clearStatus(), deleteStep(m.task, item.(*domain.Step)))
				case resources:
					m.task.Resources = append(m.task.Resources[:index], m.task.Resources[index+1:]...)
					cmds = append(cmds, clearStatus(), deleteResource(m.task, item.(*domain.Resource)))
				case status:
					m.task.Status = append(m.task.Status[:index], m.task.Status[index+1:]...)
					cmds = append(cmds, clearStatus(), deleteStatus(m.task, item.(*domain.Status)))
				}
				m.lists[m.active].RemoveItem(index)
				m.statusMessage = "removed item!"
			}
		case Add:
			if !m.showForm && m.active < header {
				m.showForm = true
			}
		case MoveFocus:
			if m.showForm {
				break
			}

			if m.active < header {
				m.lists[m.active].Styles.Title = listTitleStyle
			}
			m.active = nextSection(m.active)
			if m.active < header {
				m.lists[m.active].Styles.Title = activeListStyle
			}
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
		return shared.SaveStateMsg{
			Update: *task,
			Type:   shared.ModifyTask,
		}
	}
}

func updateStep(step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Update: *step,
			Type:   shared.ModifyStep,
		}
	}
}

func deleteStep(task *domain.Task, step *domain.Step) tea.Cmd {
	return func() tea.Msg {
		return shared.DeleteStateMsg{
			Type:   shared.ModifyStep,
			Parent: task,
			Child:  step,
		}
	}
}

func deleteResource(task *domain.Task, resource *domain.Resource) tea.Cmd {
	return func() tea.Msg {
		return shared.DeleteStateMsg{
			Type:   shared.ModifyResource,
			Parent: task,
			Child:  resource,
		}
	}
}

func deleteStatus(task *domain.Task, status *domain.Status) tea.Cmd {
	return func() tea.Msg {
		return shared.DeleteStateMsg{
			Type:   shared.ModifyStatus,
			Parent: task,
			Child:  status,
		}
	}
}

type clearStatusMessage struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*5, func(_ time.Time) tea.Msg {
		return clearStatusMessage{}
	})
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
	lists[steps].SetFilteringEnabled(false)
	lists[steps].Styles.Title = listTitleStyle
	lists[steps].SetShowHelp(false)
	lists[steps].SetFilteringEnabled(false)
	lists[steps].SetShowStatusBar(false)
	lists[steps].KeyMap.Quit.Unbind()

	lists[resources] = list.New(_resources, resourceDelegate{}, 80, 7)
	lists[resources].Title = "Resources"
	lists[resources].SetFilteringEnabled(false)
	lists[resources].Styles.Title = listTitleStyle
	lists[resources].SetShowHelp(false)
	lists[resources].SetFilteringEnabled(false)
	lists[resources].SetShowStatusBar(false)
	lists[resources].KeyMap.Quit.Unbind()

	lists[status] = list.New(_status, statusDelegate{}, 80, 5)
	lists[status].Title = "Status"
	lists[status].SetFilteringEnabled(false)
	lists[status].Styles.Title = listTitleStyle
	lists[status].SetShowHelp(false)
	lists[status].SetFilteringEnabled(false)
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
