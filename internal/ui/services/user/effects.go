package user

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/services/toast"
	"log"
)

func New() *Effects {
	return &Effects{}
}

type Effects struct{}

func (e *Effects) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case copyMsg:
		cmd = e.copy(msg.content)
	}

	return cmd
}

func (e *Effects) copy(s string) tea.Cmd {
	// copy content into clipboard
	err := clipboard.WriteAll(s)
	if err != nil {
		log.Printf("failed to copy to clipboard: %v", err)
		return toast.ShowToast("failed to copy to clipboard", toast.Warn)
	}
	return toast.ShowToast("copied to clipboard!", toast.Info)
}
