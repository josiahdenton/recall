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
	ActionPageBack

	// NOTE Cannot change these

	ActionSearch
	ActionShowKeybindings
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
)

func DefaultKeybindings() Keybindings {
	keys := make(map[string]Keybinding)
	keys["enter"] = Keybinding{
		Key:         "enter",
		Description: "Interact with the UI element",
		Action:      ActionInteract,
	}
	keys["space"] = Keybinding{
		Key:         "space",
		Description: "Copy a UI element's text",
		Action:      ActionCopy,
	}
	keys["tab"] = Keybinding{
		Key:         "tab",
		Description: "Move focus to the next UI element",
		Action:      ActionMoveFocusForward,
	}
	keys["a"] = Keybinding{
		Key:         "a",
		Description: "Add an item",
		Action:      ActionAdd,
	}
	keys["e"] = Keybinding{
		Key:         "e",
		Description: "Edit",
		Action:      ActionEdit,
	}
	keys["d"] = Keybinding{
		Key:         "d",
		Description: "Delete the selected item",
		Action:      ActionDelete,
	}
	keys["x"] = Keybinding{
		Key:         "x",
		Description: "Export (limited support)",
		Action:      ActionExport,
	}
	keys["c"] = Keybinding{
		Key:         "c",
		Description: "Complete/Archive the item",
		Action:      ActionArchive,
	}
	keys["f"] = Keybinding{
		Key:         "f",
		Description: "Favorite the selected item",
		Action:      ActionFavorite,
	}
	keys["u"] = Keybinding{
		Key:         "u",
		Description: "Undo the last io (limited support)",
		Action:      ActionUndo,
	}
	keys["r"] = Keybinding{
		Key:         "r",
		Description: "Restore items",
		Action:      ActionRestore,
	}
	keys["/"] = Keybinding{
		Key:         "/",
		Description: "Search/Filter items",
		Action:      ActionSearch,
	}
	keys["?"] = Keybinding{
		Key:         "?",
		Description: "Show the keybindings",
		Action:      ActionShowKeybindings,
	}
	keys["esc"] = Keybinding{
		Key:         "esc",
		Description: "Go back one page",
		Action:      ActionPageBack,
	}
	// keys["j"] = Keybinding{
	// 	Description: "Move up",
	// 	Action:      ActionUp,
	// }
	// keys["k"] = Keybinding{
	// 	Description: "Move down",
	// 	Action:      ActionDown,
	// }
	// keys["?"] = Keybinding{
	// 	Description: "Show the keybindings",
	// 	Action:      ActionShowKeybindings,
	// }
	// keys["?"] = Keybinding{
	// 	Description: "Show the keybindings",
	// 	Action:      ActionShowKeybindings,
	// }

	return Keybindings{
		Keys: keys,
	}
}

type Keybindings struct {
	Keys map[string]Keybinding `json:"keys"`
}

func (k Keybindings) FormatActionToKeyBind(action Action) Keybinding {
	for _, bind := range k.Keys {
		if bind.Action == action {
			return bind
		}
	}
	switch action {
	case ActionUp:
		return Keybinding{
			Key:         "k/\U000F005D",
			Description: "move up",
			Action:      ActionUp,
		}
	case ActionDown:
		return Keybinding{
			Key:         "j/\U000F0045",
			Description: "move down",
			Action:      ActionDown,
		}
	case ActionLeft:
		return Keybinding{
			Key:         "h/\U000F004D",
			Description: "move left",
			Action:      ActionLeft,
		}
	case ActionRight:
		return Keybinding{
			Key:         "j/\U000F0054",
			Description: "move right",
			Action:      ActionRight,
		}
	}
	return Keybinding{
		Description: "UNSUPPORTED",
		Action:      ActionNOP,
	}
}

func (k Keybindings) ParseKeyPress(keyPress string) Action {
	binding, ok := k.Keys[keyPress]
	if !ok {
		return ActionNOP
	}
	return binding.Action
}

type Keybinding struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Action      Action `json:"io"`
}
