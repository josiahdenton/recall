package shared

type ModifyType = int

const (
	ModifyTask ModifyType = iota
	ModifyStep
	ModifyResource
	ModifyStatus
	ModifyCycle
	ModifyZettel
	ModifyAccomplishment
	ModifySettings
)

// SaveStateMsg should be sent anytime any state is modified
type SaveStateMsg struct {
	Update any
	Type   ModifyType
}

type DeleteStateMsg struct {
	Type ModifyType
	ID   uint
}

type LoadType = int

const (
	LoadZettel = iota
)

type RequestStateMsg struct {
	Type LoadType
	ID   uint // if 0, we load everything
}

type LoadedStateMsg struct {
	State any
}
