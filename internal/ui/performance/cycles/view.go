package cycles

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/performance/cycles/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

type Model struct {
	cycles   list.Model
	ready    bool
	form     forms.CycleFormModel
	showForm bool
}

func New() Model {
	return Model{
		form: forms.NewCycleForm(),
	}
}

type loadCycles struct {
	cycles []list.Item
}

// LoadCycles must be called before this Model is ready
// core is responsible for this call
func LoadCycles(cycles []domain.Cycle) tea.Cmd {
	return func() tea.Msg {
		items := make([]list.Item, len(cycles))
		for i := range cycles {
			cycle := &cycles[i]
			items[i] = cycle
		}
		return loadCycles{cycles: items}
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
	case loadCycles:
		m.cycles = list.New(msg.cycles, cycleDelegate{}, 50, 20)
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
				cmd = router.GotoPage(domain.MenuPage, nil, "")
				cmds = append(cmds, cmd)
			}
		case tea.KeyEnter:
			if !m.showForm {
				cycle := m.cycles.SelectedItem().(*domain.Cycle)
				cmd = router.GotoPage(domain.AccomplishmentsPage, cycle, "")
				cmds = append(cmds, cmd)
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
