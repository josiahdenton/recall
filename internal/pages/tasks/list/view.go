package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/pages/tasks"
	"io"
	"time"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

type redrawTick time.Time

type Model struct {
	ready bool
	tasks list.Model
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
	task, ok := item.(*tasks.Task)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderTask(task, index == m.Index()))
}

func loadTasks() tea.Msg {
	exampleDuration, _ := time.ParseDuration("2hr")
	taskList := make([]list.Item, 7)
	tl := []*tasks.Task{
		{Title: "update PM for EOY", Priority: tasks.Low, Due: "11/22/23", LastActivatedTime: time.Now()},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23", LastActivatedTime: time.Now()},
		{Title: "pickup toys", LastActivatedTime: time.Now()},
		{Title: "update docs", Priority: tasks.None, Active: true, TotalActiveTime: exampleDuration, LastActivatedTime: time.Now()},
		{Title: "message Sesha", Priority: tasks.High, LastActivatedTime: time.Now()},
		{Title: "update PM for EOY", Priority: tasks.Low, Due: "11/22/23", LastActivatedTime: time.Now()},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23", LastActivatedTime: time.Now()},
	}
	for i, item := range tl {
		taskList[i] = item
	}
	return LoadTasks{Tasks: taskList}
}

func (m Model) Init() tea.Cmd {
	return loadTasks
}

func (m Model) View() string {
	return "\n" + m.tasks.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case LoadTasks:
		m.tasks = list.New(msg.Tasks, taskDelegate{}, 50, 20)
		m.tasks.SetShowStatusBar(false)
		m.tasks.SetFilteringEnabled(false)
		m.tasks.Title = "Tasks"
		m.tasks.Styles.PaginationStyle = paginationStyle
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
	}
	return m, cmd
}
