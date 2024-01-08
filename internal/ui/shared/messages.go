package shared

type UpdateType = int

const (
	TaskUpdate UpdateType = iota
	StepUpdate
	ResourceUpdate
	StatusUpdate
	CycleUpdate
	AccomplishmentUpdate
	SettingsUpdate
)

// SaveStateMsg should be sent anytime any state is modified
type SaveStateMsg struct {
	// Update always has a value, not a pointer
	Update   any
	Type     UpdateType
	ParentId string
}

type LoadRepositoryMsg struct{}
