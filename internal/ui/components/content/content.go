package content

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/ui/services/router"
	"github.com/josiahdenton/recall/internal/ui/services/state"
	"github.com/josiahdenton/recall/internal/ui/services/user"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

type attachContentMsg struct {
	content string
}

func AttachContent(content string) tea.Cmd {
	return func() tea.Msg {
		return attachContentMsg{
			content: content,
		}
	}
}

func New(size styles.Size) Model {
	hoverContentStyle := styles.Box(styles.BoxOptions{
		Size:        size,
		BorderColor: styles.PrimaryColor,
	})
	defaultContentStyle := styles.Box(styles.BoxOptions{
		Size:        size,
		TextColor:   styles.SecondaryGray,
		BorderColor: styles.SecondaryColor,
	})

	return Model{
		mode:                state.View,
		hoverContentStyle:   hoverContentStyle,
		defaultContentStyle: defaultContentStyle,
	}
}

type Location = int

type Model struct {
	focused             bool
	content             string
	ready               bool
	mode                state.Mode
	hoverContentStyle   lipgloss.Style
	defaultContentStyle lipgloss.Style
	// form
	area textarea.Model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var b strings.Builder
	if m.mode == state.Edit {
		// ... write string for form
	} else {
		// TODO - leave it up to the parent to use boxes
		if m.focused {
			b.WriteString(m.hoverContentStyle.Render(m.content))
		} else {
			b.WriteString(m.defaultContentStyle.Render(m.content))
		}
	}

	return b.String()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	cmd = m.onGlobalEvent(msg)
	cmds = append(cmds, cmd)

	cmd = m.onLocalEvent(msg)
	cmds = append(cmds, cmd)

	if m.mode == state.View && m.focused {
		cmd = m.onInput(msg)
		cmds = append(cmds, cmd)
	} else if m.mode == state.Edit && m.focused {
		cmd = m.onFormInput(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

// onGlobalEvent for any event defined outside this component
func (m *Model) onGlobalEvent(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case state.ModeSwitchMsg:
		m.mode = msg.Current
	}

	return nil
}

// onLocalEvent for any event defined inside this component
func (m *Model) onLocalEvent(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case attachContentMsg:
		m.content = msg.content
		m.area.SetValue(msg.content)
	}
	return nil
}

// onInput used when view component is in view mode
func (m *Model) onInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			cmd = state.SwitchMode(state.Edit)
		case "space":
			cmd = user.Copy(m.content)
		case "enter":
			cmd = state.SwitchMode(state.Focus)
		}
	}

	return cmd
}

// onFormInput used when view component is in edit mode
func (m *Model) onFormInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// keyDC3: "ctrl+s",
		switch msg.String() {
		case "esc":
			//TODO - esc should switch us back to view mode
			// I need Blur and Focus so I can know whether
			// my component should react to a SwitchMode
			// the slightly extra delay is fine for "actions"
			// typing is the only thing that needs to be "instant"
			if m.mode == state.Edit {
				cmds = append(cmds, state.SwitchMode(state.View))
			} else {
				cmds = append(cmds, router.Back())
			}
		case "enter":
		case "ctrl+s":
		}
	}

	m.area, cmd = m.area.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) Reset() {
	m.content = ""
	m.ready = false
	m.mode = state.View
}
