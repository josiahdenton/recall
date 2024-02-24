package user

import (
	tea "github.com/charmbracelet/bubbletea"
)

type copyMsg struct {
	content string
}

// Copy will copy the string s to the clipboard
func Copy(s string) tea.Cmd {
	return func() tea.Msg {
		return copyMsg{content: s}
	}
}
