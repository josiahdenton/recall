package task

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

type statusDelegate struct{}

func (d statusDelegate) Height() int  { return 1 }
func (d statusDelegate) Spacing() int { return 0 }
func (d statusDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d statusDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	status, ok := item.(*domain.Status)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStatus(status, index == m.Index()))
}

type stepDelegate struct{}

func (d stepDelegate) Height() int  { return 1 }
func (d stepDelegate) Spacing() int { return 0 }
func (d stepDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d stepDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	step, ok := item.(*domain.Step)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderStep(step, index == m.Index()))
}
