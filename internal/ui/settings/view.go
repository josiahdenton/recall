package settings

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/shared"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

var (
	titleStyle     = styles.PrimaryColor.Copy()
	formLabelStyle = styles.SecondaryGray.Copy()
	errorStyle     = styles.PrimaryColor.Copy()
)

const (
	location = iota
	// add more here
	inputCount
)

type Model struct {
	settings *domain.Settings
	inputs   []textinput.Model
	ready    bool
}

func New() Model {

	inputs := make([]textinput.Model, inputCount)
	inputs[location] = textinput.New()
	inputs[location].Focus()
	inputs[location].Width = 60
	inputs[location].CharLimit = 60
	inputs[location].Prompt = "Location: "
	inputs[location].PromptStyle = formLabelStyle
	inputs[location].Placeholder = "..." // TODO replace with existing value

	return Model{inputs: inputs}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString(titleStyle.Render("Settings"))
	for _, input := range m.inputs {
		builder.WriteString(input.View())
	}
	return styles.WindowStyle.Render(builder.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		settings := msg.State.(*domain.Settings)
		m.settings = settings
		m.inputs[location].Placeholder = m.settings.Location
		m.ready = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, saveSettings(*m.settings), router.GotoPage(domain.MenuPage, 0))
		}
	}

	return m, cmd
}

func saveSettings(settings domain.Settings) tea.Cmd {
	return func() tea.Msg {
		return shared.SaveStateMsg{
			Update: settings,
			Type:   shared.ModifySettings,
		}
	}
}
