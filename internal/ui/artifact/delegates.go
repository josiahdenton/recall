package artifact

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type resourceDelegate struct{}

func (d resourceDelegate) Height() int  { return 1 }
func (d resourceDelegate) Spacing() int { return 0 }
func (d resourceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d resourceDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	resource, ok := item.(*domain.Resource)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderResource(resource, index == m.Index()))
}

type releaseDelegate struct{}

func (d releaseDelegate) Height() int  { return 1 }
func (d releaseDelegate) Spacing() int { return 1 }
func (d releaseDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d releaseDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	release, ok := item.(*domain.Release)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderRelease(release, index == m.Index()))
}
