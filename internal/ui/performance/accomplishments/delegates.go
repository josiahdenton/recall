package accomplishments

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/domain"
	"io"
)

type accomplishmentDelegate struct{}

func (d accomplishmentDelegate) Height() int  { return 1 }
func (d accomplishmentDelegate) Spacing() int { return 1 }
func (d accomplishmentDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d accomplishmentDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	accomplishment, ok := item.(*domain.Accomplishment)
	if !ok {
		return
	}
	fmt.Fprintf(w, renderAccomplishment(accomplishment, index == m.Index()))
}
