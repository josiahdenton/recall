package forms

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui/toast"
)

func validateForm(errs ...error) tea.Cmd {
	for _, err := range errs {
		if err != nil {
			return toast.ShowToast(fmt.Sprintf("%v", err), toast.Warn)
		}
	}
	return nil
}
