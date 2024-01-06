package cycles

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type cycleDelegate struct{}

func (d cycleDelegate) Height() int  { return 1 }
func (d cycleDelegate) Spacing() int { return 0 }
func (d cycleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d cycleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	cycle, ok := item.(*domain.Cycle)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderCycleOption(cycle, index == m.Index()))
}
