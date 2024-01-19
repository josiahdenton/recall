package forms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type createZettelOptionDelegate struct{}

func (d createZettelOptionDelegate) Height() int  { return 1 }
func (d createZettelOptionDelegate) Spacing() int { return 1 }
func (d createZettelOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d createZettelOptionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	option, ok := item.(*linkZettelOption)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderLinkOption(option, index == m.Index()))
}

type zettelDelegate struct{}

func (d zettelDelegate) Height() int  { return 1 }
func (d zettelDelegate) Spacing() int { return 0 }
func (d zettelDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d zettelDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	zettel, ok := item.(*domain.Zettel)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderZettel(zettel, index == m.Index()))
}

type createResourceOptionDelegate struct{}

func (d createResourceOptionDelegate) Height() int  { return 1 }
func (d createResourceOptionDelegate) Spacing() int { return 1 }
func (d createResourceOptionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d createResourceOptionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	option, ok := item.(*createResourceOption)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderResourceOption(option, index == m.Index()))
}

type resourceDelegate struct{}

func (d resourceDelegate) Height() int  { return 1 }
func (d resourceDelegate) Spacing() int { return 1 }
func (d resourceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d resourceDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	zettel, ok := item.(*domain.Resource)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderResource(zettel, index == m.Index()))
}
