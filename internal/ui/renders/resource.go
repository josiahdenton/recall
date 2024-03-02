package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
)

type ResourceDelegate struct{}

func (d ResourceDelegate) Height() int  { return 1 }
func (d ResourceDelegate) Spacing() int { return 0 }
func (d ResourceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d ResourceDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	resource, ok := item.(*domain.Resource)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderResource(resource, index == m.Index()))
}

func ResourcesToListItems(resources []domain.Resource) []list.Item {
	items := make([]list.Item, len(resources))
	for i := range resources {
		item := &resources[i]
		items[i] = item
	}
	return items
}

func renderResource(resource *domain.Resource, selected bool) string {
	style := styles.DefaultItemStyle
	marker := ""
	cursorStyle := styles.InactiveCursorStyle
	if selected {
		style = styles.SelectedItemStyle
		marker = ">"
		cursorStyle = styles.ActiveCursorStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(marker), style.Render(resource.Name), style.Render(resource.Tags))
}
