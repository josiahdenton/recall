package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/views"
	"os"
)

func main() {
	if err := Run(); err != nil {
		fmt.Printf("failed to brew tea: %v", err)
		os.Exit(1)
	}
}

func Run() error {
	p := tea.NewProgram(views.New())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
