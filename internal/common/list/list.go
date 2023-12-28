package list

//
//import (
//	"fmt"
//	"github.com/charmbracelet/bubbles/paginator"
//	tea "github.com/charmbracelet/bubbletea"
//	"strings"
//)
//
//type Renderable interface {
//	Render(selected bool) string
//	Select()
//}
//
//type Model struct {
//	paginator paginator.Model
//	items []Renderable
//}
//
//// TODO add option to include callback here...
//
//func (m *Model) Init() tea.Cmd {
//	return nil
//}
//
//func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	// TODO add the ability to move across status updates, and edit them
//	var cmd tea.Cmd
//	m.paginator, cmd = m.paginator.Update(msg)
//	return m, cmd
//}
//
//func (m *Model) View() string {
//	// TODO I should generalize this list...
//	var b strings.Builder
//	b.WriteString("Status Updates ")
//	b.WriteString("  " + m.paginator.View())
//	b.WriteString(fmt.Sprintf("  (%d/%d)", m.paginator.Page+1, m.paginator.TotalPages))
//	b.WriteString("\n")
//	b.WriteString("\n")
//	i, n := m.paginator.GetSliceBounds(len(m.Updates))
//	for _, update := range m.Updates[i:n] {
//		b.WriteString(update.Render())
//	}
//	b.WriteString("\n")
//	b.WriteString("\n")
//	return ""
//}
