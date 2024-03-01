package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
)

type StatusDelegate struct{}

func (d StatusDelegate) Height() int  { return 1 }
func (d StatusDelegate) Spacing() int { return 0 }
func (d StatusDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d StatusDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	status, ok := item.(*domain.Status)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStatus(status, index == m.Index()))
}

func StatusToListItems(status []domain.Status) []list.Item {
	items := make([]list.Item, len(status))
	for i := range status {
		item := &status[i]
		items[i] = item
	}
	return items
}

var (
	selectedStatusStyle = styles.Box(styles.BoxOptions{
		BorderColor: styles.PrimaryColor,
		BoxSize: styles.BoxSize{
			Width: 60,
		},
	})
	defaultStatusStyle = styles.Box(styles.BoxOptions{
		BorderColor: styles.PrimaryColor,
		BoxSize: styles.BoxSize{
			Width: 60,
		},
	})
)

func renderStatus(status *domain.Status, selected bool) string {
	style := defaultStatusStyle
	if selected {
		style = selectedStatusStyle
	}

	return style.Render(status.Description)
}
