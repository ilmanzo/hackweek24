package main

// purpose of this program is to test a "canvas" object
// where I can draw items, measure how fast is it
// and check how it behaves on different terminal size

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		newBall(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
}
