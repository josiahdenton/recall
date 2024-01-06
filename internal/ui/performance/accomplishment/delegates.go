package accomplishment

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type shortTaskDelegate struct{}

func (d shortTaskDelegate) Height() int  { return 1 }
func (d shortTaskDelegate) Spacing() int { return 0 }
func (d shortTaskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d shortTaskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	task, ok := item.(*domain.Task)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderMinimumTask(task, index == m.Index()))
}
