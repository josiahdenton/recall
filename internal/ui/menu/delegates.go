package menu

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type menuDelegate struct{}

func (d menuDelegate) Height() int  { return 1 }
func (d menuDelegate) Spacing() int { return 1 }
func (d menuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d menuDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	menuOption, ok := item.(*domain.MenuOption)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderMenuOption(menuOption, index == m.Index()))
}
