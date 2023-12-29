package detailed

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/pages/styles"
	"github.com/josiahdenton/recall/internal/pages/tasks"
	tasklist "github.com/josiahdenton/recall/internal/pages/tasks/list"
	"log"
	"strings"
	"time"
)

var (
	titleStyle         = styles.PrimaryColor.Copy()
	listTitleStyle     = styles.SecondaryGray.Copy()
	activeListStyle    = styles.SecondaryColor.Copy()
	statusMessageStyle = styles.PrimaryColor.Copy().PaddingLeft(1)
)

const (
	steps = iota
	resources
	status
)

type Model struct {
	ready         bool
	statusMessage string
	task          *tasks.Task
	lists         []list.Model
	//steps      list.Model
	//resources  list.Model
	//status     list.Model
	activeList int
}

func New() *Model {
	return &Model{
		activeList: -1,
		lists:      make([]list.Model, 3),
	}
}

func (m *Model) Init() tea.Cmd {
	// TODO - this may have to load more?
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	if m.ready {
		b.WriteString(titleStyle.Render(m.task.Title) + "\n")
		b.WriteString(fmt.Sprintf("%s  %s\n\n", m.task.Due, m.task.Runtime()))
		b.WriteString(m.lists[steps].View() + "\n")
		b.WriteString(m.lists[resources].View() + "\n")
		b.WriteString(m.lists[status].View() + "\n")
		b.WriteString(statusMessageStyle.Render(m.statusMessage))
	} else {
		b.WriteString("loading...")
	}
	return styles.WindowStyle.Render(b.String())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tasklist.ShowDetailedMsg:
		m.task = msg.Task
		m.setupLists(msg.Task)
		m.ready = true
	case clearStatusMessage:
		m.statusMessage = ""
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.activeList > -1 {
				m.lists[m.activeList%len(m.lists)].Styles.Title = listTitleStyle
			}
			m.activeList++
			m.lists[m.activeList%len(m.lists)].Styles.Title = activeListStyle
		case " ":
			// TODO every time an action occurs, I should flash a msg somewhere...
			// I would need "add a msg" to my model and then also add a cmd to the Batch
			// in order to clear that message
			switch m.activeList % len(m.lists) {
			case steps:
				step := m.lists[steps].SelectedItem().(*tasks.Step)
				step.Complete = !step.Complete
				m.statusMessage = "step completed!"
			case resources:
				// this should open the resource
				// refer to the other golang app for that...
				// for now, also copy it to clipboard
				resource := m.lists[resources].SelectedItem().(*tasks.Resource)
				err := clipboard.WriteAll(resource.Source)
				if err != nil {
					log.Printf("failed to copy to clipboard %v", err)
				}
				m.statusMessage = "copied to clipboard!"
			case status:
				status := m.lists[status].SelectedItem().(*tasks.Status)
				err := clipboard.WriteAll(status.Description)
				if err != nil {
					log.Printf("failed to copy to clipboard %v", err)
				}
				m.statusMessage = "copied to clipboard!"
			}
			cmds = append(cmds, clearStatus())
		}
	}
	if m.ready {
		var cmd tea.Cmd
		switch m.activeList % len(m.lists) {
		case steps:
			m.lists[steps], cmd = m.lists[steps].Update(msg)
			cmds = append(cmds, cmd)
		case resources:
			m.lists[resources], cmd = m.lists[resources].Update(msg)
			cmds = append(cmds, cmd)
		case status:
			m.lists[status], cmd = m.lists[status].Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) setupLists(task *tasks.Task) {
	_steps := make([]list.Item, len(task.Steps))
	_resources := make([]list.Item, len(task.Resources))
	_status := make([]list.Item, len(task.Status))
	for i := range task.Steps {
		s := task.Steps[i]
		_steps[i] = &s
	}
	for i := range task.Resources {
		r := task.Resources[i]
		_resources[i] = &r
	}
	for i := range task.Status {
		s := task.Status[i]
		_status[i] = &s
	}

	m.lists[steps] = list.New(_steps, stepDelegate{}, 50, 8)
	m.lists[steps].Title = "Steps"
	m.lists[steps].SetFilteringEnabled(false)
	m.lists[steps].Styles.Title = listTitleStyle
	m.lists[steps].SetShowHelp(false)
	m.lists[steps].SetFilteringEnabled(false)
	m.lists[steps].SetShowStatusBar(false)

	m.lists[resources] = list.New(_resources, resourceDelegate{}, 50, 8)
	m.lists[resources].Title = "Resources"
	m.lists[resources].SetFilteringEnabled(false)
	m.lists[resources].Styles.Title = listTitleStyle
	m.lists[resources].SetShowHelp(false)
	m.lists[resources].SetFilteringEnabled(false)
	m.lists[resources].SetShowStatusBar(false)

	m.lists[status] = list.New(_status, statusDelegate{}, 50, 8)
	m.lists[status].Title = "Status"
	m.lists[status].SetFilteringEnabled(false)
	m.lists[status].Styles.Title = listTitleStyle
	m.lists[status].SetShowHelp(false)
	m.lists[status].SetFilteringEnabled(false)
	m.lists[status].SetShowStatusBar(false)
}

type clearStatusMessage struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*5, func(_ time.Time) tea.Msg {
		return clearStatusMessage{}
	})
}
