package toast

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"time"
)

type ShowToastMsg struct {
	Message string
	Toast   Type
}

const (
	Info = iota
	Warn
)

type Type = int

func ShowToast(message string, toast Type) tea.Cmd {
	return func() tea.Msg {
		return ShowToastMsg{Message: message, Toast: toast}
	}
}

func New() *Model {
	return &Model{}
}

type Model struct {
	message string
	toast   Type
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	if len(m.message) > 0 && m.toast == Info {
		return styles.InfoToastStyle.Render(m.message)
	} else if len(m.message) > 0 && m.toast == Warn {
		return styles.WarnToastStyle.Render(m.message)
	}
	return ""
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ShowToastMsg:
		m.message = msg.Message
		m.toast = msg.Toast
		cmd = clearStatus()
	case clearStatusMsg:
		m.message = ""
	}

	return m, cmd
}

type clearStatusMsg struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*5, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}
