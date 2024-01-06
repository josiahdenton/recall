package domain

type MenuOption struct {
	Title string
	Page  Page
}

func (m *MenuOption) FilterValue() string {
	return ""
}
