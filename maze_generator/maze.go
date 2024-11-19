package main

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maze struct {
	cells  []byte
	width  int
	height int
}

var baseStyle = lipgloss.NewStyle().Background(lipgloss.Color("#173f4f")).Foreground(lipgloss.Color("#73ba25"))

// 0 = empty space
// 1 = only bottom filled
// 2 = only top filled
// 3 = both filled
var valToRune = [4]rune{' ', '\u2584', '\u2580', '\u2588'}

func NewMaze(w, h int) maze {
	cells := make([]byte, w*h)
	m := maze{cells: cells, width: w, height: h}
	// draw borders
	for x := 0; x < w; x++ {
		m.set(x, 0, 1)
		m.set(x, h-1, 1)
	}
	for y := 0; y < h; y++ {
		m.set(0, y, 1)
		m.set(w-1, y, 1)
	}

	// draw random dots for debug
	for i := 0; i < 100; i++ {
		m.set(rand.Intn(w), rand.Intn(h), 1)
	}
	// draw reference dots
	m.set(1, 1, 1)
	m.set(w-2, h-2, 1)
	m.set(w-2, 1, 1)
	m.set(1, h-2, 1)

	//	m.set(2, 1, 1)
	//	m.set(2, 2, 1)
	return m
}

func (m *maze) set(x, y int, value byte) {
	i := y*m.width + x
	if i > len(m.cells)-1 || x < 0 || y < 0 || x >= m.width || y >= m.height {
		return
	}
	m.cells[i] = value
}

// returns a string representing our model
func (m maze) toString() string {
	var sb strings.Builder
	y := 0
	for y < m.height {
		for x := 0; x < m.width; x++ {
			up := x + y*m.width
			down := up + m.width
			sb.WriteRune(valToRune[2*m.cells[up]+m.cells[down]])
		}
		y += 2
		if y < m.height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func (m maze) Init() tea.Cmd {
	return nil
}

func (m maze) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		return NewMaze(msg.Width, msg.Height*2), nil
	default:
		return m, nil
	}
}

func (m maze) View() string {
	content := m.toString()
	return baseStyle.Render(content)
}

func (m maze) print() {
	for i := 0; i < len(m.cells); i++ {
		fmt.Print(m.cells[i])
	}
}
