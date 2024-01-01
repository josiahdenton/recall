package domain

type Accomplishment struct {
	Description string
	Impact      string
	Strength    string
}

func (a *Accomplishment) FilterValue() string {
	return a.Description
}
