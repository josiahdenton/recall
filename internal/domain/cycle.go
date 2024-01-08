package domain

import (
	"github.com/google/uuid"
	"time"
)

type Cycle struct {
	Id                string    `json:"id"`
	Title             string    `json:"title"`
	AccomplishmentIds []string  `json:"accomplishment_ids"`
	Active            bool      `json:"active"`
	StartDate         time.Time `json:"start_date"`
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
