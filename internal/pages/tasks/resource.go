package tasks

const (
	WebLinkType = iota
	ZettelType
	File
)

type Type = int

type Resource struct {
	Name   string
	Source string
	Type   Type
}

func (r *Resource) FilterValue() string {
	return r.Name
}
