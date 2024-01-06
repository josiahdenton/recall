package shared

type UpdateType = int

const (
	TaskUpdate UpdateType = iota
	StepUpdate
	ResourceUpdate
	StatusUpdate
	CycleUpdate
	AccomplishmentUpdate
)

// SaveStateMsg should be sent anytime any state is modified
type SaveStateMsg struct {
	Update   any
	Type     UpdateType
	ParentId string
}

type SelectedType = int

const (
	TaskSelected SelectedType = iota
)

// SelectedItemMsg is for passing down selected items to forms
type SelectedItemMsg struct {
	Item any
	Type SelectedType
}
