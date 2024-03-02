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

type statusDelegate struct{}

func (d statusDelegate) Height() int  { return 1 }
func (d statusDelegate) Spacing() int { return 0 }
func (d statusDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d statusDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	step, ok := item.(*domain.Step)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStep(step, index == m.Index()))
}

func StatusToListItems(status []domain.Status) []list.Item {
	items := make([]list.Item, len(status))
	for i := range status {
		item := &status[i]
		items[i] = item
	}
	return items
}

func renderStatus(status *domain.Status, selected bool) string {
	style := styles.DefaultItemStyle
	marker := ""
	cursorStyle := styles.InactiveCursorStyle
	if selected {
		style = styles.SelectedItemStyle
		marker = ">"
		cursorStyle = styles.ActiveCursorStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(marker), style.Render(status.Name), style.Render(status.Tags))
}
