package tasks

type Step struct {
	Description string
}

func (s *Step) FilterValue() string {
	return s.Description
}
