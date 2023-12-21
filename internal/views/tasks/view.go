package tasks

import (
	"fmt"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

// ------------------------------- //
//             styles              //
// ------------------------------- //

func New() Model {
	return Model{
		active: -1,
	}
}

// ------------------ //
//    tea messages    //
// ------------------ //

type activeTaskTick time.Time

type Model struct {
	Tasks     []Task // for now, just a string
	Selected  int
	paginator paginator.Model
	ready     bool
	active    int
}

type LoadTasks struct {
	Tasks []Task
}

func loadTasks() tea.Msg {
	return LoadTasks{[]Task{
		{Title: "update PM for EOY", Priority: Low, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "pickup toys"},
		{Title: "update docs", Priority: None, Active: true},
		{Title: "message Sesha", Priority: High},
		{Title: "update PM for EOY", Priority: Low, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "pickup toys"},
		{Title: "update docs", Priority: None},
		{Title: "message Sesha", Priority: High},
		{Title: "update PM for EOY", Priority: Low, Due: "11/22/23"},
		{Title: "clean dishes and then take the trash out my dude", Due: "11/25/23"},
		{Title: "pickup toys"},
		{Title: "update docs", Priority: None},
		{Title: "message Sesha", Priority: High},
	}}
}

func (m Model) Init() tea.Cmd {
	return loadTasks
}

var (
	footerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	footerBoldStyle = footerStyle.Copy().Bold(true)
	headerKeyStyle  = lipgloss.NewStyle().Width(10).Align(lipgloss.Right).Italic(true).Foreground(lipgloss.Color("#FF06B7"))
)

func (m Model) View() string {
	// let's add a headerStyle with the active task info
	var b strings.Builder
	b.WriteString(m.header())
	b.WriteString(headerKeyStyle.Render("Tasks: "))
	b.WriteString("  " + m.paginator.View())
	b.WriteString(fmt.Sprintf("  (%d/%d)", m.paginator.Page+1, m.paginator.TotalPages))
	b.WriteString("\n")
	b.WriteString("\n")
	// start end bounds
	i, n := m.paginator.GetSliceBounds(len(m.Tasks))
	for index, task := range m.Tasks[i:n] {
		b.WriteString(fmt.Sprintf("%s\n", task.Render((index+i) == m.Selected)))
	}
	b.WriteString("\n")
	b.WriteString("\n")
	b.WriteString(m.footer())
	return b.String()
}

func (m Model) header() string {
	var b strings.Builder
	b.WriteString(headerKeyStyle.Render("Active:"))
	if m.active >= 0 {
		b.WriteString(" " + m.Tasks[m.active].Title)
	}
	b.WriteString("\n")
	b.WriteString(headerKeyStyle.Render("Runtime:"))
	if m.active >= 0 {
		b.WriteString(" " + m.Tasks[m.active].Runtime())
	}
	b.WriteString("\n")
	b.WriteString(headerKeyStyle.Render("Sub-Tasks:"))
	b.WriteString("\n")
	b.WriteString(headerKeyStyle.Render("Tags:"))
	b.WriteString("\n")
	b.WriteString("\n")

	return b.String()
}

func (m Model) footer() string {
	var b strings.Builder
	// TODO look into using the help bubble
	b.WriteString(footerStyle.Render(fmt.Sprintf("(%d/%d)", m.Selected+1, len(m.Tasks))))
	b.WriteString(fmt.Sprintf("  %s %s %s ", footerBoldStyle.Render("a"), footerStyle.Render("add"), footerBoldStyle.Render("\uF444")))
	b.WriteString(fmt.Sprintf("%s %s %s ", footerBoldStyle.Render("e"), footerStyle.Render("edit"), footerBoldStyle.Render("\uF444")))
	b.WriteString(fmt.Sprintf("%s %s %s ", footerBoldStyle.Render("t"), footerStyle.Render("toggle"), footerBoldStyle.Render("\uF444")))
	b.WriteString(fmt.Sprintf("%s %s %s ", footerBoldStyle.Render("c"), footerStyle.Render("complete"), footerBoldStyle.Render("\uF444")))
	return b.String()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case LoadTasks:
		p := paginator.New()
		p.Type = paginator.Dots
		p.PerPage = 10
		p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
		p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
		p.SetTotalPages(len(msg.Tasks))
		m.paginator = p
		m.ready = true
		m.Tasks = msg.Tasks
		m.active = activeTask(m.Tasks)
	case tea.KeyMsg:
		if m.ready {
			i, n := m.paginator.GetSliceBounds(len(m.Tasks))
			switch msg.String() {
			case "j":
				// down
				if m.Selected < (n - 1) {
					m.Selected++
				}
			case "k":
				if m.Selected > i {
					m.Selected--
				}
			case "t":
				// BUG fix timer, it should just be relative instead of increasing a value...
				// set selected task to be active
				if !hasActiveTask(m.Tasks) {
					m.Tasks[m.Selected].Active = true
					m.active = m.Selected
					// restart timer
					cmds = append(cmds, Tick())
				} else if m.Tasks[m.Selected].Active {
					m.Tasks[m.Selected].Active = false
					m.active = -1
				}
			case "enter":
				// this will enter detailed task view
			}
		}
	case activeTaskTick:
		if m.active >= 0 {
			m.Tasks[m.active].Tick()
			cmds = append(cmds, Tick())
		}
	}

	if m.ready {
		var cmd tea.Cmd
		m.paginator, cmd = m.paginator.Update(msg)
		cmds = append(cmds, cmd)
		i, n := m.paginator.GetSliceBounds(len(m.Tasks))
		if m.Selected < i || m.Selected >= n {
			m.Selected = i
		}
	}

	return m, tea.Batch(cmds...)
}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return activeTaskTick(t)
	})
}
