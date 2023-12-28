package tasks

type Status struct {
	Description string
}

func (s *Status) FilterValue() string {
	return s.Description
}
