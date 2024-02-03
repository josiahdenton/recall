package keybindings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
)

var order = []domain.Action{
	domain.ActionInteract,
	domain.ActionCopy,
	domain.ActionAdd,
	domain.ActionDelete,
	domain.ActionEdit,
	domain.ActionMoveFocusForward,
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

func New() Model {
	return Model{}
}

type Model struct {
	inputs map[domain.Action]tea.Model
}
