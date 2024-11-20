package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	maze "github.com/ilmanzo/hackweek24/a_maze/game/internal/maze"
)

// main purpose of this projects is to learn and explore the Go import rules and directory structure

func main() {
	p := tea.NewProgram(
		maze.NewMaze(20, 20), tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
}
