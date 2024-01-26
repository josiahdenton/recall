package accomplishment

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

var (
	// TODO make these std across the whole app...
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = styles.PrimaryColor.Copy().PaddingLeft(1)
	fadedTitleStyle = styles.SecondaryGray.Copy().Width(16).Align(lipgloss.Right)
	headerStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b")).Width(80)
)

type Model struct {
	accomplishment *domain.Accomplishment
	tasks          list.Model
	ready          bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(fadedTitleStyle.Render("What: "))
	b.WriteString("\n\n")
	b.WriteString(fadedTitleStyle.Render("Impact: "))
	b.WriteString("\n\n")
	b.WriteString(fadedTitleStyle.Render("Strength: "))
	left := b.String()
	b.Reset()

	b.WriteString(titleStyle.Render(m.accomplishment.Description))
	b.WriteString("\n\n")
	b.WriteString(titleStyle.Render(m.accomplishment.Impact))
	b.WriteString("\n\n")
	b.WriteString(titleStyle.Render(m.accomplishment.Strength))
	right := b.String()
	b.Reset()

	top := lipgloss.JoinHorizontal(lipgloss.Right, left, right)
	header := headerStyle.Render(top)

	b.WriteString("\n\n")
	b.WriteString(m.tasks.View())
	return styles.WindowStyle.Render(lipgloss.JoinVertical(lipgloss.Top, header, b.String()))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.ready {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		accomplishment := msg.State.(*domain.Accomplishment)
		m.accomplishment = accomplishment
		m.tasks = list.New(toItemList(m.accomplishment.Tasks), shortTaskDelegate{}, 50, 10)
		m.tasks.SetShowStatusBar(false)
		m.tasks.SetFilteringEnabled(false)
		m.tasks.Title = "Related Tasks"
		m.tasks.Styles.PaginationStyle = paginationStyle
		m.tasks.Styles.Title = fadedTitleStyle
		m.tasks.SetShowHelp(false)
		m.tasks.KeyMap.Quit.Unbind()
		m.ready = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			task := m.tasks.SelectedItem().(*domain.Task)
			cmds = append(cmds, router.GotoPage(domain.TaskDetailedPage, task.ID))
		case tea.KeyEsc:
			cmds = append(cmds, router.GotoPreviousPage())
		}
	}

	return m, tea.Batch(cmds...)
}

func toItemList(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i := range tasks {
		item := &tasks[i]
		items[i] = item
	}
	return items
}
