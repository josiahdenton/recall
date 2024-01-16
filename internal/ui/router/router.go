package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
)

// TODO increase the utility of this package

type GotoPageMsg struct {
	Page            domain.Page
	RequestedItemId uint // the requested ID of the page / the parent ID??
}

type History struct {
	Pages []GotoPageMsg
}

type LoadPageMsg struct {
	Page  domain.Page
	State any // any data needed for that page, attached from core
}

type PreviousPageMsg struct{}

func GotoPage(page domain.Page, id uint) tea.Cmd {
	return func() tea.Msg {
		return GotoPageMsg{Page: page, RequestedItemId: id}
	}
}

func GotoPreviousPage() tea.Cmd {
	return func() tea.Msg {
		return PreviousPageMsg{}
	}
}
