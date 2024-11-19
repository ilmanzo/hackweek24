package main

import (
	"fmt"
	"strings"

	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type maze struct {
	cells  []byte
	width  int
	height int
}

// use official SUSE colors for default background and foreground
var baseStyle = lipgloss.NewStyle().Background(lipgloss.Color("#192072")).Foreground(lipgloss.Color("#2453ff"))

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
	// draw grid on alternate cells
	for x := 2; x < w-2; x += 2 {
		for y := 1; y < h-1; y++ {
			m.set(x, y, 1)
		}
	}
	for y := 2; y < h-2; y += 2 {
		for x := 1; x < w-1; x++ {
			m.set(x, y, 1)
		}
	}
	// starting on the upper row, for each cell
	// flip a coin in order to decide which direction to carve
	for y := 1; y < h-1; y += 2 {
		for x := 1; x < w-1; x += 2 {
			if rand.Intn(2) == 1 {
				m.set(x+1, y, 0)
			} else {
				m.set(x, y+1, 0)
			}
		}
	}
	// carve some extra random spots (20%)
	for i := 0; i < (w*h)/5; i++ {
		m.set(2+rand.Intn(w-3), 2+rand.Intn(h-3), 0)
	}

	return m
}

// value 1 = turn on pixel
// value 0 = turn off pixel
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

// nothing to do on startup
func (m maze) Init() tea.Cmd {
	return nil
}

// this function is automatically called by framework
// with a message as parameter
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

// this must return a string rapresentation of our model
func (m maze) View() string {
	return baseStyle.Render(m.toString())
}

// utility func for debugging
func (m maze) print() {
	for i := 0; i < len(m.cells); i++ {
		fmt.Print(m.cells[i])
	}
}
