package toast

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"time"
)

var (
	warnStatusStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF")).Width(25).Align(lipgloss.Center)
	infoStatusStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#2dd4bf")).Width(25).Align(lipgloss.Center)
)

type showToastMsg struct {
	message string
	toast   ToastType
}

const (
	Info = iota
	Warn
)

type ToastType = int

func ShowToast(message string, toast ToastType) tea.Cmd {
	return func() tea.Msg {
		return showToastMsg{message: message, toast: toast}
	}
}

func New() Model {
	return Model{}
}

type Model struct {
	message string
	toast   ToastType
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if len(m.message) > 0 && m.toast == Info {
		return infoStatusStyle.Render(m.message)
	} else if len(m.message) > 0 && m.toast == Warn {
		return warnStatusStyle.Render(m.message)
	}
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case showToastMsg:
		m.message = msg.message
		m.toast = msg.toast
		cmd = clearStatus()
	case clearStatusMsg:
		m.message = ""
	}

	return m, cmd
}

type clearStatusMsg struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*3, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}
