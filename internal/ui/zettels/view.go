package zettels

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/forms"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#3a3b5b"))
)

func New() Model {
	return Model{
		form: forms.NewZettelForm(),
	}
}

type Model struct {
	ready    bool
	zettels  list.Model
	showForm bool
	form     tea.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var b strings.Builder
	if m.showForm && m.ready {
		b.WriteString(m.form.View())
	} else if m.ready {
		b.WriteString(m.zettels.View())
	}
	return styles.WindowStyle.Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		zettels := msg.State.([]domain.Zettel)
		m.zettels = list.New(toItemList(zettels), zettelDelegate{}, 50, 30)
		m.zettels.Title = "Zettels"
		m.zettels.Styles.PaginationStyle = paginationStyle
		m.zettels.Styles.Title = titleStyle
		m.zettels.SetShowHelp(false)
		m.zettels.KeyMap.Quit.Unbind()
		m.ready = true
	case shared.SaveStateMsg:
		m.showForm = false
		if msg.Type == shared.ModifyZettel {
			zettel := msg.Update.(domain.Zettel)
			m.zettels.InsertItem(len(m.zettels.Items()), &zettel)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.showForm {
			m.showForm = false
		} else if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		}
	}

	if m.showForm && m.ready {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && len(m.zettels.Items()) > 1 {
			selected := m.zettels.SelectedItem().(*domain.Zettel)
			cmds = append(cmds, router.GotoPage(domain.ZettelPage, selected.ID))
		}
	}

	if m.zettels.FilterState() != list.Unfiltered {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.showForm = true
		case "d":
			// TODO
		}
	}

	return m, tea.Batch(cmds...)
}

func toItemList(zettels []domain.Zettel) []list.Item {
	items := make([]list.Item, len(zettels))
	for i := range zettels {
		item := &zettels[i]
		items[i] = item
	}
	return items
}
