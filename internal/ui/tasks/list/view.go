package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

type Model struct {
	ready bool
	tasks list.Model
}

type GotoDetailedPageMsg struct {
	Task *domain.Task
}

type LoadTasks struct {
	Tasks []list.Item
}

type taskDelegate struct{}

func (d taskDelegate) Height() int  { return 1 }
func (d taskDelegate) Spacing() int { return 0 }
func (d taskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	task, ok := item.(*domain.Task)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderTask(task, index == m.Index()))
}

// I feel like this should move to core...
func loadTasks() tea.Msg {
	tl := []domain.Task{
		{Title: "update PM for EOY", Priority: domain.TaskPriorityLow, Due: "11/22/23",
			Steps: []domain.Step{
				{
					Description: "tnsteisieraTIRSTREIAns",
				},
				{
					Description: "sintnrneirn",
				},
				{
					Description: "snitrnseeisintnrneirrnitern",
				},
			},
			Status: []domain.Status{
				{
					Description: "Needed to meet with my manager to discuss where this is headed hNeeded to meet with my manager to discuss where this is headed h",
				},
				{
					Description: "After meeting, moved toward creating a more sensible position for this",
				},
				{
					Description: "Needed to meet with my manager to discuss where this is headed",
				},
				{
					Description: "After meeting, moved toward creating a more sensible position for this",
				},
				{
					Description: "Needed to meet with my manager to discuss where this is headed",
				},
				{
					Description: "After meeting, moved toward creating a more sensible position for this",
				},
			},
		},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "pickup toys"},
		{Title: "update docs", Priority: domain.TaskPriorityNone, Active: true},
		{Title: "message Sesha", Priority: domain.TaskPriorityHigh},
		{Title: "update PM for EOY", Priority: domain.TaskPriorityLow, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "update PM for EOY", Priority: domain.TaskPriorityLow, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "pickup toys"},
		{Title: "update docs", Priority: domain.TaskPriorityNone, Active: true},
		{Title: "message Sesha", Priority: domain.TaskPriorityHigh},
		{Title: "update PM for EOY", Priority: domain.TaskPriorityLow, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
	}
	taskList := make([]list.Item, len(tl))
	for i := range tl {
		task := &tl[i]
		taskList[i] = task
	}
	return LoadTasks{Tasks: taskList}
}

func (m *Model) Init() tea.Cmd {
	return loadTasks
}

func (m *Model) View() string {
	return styles.WindowStyle.Render(m.tasks.View())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case LoadTasks:
		// TODO split this setup into its own func
		m.tasks = list.New(msg.Tasks, taskDelegate{}, 50, 20)
		m.tasks.SetShowStatusBar(false)
		m.tasks.SetFilteringEnabled(false)
		m.tasks.Title = "Tasks"
		m.tasks.Styles.PaginationStyle = paginationStyle
		m.tasks.Styles.Title = titleStyle
		m.tasks.SetShowHelp(false)
		m.tasks.KeyMap.Quit.Unbind()
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item := m.tasks.SelectedItem().(*domain.Task)
			cmds = append(cmds, ShowDetailedView(item))
		}
	}

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// ShowDetailedView TODO refactor this to the router
func ShowDetailedView(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		return GotoDetailedPageMsg{
			Task: task,
		}
	}
}
