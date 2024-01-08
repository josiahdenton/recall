package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/josiahdenton/recall/internal"
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
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(">>>>>>>>>>>>>>>>> STARTING LOGGER!")
	p := tea.NewProgram(internal.New())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
