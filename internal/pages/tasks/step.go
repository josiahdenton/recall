package tasks

type Step struct {
	Description string
	Complete    bool
}

func (s *Step) FilterValue() string {
	return s.Description
}

func (s *Step) ToggleStatus() {
	s.Complete = !s.Complete
}
