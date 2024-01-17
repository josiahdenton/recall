package menu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
)

var (
	pageOptions = []domain.MenuOption{
		{
			Title: "Tasks",
			Page:  domain.TaskListPage,
		},
		{
			Title: "Zettels",
			Page:  domain.ZettelsPage,
		},
		{
			Title: "Resources",
			Page:  domain.ResourcesPage,
		},
		{
			Title: "Performance Cycles",
			Page:  domain.CyclesPage,
		},
		{
			Title: "Releases",
			Page:  domain.CyclesPage,
		},
		{
			Title: "Artifacts",
			Page:  domain.CyclesPage,
		},
		{
			Title: "Settings",
			Page:  domain.CyclesPage,
		},
	}
	paginationStyle = list.DefaultStyles().PaginationStyle
)

func New() Model {
	pages := list.New(toItemList(pageOptions), menuDelegate{}, 50, 20)
	pages.SetShowStatusBar(false)
	pages.SetFilteringEnabled(false)
	pages.Title = "Recall"
	pages.Styles.PaginationStyle = paginationStyle
	pages.Styles.Title = styles.SecondaryColor.Copy()
	pages.SetShowHelp(false)
	pages.KeyMap.Quit.Unbind()

	return Model{pages: pages}
}

type Model struct {
	pages list.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return styles.WindowStyle.Render(m.pages.View())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.pages, cmd = m.pages.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			option := m.pages.SelectedItem().(*domain.MenuOption)
			cmds = append(cmds, router.GotoPage(option.Page, 0))
		}
	}

	return m, tea.Batch(cmds...)
}

func toItemList(options []domain.MenuOption) []list.Item {
	items := make([]list.Item, len(options))
	for i := range options {
		item := &options[i]
		items[i] = item
	}
	return items
}
