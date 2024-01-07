package accomplishments

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"log"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

type Model struct {
	ready           bool
	accomplishments list.Model
	// you cannot add, this will all be just read only
	// will need to add an export option...
	// maybe have an option to choose to delete?
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return styles.WindowStyle.Render(m.accomplishments.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case router.LoadPageMsg:
		cycle := msg.State.(*domain.Cycle)
		m.accomplishments = list.New(toItemList(cycle.Accomplishments()), accomplishmentDelegate{}, 50, 20)
		m.accomplishments.SetShowStatusBar(false)
		m.accomplishments.SetFilteringEnabled(false)
		m.accomplishments.Title = "Accomplishments"
		m.accomplishments.Styles.PaginationStyle = paginationStyle
		m.accomplishments.Styles.Title = titleStyle
		m.accomplishments.SetShowHelp(false)
		m.accomplishments.KeyMap.Quit.Unbind()
		m.ready = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// go into accomplishment details
			accomplishment := m.accomplishments.SelectedItem().(*domain.Accomplishment)
			cmds = append(cmds, router.GotoPage(domain.AccomplishmentPage, accomplishment.Id))
		case tea.KeyEsc:
			log.Printf("tea esc")
			cmd = router.GotoPage(domain.CyclesPage, "")
			cmds = append(cmds, cmd)
		}
	}

	if m.ready {
		m.accomplishments, cmd = m.accomplishments.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func toItemList(accomplishments []domain.Accomplishment) []list.Item {
	items := make([]list.Item, len(accomplishments))
	for i := range accomplishments {
		item := &accomplishments[i]
		items[i] = item
	}
	return items
}
