package accomplishment

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
)

var (
	// TODO make these std across the whole app...
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

type loadAccomplishmentMsg struct {
	accomplishment domain.Accomplishment
}

func LoadAccomplishment(accomplishment domain.Accomplishment) tea.Cmd {
	return func() tea.Msg {
		return loadAccomplishmentMsg{accomplishment: accomplishment}
	}
}

type Model struct {
	accomplishment domain.Accomplishment
	tasks          list.Model
	ready          bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case router.GotoPageMsg:
		if msg.Page == domain.AccomplishmentPage {
			accomplishment := msg.Parameter.(domain.Accomplishment)
			m.accomplishment = accomplishment
			items := setupList(m.accomplishment.AssociatedTasks)
			m.tasks = list.New(items, shortTaskDelegate{}, 50, 10)
			m.tasks.SetShowStatusBar(false)
			m.tasks.SetFilteringEnabled(false)
			m.tasks.Title = "Related Tasks"
			m.tasks.Styles.PaginationStyle = paginationStyle
			m.tasks.Styles.Title = titleStyle
			m.tasks.SetShowHelp(false)
			m.tasks.KeyMap.Quit.Unbind()
			m.ready = true
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			task := m.tasks.SelectedItem().(*domain.Task)
			cmds = append(cmds, router.GotoPage(domain.TaskDetailedPage, *task, ""))
		}
	}

	return m, tea.Batch(cmds...)
}

func setupList(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i := range tasks {
		item := &tasks[i]
		items[i] = item
	}
	return items
}
