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

type ZettelDelegate struct{}

func (d ZettelDelegate) Height() int  { return 1 }
func (d ZettelDelegate) Spacing() int { return 0 }
func (d ZettelDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d ZettelDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	step, ok := item.(*domain.Step)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStep(step, index == m.Index()))
}

func ZettelsToListItems(zettels []domain.Zettel) []list.Item {
	items := make([]list.Item, len(zettels))
	for i := range zettels {
		item := &zettels[i]
		items[i] = item
	}
	return items
}

func renderZettel(zettel *domain.Zettel, selected bool) string {
	style := styles.DefaultItemStyle
	marker := ""
	cursorStyle := styles.InactiveCursorStyle
	if selected {
		style = styles.SelectedItemStyle
		marker = ">"
		cursorStyle = styles.ActiveCursorStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cursorStyle.Render(marker), style.Render(zettel.Name), style.Render(zettel.Tags))
}
