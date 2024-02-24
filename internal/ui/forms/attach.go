package forms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

const (
	existingItem = "Existing"
	newItem      = "New"
)

func New() *AttachModel {
	selectedBox := styles.Box(styles.BoxOptions{
		BorderColor: styles.PrimaryColor,
		TextColor:   styles.SecondaryGray,
		BoxSize: styles.BoxSize{
			Width:  20,
			Height: 1,
		},
	})
	defaultBox := styles.Box(styles.BoxOptions{
		BorderColor: styles.SecondaryGray,
		TextColor:   styles.SecondaryGray,
		BoxSize: styles.BoxSize{
			Width:  20,
			Height: 1,
		},
	})

	return &AttachModel{
		selectedOptionStyle: selectedBox,
		defaultOptionStyle:  defaultBox,
		options:             []string{existingItem, newItem},
	}
}

// AttachModel is used for when
// you want to attach an item to another
type AttachModel struct {
	selectedOptionStyle lipgloss.Style
	defaultOptionStyle  lipgloss.Style
	options             []string
	active              int
	choice              string
}

func (m *AttachModel) Choice() string {
	return m.options[m.active]
}

func (m *AttachModel) Init() tea.Cmd {
	return nil
}

func (m *AttachModel) View() string {
	var b strings.Builder
	for i, choice := range m.options {
		if i == (m.active % len(m.options)) {
			b.WriteString(m.selectedOptionStyle.Render(choice))
		} else {
			b.WriteString(m.defaultOptionStyle.Render(choice))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *AttachModel) Update(msg tea.Msg) (*AttachModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.active++
		case "down", "j":
			m.active++
		case "enter":
			return m, choose(m.options[m.active%len(m.options)])
		}
	}

	return m, nil
}

type AttachTypeOptionMsg struct {
	Choice string
}

func choose(s string) tea.Cmd {
	return func() tea.Msg {
		return AttachTypeOptionMsg{Choice: s}
	}
}
