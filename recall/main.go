package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal/ui"
	"log"
	"os"
)

func main() {
	if err := Run(); err != nil {
		fmt.Printf("failed to brew tea: %v", err)
		os.Exit(1)
	}
}

func Run() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", home, "recall-log"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// on non debug mode, set out to null dest?
	log.SetOutput(f)
	log.Println("--------------- Recall! ---------------")

	path := fmt.Sprintf("%s/%s", home, ".recall")

	p := tea.NewProgram(ui.New(path))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
