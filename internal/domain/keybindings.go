package domain

type Action = int

const (
	ActionNOP Action = iota
	ActionInteract
	ActionCopy
	ActionMoveFocusForward
	ActionAdd
	ActionEdit
	ActionDelete
	ActionExport
	ActionArchive
	ActionFavorite
	ActionUndo
	ActionRestore
	ActionSearch
	ActionShowKeybindings
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
)

type Keybindings struct {
	Keys []Keybinding `json:"keys"`
}

func (k Keybindings) ParseKeyPress(keyPress string) (Action, bool) {
	for _, binding := range k.Keys {
		if binding.Key == keyPress {
			return binding.Action, true
		}
	}
	return ActionNOP, false
}

type Keybinding struct {
	Key    string `json:"key"`
	Action Action `json:"action"`
}
