package tasks

type Step struct {
	Title     string
	Completed bool
}

func (s *Step) Render() string {
	return s.Title
}
