package forms

import tea "github.com/charmbracelet/bubbletea"

type Form interface {
	Init() tea.Cmd
	View() string
	Update(msg tea.Msg) (Form, tea.Cmd)
	Reset()
}
