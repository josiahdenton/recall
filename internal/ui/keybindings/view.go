package keybindings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/router"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"strings"
)

type SaveKeyBindingsMsg struct {
	keybindings domain.Keybindings
}

var (
	keyStyle               = styles.AccentColor.Copy().Width(6).Bold(true)
	actionDescriptionStyle = styles.SecondaryGray.Copy()
)

var actions = []domain.Action{
	domain.ActionInteract,
	domain.ActionMoveFocusForward,
	domain.ActionCopy,
	domain.ActionPageBack,
	domain.ActionAdd,
	domain.ActionDelete,
	domain.ActionEdit,
	domain.ActionExport,
	domain.ActionArchive,
	domain.ActionFavorite,
	domain.ActionUndo,
	domain.ActionRestore,
	domain.ActionSearch,
	domain.ActionShowKeybindings,
	domain.ActionUp,
	domain.ActionDown,
	domain.ActionLeft,
	domain.ActionRight,
}

func New(keybindings domain.Keybindings) Model {
	return Model{
		keybindings: keybindings,
	}
}

type Model struct {
	ready       bool
	keybindings domain.Keybindings
}

// tab to switch around focus...
// inputs locked until we hit "Enter" on that input
// sets value to preset "[]"
// "Esc" will put the value back to the way it was
// "Enter" on the "[]" state will throw an error

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var b strings.Builder
	for _, action := range actions {
		bind := m.keybindings.FormatActionToKeyBind(action)
		b.WriteString(keyStyle.Render(bind.Key))
		b.WriteString(actionDescriptionStyle.Render(bind.Description))
		b.WriteString("\n\n")
	}
	return styles.WindowStyle.Render(b.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			cmds = append(cmds, router.GotoPreviousPage())
		} else if msg.String() == "?" {
			cmds = append(cmds, router.GotoPreviousPage())
		}
	}

	return m, tea.Batch(cmds...)
}
