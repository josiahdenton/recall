package zettel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

var (
	paginationStyle   = list.DefaultStyles().PaginationStyle
	titleStyle        = styles.SecondaryColor.Copy().PaddingLeft(8)
	defaultTitleStyle = styles.SecondaryGray.Copy().PaddingLeft(8)
	defaultListTitle  = styles.SecondaryGray.Copy()
	activeListTitle   = styles.SecondaryColor.Copy()
	// windows for zettel
	alignContent              = lipgloss.NewStyle().Width(100).Align(lipgloss.Center).PaddingRight(4)
	leftPad                   = lipgloss.NewStyle().PaddingLeft(8)
	defaultConceptWindowStyle = lipgloss.NewStyle().Padding(1).Width(80).Height(20).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	activeConceptWindowStyle  = lipgloss.NewStyle().Padding(1).Width(80).Height(20).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#fcd34d"))
)

const (
	header = iota
	content
	links
	resources
)

type section = int

func New() Model {
	return Model{
		zettelForm:     forms.NewZettelForm(),
		linkZettelForm: forms.NewLinkForm(),
		conceptForm:    forms.NewConceptForm(),
		resourceForm:   forms.NewResourceForm(),
	}
}

type Model struct {
	zettel         *domain.Zettel
	linkZettelForm tea.Model
	conceptForm    tea.Model
	resourceForm   tea.Model
	zettelForm     tea.Model
	links          list.Model
	resources      list.Model
	showForm       bool
	ready          bool
	active         section
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// TODO - tie in glamour for displaying the content
	var b strings.Builder
	if m.showForm && m.active == header {
		b.WriteString(m.zettelForm.View())
	} else if m.showForm && m.active == links {
		b.WriteString(m.linkZettelForm.View())
	} else if m.showForm && m.active == resources {
		b.WriteString(m.resourceForm.View())
	} else if m.showForm && m.active == content {
		b.WriteString(m.conceptForm.View())
	} else {
		if m.active == header {
			b.WriteString(titleStyle.Render(m.zettel.Name))
		} else {
			b.WriteString(defaultTitleStyle.Render(m.zettel.Name))
		}
		b.WriteString("\n")
		b.WriteString(defaultTitleStyle.Render(m.zettel.Tags))
		b.WriteString("\n\n")
		if m.active == content {
			b.WriteString(alignContent.Render(activeConceptWindowStyle.Render(m.zettel.Concept)))
		} else if m.active != content {
			b.WriteString(alignContent.Render(defaultConceptWindowStyle.Render(m.zettel.Concept)))
		}
		b.WriteString("\n\n")
		b.WriteString(leftPad.Render(m.links.View()))
		b.WriteString("\n")
		b.WriteString(leftPad.Render(m.resources.View()))
	}
	return styles.WindowStyle.Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	focusMoved := false

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		zettel := msg.State.(*domain.Zettel)
		m.zettel = zettel
		m.links = list.New(linksToItemList(m.zettel.Links), zettelDelegate{}, 80, 7)
		m.links.Title = "Links"
		m.links.Styles.PaginationStyle = paginationStyle
		if m.active == links {
			m.links.Styles.Title = activeListTitle
		} else {
			m.links.Styles.Title = defaultListTitle
		}
		m.links.SetShowHelp(false)
		m.links.SetShowStatusBar(false)
		m.links.KeyMap.Quit.Unbind()

		m.resources = list.New(resourcesToItemList(m.zettel.Resources), resourceDelegate{}, 80, 7)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		if m.active == resources {
			m.resources.Styles.Title = activeListTitle
		} else {
			m.resources.Styles.Title = defaultListTitle
		}
		m.resources.SetShowHelp(false)
		m.resources.SetShowStatusBar(false)
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true

	case forms.ConceptFormMsg:
		m.zettel.Concept = msg.Concept
		cmds = append(cmds, modifyZettel(*m.zettel))
	case forms.ResourceFormMsg:
		m.zettel.Resources = append(m.zettel.Resources, msg.Resource)
		m.resources.InsertItem(len(m.zettel.Resources), &m.zettel.Resources[len(m.zettel.Resources)-1])
		cmds = append(cmds, modifyZettel(*m.zettel))
	case forms.LinkFormMsg:
		m.zettel.Links = append(m.zettel.Links, &msg.Zettel)
		m.links.InsertItem(len(m.zettel.Links), m.zettel.Links[len(m.zettel.Links)-1])
		cmds = append(cmds, modifyZettel(*m.zettel))
	case state.SaveStateMsg:
		if msg.Type == state.ModifyZettel {
			m.showForm = false
		}
	}

	if !m.ready {
		return m, tea.Batch(cmds...)
	}

	if m.showForm && m.active == links {
		m.linkZettelForm, cmd = m.linkZettelForm.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.showForm && m.active == content {
		m.conceptForm, cmd = m.conceptForm.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.showForm && m.active == resources {
		m.resourceForm, cmd = m.resourceForm.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.showForm && m.active == header {
		m.zettelForm, cmd = m.zettelForm.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.showForm {
				m.showForm = false
			} else {
				cmds = append(cmds, router.GotoPreviousPage())
				m.active = header
			}
		}
	}

	if m.showForm {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.active = nextSection(m.active)
			focusMoved = true
		case tea.KeyEnter:
			if m.active == content {
				cmds = append(cmds, forms.AttachConcept(m.zettel.Concept))
				m.showForm = true
			} else if m.active == links && len(m.links.Items()) > 0 {
				selected := m.links.SelectedItem().(*domain.Zettel)
				cmds = append(cmds, router.GotoPage(domain.ZettelPage, selected.ID))
			} else if m.active == resources && len(m.resources.Items()) > 0 {
				selected := m.resources.SelectedItem().(*domain.Resource)
				selected.Open()
			}
		}

		if m.resources.FilterState() == list.Filtering || m.links.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "a": // add zettel
			if m.active == links || m.active == resources {
				m.showForm = true
			}
		case "e":
			if m.active == header {
				m.showForm = true
				cmds = append(cmds, forms.EditZettel(m.zettel))
			}
		case "d":
			if m.active == links {
				selected := m.links.SelectedItem().(*domain.Zettel)
				index := m.links.Index()
				m.zettel.Links = append(m.zettel.Links[:index], m.zettel.Links[index+1:]...)
				m.links.SetItems(linksToItemList(m.zettel.Links))
				cmds = append(cmds, unlinkZettel(m.zettel, selected), toast.ShowToast("unlinked zettel!", toast.Warn))
			} else if m.active == resources {
				selected := m.resources.SelectedItem().(*domain.Resource)
				index := m.resources.Index()
				m.zettel.Resources = append(m.zettel.Resources[:index], m.zettel.Resources[index+1:]...)
				m.resources.SetItems(resourcesToItemList(m.zettel.Resources))
				cmds = append(cmds, unlinkZettelResource(m.zettel, selected))
			}
		}
	}

	if focusMoved {
		switch m.active {
		case header:
			m.links.Styles.Title = defaultListTitle
			m.resources.Styles.Title = defaultListTitle
		case content:
			m.links.Styles.Title = defaultListTitle
			m.resources.Styles.Title = defaultListTitle
		case links:
			m.links.Styles.Title = activeListTitle
			m.resources.Styles.Title = defaultListTitle
		case resources:
			m.links.Styles.Title = defaultListTitle
			m.resources.Styles.Title = activeListTitle
		}
	}

	if m.active == links {
		m.links, cmd = m.links.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.active == resources {
		m.resources, cmd = m.resources.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func modifyZettel(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: zettel,
			Type:   state.ModifyZettel,
		}
	}
}

func unlinkZettel(parent *domain.Zettel, child *domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.ModifyLink,
			Parent: parent,
			Child:  child,
		}
	}
}

func unlinkZettelResource(parent *domain.Zettel, child *domain.Resource) tea.Cmd {
	return func() tea.Msg {
		return state.DeleteStateMsg{
			Type:   state.UnlinkZettelResource,
			Parent: parent,
			Child:  child,
		}
	}
}

func nextSection(section section) section {
	if section == header {
		return content
	} else if section == content {
		return links
	} else if section == links {
		return resources
	} else if section == resources {
		return header
	}
	return header
}

func linksToItemList(links []*domain.Zettel) []list.Item {
	items := make([]list.Item, len(links))
	for i := range links {
		items[i] = links[i]
	}
	return items
}

func resourcesToItemList(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}
