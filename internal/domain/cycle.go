package domain

import (
	"github.com/google/uuid"
	"time"
)

type Cycle struct {
	Id                string
	Title             string
	AccomplishmentIds []string
	Active            bool
	StartDate         time.Time
	accomplishments   []Accomplishment
}

func (c *Cycle) FilterValue() string {
	return c.Title
}

func (c *Cycle) Accomplishments() []Accomplishment {
	return c.accomplishments
}

func (c *Cycle) AttachAccomplishments(accomplishments []Accomplishment) {
	c.accomplishments = accomplishments
}

func NewCycle(title string, startDate time.Time) Cycle {
	id, err := uuid.NewRandom()
	if err != nil {
		return Cycle{}
	}

	return Cycle{
		Id:                id.String(),
		Title:             title,
		AccomplishmentIds: make([]string, 0),
		StartDate:         startDate,
	}

}
