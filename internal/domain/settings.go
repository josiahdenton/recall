package domain

const (
	AccomplishmentsFileName = "accomplishments.json"
	TasksFileName           = "tasks.json"
	CyclesFileName          = "cycles.json"
	SettingsFileName        = "settings.json"
)

type Settings struct {
	Location string `json:"location"`
}
