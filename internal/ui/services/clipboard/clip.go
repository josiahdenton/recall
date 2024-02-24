package clipboard

import (
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
)

func New() *Clip {
	return &Clip{}
}

type Clip struct{}

func (c *Clip) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case copyMsg:
		err := clipboard.WriteAll(msg.content)
		if err != nil {
			return toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn)
		} else {
			return toast.ShowToast("copied to clipboard", toast.Info)
		}
	}
	return nil
}

type copyMsg struct {
	content string
}

func Copy(content string) tea.Cmd {
	return func() tea.Msg {
		return copyMsg{content: content}
	}
}
