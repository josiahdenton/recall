package domain

import (
	"gorm.io/gorm"
	"time"
)

type Release struct {
	gorm.Model
	Date       time.Time
	Owner      string
	Steps      []Step
	Outcome    Status
	Completed  bool
	ArtifactID uint
}

func (r *Release) FilterValue() string {
	return ""
}

func (r *Release) ToggleComplete() {
	r.Completed = !r.Completed
}
