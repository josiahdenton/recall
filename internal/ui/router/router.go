package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
)

// TODO increase the utility of this package

type GotoPageMsg struct {
	Page      domain.Page
	ParentId  string // the requested ID of the page / the parent ID??
	Parameter any    // any data needed for that page, attached from core
}

func GotoPage(page domain.Page, parameter any, id string) tea.Cmd {
	return func() tea.Msg {
		return GotoPageMsg{Page: page, Parameter: parameter, ParentId: id}
	}
}
