package shared

type UpdateType = int

const (
	ModifyTask UpdateType = iota
	ModifyStep
	ModifyResource
	ModifyStatus
	ModifyCycle
	ModifyAccomplishment
	ModifySettings
)

// SaveStateMsg should be sent anytime any state is modified
type SaveStateMsg struct {
	Update any
	Type   UpdateType
}

type LoadRepositoryMsg struct{}

type DeleteStateMsg struct {
	Type UpdateType
	ID   uint
}
