package cycles

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

type Model struct {
	keyBinds domain.Keybindings
	cycles   list.Model
	ready    bool
	form     forms.CycleFormModel
	showForm bool
}

func New(keyBinds domain.Keybindings) Model {
	return Model{
		keyBinds: keyBinds,
		form:     forms.NewCycleForm(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var s string
	if m.showForm {
		s = styles.WindowStyle.Render(m.form.View())
	} else {
		s = styles.WindowStyle.Render(m.cycles.View())
	}
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.showForm {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		cycles := msg.State.([]domain.Cycle)
		m.cycles = list.New(toItemList(cycles), cycleDelegate{}, 50, 20)
		m.cycles.SetShowStatusBar(false)
		m.cycles.SetFilteringEnabled(false)
		m.cycles.Title = "Performance Cycles"
		m.cycles.Styles.PaginationStyle = paginationStyle
		m.cycles.Styles.Title = titleStyle
		m.cycles.SetShowHelp(false)
		m.cycles.KeyMap.Quit.Unbind()
		m.ready = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.showForm {
				m.showForm = false
			} else {
				cmd = router.GotoPage(domain.MenuPage, 0)
				cmds = append(cmds, cmd)
			}
		case tea.KeyEnter:
			if !m.showForm {
				cycle := m.cycles.SelectedItem().(*domain.Cycle)
				cmd = router.GotoPage(domain.AccomplishmentsPage, cycle.ID)
				cmds = append(cmds, cmd)
			}
		case tea.KeySpace:
			if !m.showForm {
				cycle := m.cycles.SelectedItem().(*domain.Cycle)
				cmds = append(cmds, toggleCycleActive(cycle))
			}
		}
		if !m.showForm {
			switch msg.String() {
			case "a":
				m.showForm = true
			}
		}
	}

	if m.ready && !m.showForm {
		m.cycles, cmd = m.cycles.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func toItemList(cycles []domain.Cycle) []list.Item {
	items := make([]list.Item, len(cycles))
	for i := range cycles {
		item := &cycles[i]
		items[i] = item
	}
	return items
}

func toggleCycleActive(cycle *domain.Cycle) tea.Cmd {
	cycle.Active = !cycle.Active
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: *cycle,
			Type:   state.ModifyCycle,
		}
	}
}
