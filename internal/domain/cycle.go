package domain

import (
	"github.com/google/uuid"
	"time"
)

type Cycle struct {
	Id              string
	Title           string
	Accomplishments []Accomplishment
	Active          bool
	StartDate       time.Time
}

func (c *Cycle) FilterValue() string {
	return c.Title
}

func NewCycle(title string, startDate time.Time) Cycle {
	id, err := uuid.NewRandom()
	if err != nil {
		return Cycle{}
	}

	return Cycle{
		Id:              id.String(),
		Title:           title,
		Accomplishments: make([]Accomplishment, 0),
		StartDate:       startDate,
	}

}
