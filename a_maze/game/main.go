package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	maze "github.com/ilmanzo/hackweek24/a_maze/game/internal"
)

// main purpose of this projects is to learn and explore the Go import rules and directory structure

func main() {
	p := tea.NewProgram(
		maze.NewMaze(20, 20), tea.WithAltScreen(),
	)
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
	// Run() returns an interface, need type assertion to get the original type
	maze, ok := m.(*maze.MazeModel)
	if !ok {
		fmt.Printf("Unexpected model type: %T\n", maze)
		os.Exit(1)
	}
	fmt.Printf("\n Good! You walked %d steps to get the ticket\n", maze.StepsDone)
}
