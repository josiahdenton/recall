package state

import tea "github.com/charmbracelet/bubbletea"

type ModifyType = int

const (
	ModifyTask ModifyType = iota
	ModifyStep
	ModifyResource
	ModifyStatus
	ModifyCycle
	ModifyZettel
	ModifyLink
	ModifyAccomplishment
	ModifySettings
	ModifyArtifact
	ModifyRelease
	UnlinkTaskStep
	UnlinkTaskResource
	UnlinkTaskStatus
	UnlinkZettelResource
)

// TODO - for now, undo only supported for deleting a whole task/zettel/accomplishment/cycle
// TODO - state stuff may need it's own

type UndoDeleteStateMsg struct{}

type History struct {
	Deletes []DeleteStateMsg
}

// SaveStateMsg should be sent anytime any state is modified
type SaveStateMsg struct {
	Update any
	Type   ModifyType
}

type DeleteStateMsg struct {
	Type ModifyType
	ID   uint
	// Parent of the association removal
	Parent any
	// Child of the association removal
	Child any
}

type LoadType = int

const (
	LoadZettel = iota
	LoadResource
	LoadCycle
)

type RequestStateMsg struct {
	Type LoadType
	ID   uint // if 0, we load everything
}

type LoadedStateMsg struct {
	State any
}

// RequestState grabs state from the repository and
// sends back a LoadedStateMsg. If "ID" is 0, RequestState
// will give back an array of those items
func RequestState(loadType LoadType, id uint) tea.Cmd {
	return func() tea.Msg {
		return RequestStateMsg{
			Type: loadType,
			ID:   id,
		}
	}
}

// UndoDeleteState will undo the last deleted item
// supports tasks, zettels, accomplishments.
func UndoDeleteState() tea.Cmd {
	return func() tea.Msg {
		return UndoDeleteStateMsg{}
	}
}
