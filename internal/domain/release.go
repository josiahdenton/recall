package domain

import (
	"gorm.io/gorm"
	"time"
)

type Outcome = int

const (
	AwaitingRelease = iota
	SuccessfulRelease
	FailedRelease
)

type Release struct {
	gorm.Model
	Date       time.Time
	Owner      string
	Outcome    Outcome
	ArtifactID uint
}

func (r *Release) FilterValue() string {
	return ""
}
