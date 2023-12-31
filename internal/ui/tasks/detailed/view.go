package detailed

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	styles2 "github.com/josiahdenton/recall/internal/ui/styles"
	forms2 "github.com/josiahdenton/recall/internal/ui/tasks/detailed/forms"
	tasklist "github.com/josiahdenton/recall/internal/ui/tasks/list"
	"log"
	"strings"
	"time"
)

var (
	listTitleStyle     = styles2.SecondaryGray.Copy()
	activeListStyle    = styles2.SecondaryColor.Copy()
	statusMessageStyle = styles2.PrimaryColor.Copy().PaddingLeft(1)
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
}

func New() *Model {
	formList := make([]tea.Model, formCount)
	formList[steps] = forms2.NewStepForm()
	formList[resources] = forms2.NewStepResourceForm()
	formList[status] = forms2.NewStatusForm()
	return &Model{
		headerActive: true,
		active:       header,
		lists:        make([]list.Model, formCount), // TODO remove magic num
		forms:        formList,
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
	return styles2.WindowStyle.Render(b.String())
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
	case tasklist.GotoDetailedPageMsg:
		m.task = msg.Task
		m.setupLists(msg.Task)
		m.ready = true
	case clearStatusMessage:
		m.statusMessage = ""
	case forms2.StepFormMsg:
		m.task.Steps = append(m.task.Steps, msg.Step)
		m.lists[steps].InsertItem(len(m.task.Steps)-1, &m.task.Steps[len(m.task.Steps)-1])
		m.showForm = false
	case forms2.ResourceFormMsg:
		m.task.Resources = append(m.task.Resources, msg.Resource)
		m.lists[resources].InsertItem(len(m.task.Resources)-1, &m.task.Resources[len(m.task.Resources)-1])
		m.showForm = false
	case forms2.StatusFormMsg:
		m.task.Status = append(m.task.Status, msg.Status)
		m.lists[status].InsertItem(len(m.task.Status)-1, &m.task.Status[len(m.task.Status)-1])
		m.showForm = false
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPage(router.TaskListPage))
				m.active = header
			}
		}

		if !m.showForm {
			cmd = m.detailedControls(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	if m.ready && !m.showForm && m.active != header {
		m.lists[m.active], cmd = m.lists[m.active].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) detailedControls(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "tab":
		// forms handle tab
		if m.active < header {
			m.lists[m.active].Styles.Title = listTitleStyle
		}
		m.active = m.nextSection()

		if m.active < header {
			m.lists[m.active].Styles.Title = activeListStyle
		}
	case "a":
		if m.active < header {
			m.showForm = true
		}
	case "d":
		if m.active < header && len(m.lists[m.active].Items()) > 0 {
			switch m.active {
			case steps:
				m.task.RemoveStep(m.lists[m.active].Index())
			case resources:
				m.task.RemoveResource(m.lists[m.active].Index())
			case status:
				m.task.RemoveStatus(m.lists[m.active].Index())
			}
			m.lists[m.active].RemoveItem(m.lists[m.active].Index())
			m.statusMessage = "removed item!"
			cmd = clearStatus()
		}
	case " ":
		switch m.active {
		case steps:
			step := m.lists[steps].SelectedItem().(*domain.Step)
			step.ToggleStatus()
			if step.Complete {
				m.statusMessage = "step completed!"
			} else {
				m.statusMessage = "step reset!"
			}
		case resources:
			// this should open the resource
			// refer to the other golang app for that...
			// for now, copy source to clipboard
			resource := m.lists[resources].SelectedItem().(*domain.Resource)
			err := clipboard.WriteAll(resource.Source)
			if err != nil {
				log.Printf("failed to copy to clipboard %v", err)
			}
			m.statusMessage = "copied to clipboard!"
		case status:
			status := m.lists[status].SelectedItem().(*domain.Status)
			err := clipboard.WriteAll(status.Description)
			if err != nil {
				log.Printf("failed to copy to clipboard %v", err)
			}
			m.statusMessage = "copied to clipboard!"
		}
		cmd = clearStatus()
	}
	return cmd
}

func (m *Model) setupLists(task *domain.Task) {
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

	m.lists[steps] = list.New(_steps, stepDelegate{}, 50, 7)
	m.lists[steps].Title = "Steps"
	m.lists[steps].SetFilteringEnabled(false)
	m.lists[steps].Styles.Title = listTitleStyle
	m.lists[steps].SetShowHelp(false)
	m.lists[steps].SetFilteringEnabled(false)
	m.lists[steps].SetShowStatusBar(false)
	m.lists[steps].KeyMap.Quit.Unbind()

	m.lists[resources] = list.New(_resources, resourceDelegate{}, 50, 7)
	m.lists[resources].Title = "Resources"
	m.lists[resources].SetFilteringEnabled(false)
	m.lists[resources].Styles.Title = listTitleStyle
	m.lists[resources].SetShowHelp(false)
	m.lists[resources].SetFilteringEnabled(false)
	m.lists[resources].SetShowStatusBar(false)
	m.lists[resources].KeyMap.Quit.Unbind()

	m.lists[status] = list.New(_status, statusDelegate{}, 50, 5)
	m.lists[status].Title = "Status"
	m.lists[status].SetFilteringEnabled(false)
	m.lists[status].Styles.Title = listTitleStyle
	m.lists[status].SetShowHelp(false)
	m.lists[status].SetFilteringEnabled(false)
	m.lists[status].SetShowStatusBar(false)
	m.lists[status].KeyMap.Quit.Unbind()
}

func (m *Model) nextSection() int {
	switch m.active {
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

type clearStatusMessage struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*5, func(_ time.Time) tea.Msg {
		return clearStatusMessage{}
	})
}
