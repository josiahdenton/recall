package settings

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/state"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"github.com/josiahdenton/recall/internal/ui/toast"
	"strings"
)

var (
	titleStyle     = styles.PrimaryColor.Copy()
	formLabelStyle = styles.SecondaryGray.Copy()
)

const (
	editor = iota
)

type Model struct {
	settings *domain.Settings
	inputs   []textinput.Model
	active   int
	ready    bool
}

func New() Model {
	inputs := make([]textinput.Model, 1)
	inputs[editor] = textinput.New()
	inputs[editor].Focus()
	inputs[editor].Width = 60
	inputs[editor].CharLimit = 60
	// Override with Artifact
	inputs[editor].Prompt = "Default Editor: "
	inputs[editor].PromptStyle = formLabelStyle
	inputs[editor].Placeholder = "(nvim)"

	return Model{inputs: inputs}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Settings"))
	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	return styles.WindowStyle.Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case router.LoadPageMsg:
		settings := msg.State.(*domain.Settings)
		m.settings = settings
		m.inputs[editor].SetValue(m.settings.Editor)
		m.ready = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			cmds = append(cmds, router.GotoPage(domain.MenuPage, 0))
		case tea.KeyEnter:
			if cmd := validateSettings(m.inputs[editor].Err); cmd != nil {
				return m, tea.Batch(cmds...)
			}

			cmds = append(cmds, saveSettings(*m.settings), router.GotoPage(domain.MenuPage, 0))
		case tea.KeyTab:
			m.inputs[m.active%len(m.inputs)].Blur()
			m.active++
			m.inputs[m.active%len(m.inputs)].Focus()
		}
	}

	m.inputs[m.active%len(m.inputs)], cmd = m.inputs[m.active%len(m.inputs)].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func validateSettings(errs ...error) tea.Cmd {
	for _, err := range errs {
		if err != nil {
			return toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn)
		}
	}
	return nil
}

func saveSettings(settings domain.Settings) tea.Cmd {
	return func() tea.Msg {
		return state.SaveStateMsg{
			Update: settings,
			Type:   state.ModifySettings,
		}
	}
}
