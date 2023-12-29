package tasks

type Step struct {
	Description string
	Complete    bool
}

func (s *Step) FilterValue() string {
	return s.Description
}
