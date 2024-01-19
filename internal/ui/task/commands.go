package task

import tea "github.com/charmbracelet/bubbletea"

type Action = int

const (
	Interact = iota
	Delete
	Add
	MoveFocus
	Back
	None
)

type Commands struct {
	keys map[string]Action
}

func DefaultCommands() Commands {
	keys := make(map[string]Action)
	keys["a"] = Add
	keys["enter"] = Interact
	keys["d"] = Delete
	keys["tab"] = MoveFocus
	keys["esc"] = Back
	return Commands{keys: keys}
}

func (c Commands) HandleInput(msg tea.KeyMsg) Action {
	action, ok := c.keys[msg.String()]
	if !ok {
		return None
	}
	return action
}
