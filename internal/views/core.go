package views

import tea "github.com/charmbracelet/bubbletea"

func New() Model {
	return Model{}
}

type Model struct {

	// Projects
	// ^ Projects will have Categories
	// tasks should have a "child" model that handles the logic here...
	// child models will use the same methods (
}

func (m Model) Init() tea.Cmd {
	// this call the Init from the active child model
	// I will need to use either tea.Batch or tea.Sequence
	return nil
}

func (m Model) View() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	}
	return nil, nil
}
