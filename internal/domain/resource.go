package domain

import "github.com/google/uuid"

const (
	WebLinkType = iota
	ZettelType  // TODO when you link a Zettel, you'll need to create some sort of lookup
	File
)

type ResourceType = int

type Resource struct {
	Id     string
	Name   string
	Source string
	Type   ResourceType
}

func NewResource(name, source string, resourceType ResourceType) Resource {
	id, err := uuid.NewRandom()
	if err != nil {
		return Resource{}
	}

	return Resource{
		Id:     id.String(),
		Name:   name,
		Source: source,
		Type:   resourceType,
	}
}

func (r *Resource) FilterValue() string {
	return r.Name
}
