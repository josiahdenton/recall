package zettel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	forms2 "github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

var (
	paginationStyle  = list.DefaultStyles().PaginationStyle
	titleStyle       = styles.PrimaryColor.Copy().Align(lipgloss.Center)
	defaultListTitle = styles.SecondaryGray.Copy()
	activeListTitle  = styles.PrimaryColor.Copy()
	// windows for zettel concepts
	defaultConceptWindowStyle = lipgloss.NewStyle().Padding(2).Width(80).Height(20).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3a3b5b"))
	activeConceptWindowStyle  = lipgloss.NewStyle().Padding(2).Width(80).Height(20).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF"))
)

const (
	content = iota
	links
	resources
)

type section = int

func New() Model {
	return Model{
		zettelForm:  forms2.NewLinkForm(),
		conceptForm: forms2.NewConceptForm(),
	}
}

type Model struct {
	zettel      *domain.Zettel
	zettelForm  tea.Model
	conceptForm tea.Model
	links       list.Model
	resources   list.Model
	showForm    bool
	ready       bool
	active      section
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// TODO - tie in glamour for displaying the content
	var b strings.Builder
	if m.showForm {
		b.WriteString(m.zettelForm.View())
	} else {
		b.WriteString(titleStyle.Render(m.zettel.Name))
		b.WriteString("\n")
		if m.active == content {
			b.WriteString(activeConceptWindowStyle.Render(m.zettel.Concept))
		} else {
			b.WriteString(defaultConceptWindowStyle.Render(m.zettel.Concept))
		}
		b.WriteString("\n")
		b.WriteString(m.links.View())
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
		m.links.Styles.Title = defaultListTitle
		m.links.SetShowHelp(false)
		m.links.SetFilteringEnabled(false)
		m.links.SetShowStatusBar(false)
		m.links.KeyMap.Quit.Unbind()

		m.resources = list.New(resourcesToItemList(m.zettel.Resources), zettelDelegate{}, 80, 5)
		m.resources.Title = "Resources"
		m.resources.Styles.PaginationStyle = paginationStyle
		m.resources.Styles.Title = defaultListTitle
		m.resources.SetShowHelp(false)
		m.resources.SetFilteringEnabled(false)
		m.resources.SetShowStatusBar(false)
		m.resources.KeyMap.Quit.Unbind()
		m.ready = true

	case forms2.ConceptFormMsg:
		m.zettel.Concept = msg.Concept
		cmds = append(cmds, modifyZettel(*m.zettel))
	case forms2.LinkFormMsg:
		// save this zettel with new link...
		m.zettel.Links = append(m.zettel.Links, &msg.Zettel)
		m.links.InsertItem(len(m.zettel.Links), m.zettel.Links[len(m.zettel.Links)-1])
		cmds = append(cmds, modifyZettel(*m.zettel))
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.showForm {
			m.showForm = false
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}

	if !m.ready {
		return m, tea.Batch(cmds...)
	}

	if m.showForm && m.active == links {
		m.zettelForm, cmd = m.zettelForm.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	} else if m.showForm && m.active == content {
		m.conceptForm, cmd = m.conceptForm.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.active == links {
		m.links, cmd = m.links.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.active == resources {
		m.resources, cmd = m.resources.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.active = nextSection(m.active)
			focusMoved = true
		case tea.KeyShiftTab:
			m.active = nextSection(m.active)
			focusMoved = true
		case tea.KeyEnter:
			if m.active == content {
				cmds = append(cmds, forms2.AttachConcept(m.zettel.Concept))
				m.showForm = true
			} else if m.active == links {
				selected := m.links.SelectedItem().(*domain.Zettel)
				cmds = append(cmds, router.GotoPage(domain.ZettelPage, selected.ID))
			}
		}

		switch msg.String() {
		case "a": // add zettel
			if m.active == links {
				m.showForm = true
			}
		}
	}

	if focusMoved {
		switch m.active {
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

	return m, tea.Batch(cmds...)
}

func modifyZettel(zettel domain.Zettel) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Update: zettel,
			Type:   shared.ModifyZettel,
		}
	}
}

func nextSection(section section) section {
	if section == content {
		return links
	}
	return content
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