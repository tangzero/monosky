package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tangzero/monosky/monosky"
)

func main() {
	p := tea.NewProgram(monosky.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ouch!!! %v", err)
		os.Exit(1)
	}
}
