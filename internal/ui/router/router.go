package router

import tea "github.com/charmbracelet/bubbletea"

type Page = int

const (
	TaskListPage = iota
	TaskDetailedPage
)

type GotoPageMsg struct {
	Page Page
}

func GotoPage(page Page) tea.Cmd {
	return func() tea.Msg {
		return GotoPageMsg{Page: page}
	}
}
