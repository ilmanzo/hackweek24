package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// this is mainly an experiment on use of "half block" unicode chars
// to have double the vertical resolution

func main() {
	p := tea.NewProgram(
		NewMaze(20, 20), tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
}
