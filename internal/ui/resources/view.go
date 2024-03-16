package resources

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = styles.SecondaryColorStyle.Copy()
)

func New(keyBinds domain.Keybindings) Model {
	return Model{
		keyBinds: keyBinds,
		form:     forms.NewResourceForm(),
	}
}

type Model struct {
	keyBinds  domain.Keybindings
	resources list.Model
	ready     bool
	form      tea.Model
	showForm  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.ready && m.showForm {
		return styles.WindowStyle.Render(m.form.View())
	} else if m.ready {
		return styles.WindowStyle.Render(m.resources.View())
	}
	return ""
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
		resources := msg.State.([]domain.Resource)
		m.resources = list.New(toItemList(resources), resourceDelegate{}, 90, 30)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		m.resources.Styles.Title = titleStyle
		m.resources.SetShowHelp(false)
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true
	case forms.ResourceFormMsg:
		cmds = append(cmds, saveResource(msg.Resource), router.RefreshPage())
		m.showForm = false
	}
	if !m.ready {
		return m, nil
	}

	if !m.showForm {
		m.resources, cmd = m.resources.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && !m.showForm {
			resource := m.resources.SelectedItem().(*domain.Resource)
			switch resource.Type {
			case domain.WebResource:
				resource.Open()
			default:
				// unsupported
			}
		} else if msg.Type == tea.KeyEsc {
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
			}
		} else if msg.String() == "e" && !m.showForm && m.resources.FilterState() != list.Filtering {
			selected, ok := m.resources.SelectedItem().(*domain.Resource)
			if ok {
				cmds = append(cmds, forms.EditResource(selected))
				m.showForm = true
			}
		} else if msg.String() == "a" && !m.showForm && m.resources.FilterState() != list.Filtering {
			//cmds = append(cmds, forms.EditResource(&domain.Resource{})) // TODO: I don't think I need this...
			m.showForm = true
		}
	}

	return m, tea.Batch(cmds...)
}

func saveResource(resource domain.Resource) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: resource,
			Type:   state.ModifyResource,
		}
	}
}

func toItemList(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}
