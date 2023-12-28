package resources

import (
	"fmt"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

const (
	WebLink = iota
	File
	Zettel
)

type Type = int

type Model struct {
	paginator paginator.Model
	Resources []Resource
}

type Resource struct {
	Name string
	Type Type
	Link string
}

func (r *Resource) Render() string {
	return r.Name
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO add the ability to move across status updates, and edit them
	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	// TODO I should generalize this list...
	var b strings.Builder
	b.WriteString("Status Updates ")
	b.WriteString("  " + m.paginator.View())
	b.WriteString(fmt.Sprintf("  (%d/%d)", m.paginator.Page+1, m.paginator.TotalPages))
	b.WriteString("\n")
	b.WriteString("\n")
	i, n := m.paginator.GetSliceBounds(len(m.Resources))
	for _, resource := range m.Resources[i:n] {
		b.WriteString(resource.Render())
	}
	b.WriteString("\n")
	b.WriteString("\n")
	return ""
}
