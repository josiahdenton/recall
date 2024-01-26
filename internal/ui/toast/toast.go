package toast

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"time"
)

var (
	statusStyle = styles.PrimaryGray.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D120AF")).Width(25).Align(lipgloss.Center)
)

type showToastMsg struct {
	message string
}

func ShowToast(message string) tea.Cmd {
	return func() tea.Msg {
		return showToastMsg{message: message}
	}
}

func New() Model {
	return Model{}
}

type Model struct {
	message string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if len(m.message) > 0 {
		return statusStyle.Render(m.message)
	}
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case showToastMsg:
		m.message = msg.message
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
