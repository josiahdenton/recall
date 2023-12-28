package detailed

import "github.com/josiahdenton/recall/internal/pages/tasks"

func renderResource(r *tasks.Resource) string {
	return r.Name
}

func renderStatus(s *tasks.Status) string {
	return s.Description
}

func renderStep(s *tasks.Step) string {
	return s.Description
}
