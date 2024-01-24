package resources

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
)

func New() Model {
	return Model{}
}

type Model struct {
	resources list.Model
	ready     bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.ready {
		return styles.WindowStyle.Render(m.resources.View())
	}
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case router.LoadPageMsg:
		resources := msg.State.([]domain.Resource)
		m.resources = list.New(toItemList(resources), resourceDelegate{}, 50, 30)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		m.resources.Styles.Title = styles.SecondaryGray.Copy()
		m.resources.SetShowHelp(false)
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true
	}
	if !m.ready {
		return m, nil
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.resources, cmd = m.resources.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			resource := m.resources.SelectedItem().(*domain.Resource)
			switch resource.Type {
			case domain.WebResource:
				resource.Open()
			default:
				// unsupported
			}
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}
	return m, tea.Batch(cmds...)
}

func toItemList(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}
