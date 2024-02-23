package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"github.com/josiahdenton/recall/internal/ui/styles"
	"io"
)

type StepDelegate struct{}

func (d StepDelegate) Height() int  { return 1 }
func (d StepDelegate) Spacing() int { return 0 }
func (d StepDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d StepDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	step, ok := item.(*domain.Step)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStep(step, index == m.Index()))
}

func StepsToListItems(steps []domain.Step) []list.Item {
	items := make([]list.Item, len(steps))
	for i := range steps {
		item := &steps[i]
		items[i] = item
	}
	return items
}

func renderStep(step *domain.Step, selected bool) string {
	style := styles.DefaultItemStyle
	if selected {
		style = styles.SelectedItemStyle
	}
	stepMarker := "\uE640"
	if step.Complete {
		stepMarker = "\U000F0856"
	}
	return style.Render(fmt.Sprintf("%s %s", stepMarker, step.Description))
}
