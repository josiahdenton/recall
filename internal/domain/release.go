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
	Date          time.Time
	ReleaseChange string
	Owner         string
	Outcome       Outcome
	ArtifactID    uint
}

func (r *Release) OpenChange() bool {
	return openWebPage(r.ReleaseChange)
}

func (r *Release) FilterValue() string {
	return ""
}
