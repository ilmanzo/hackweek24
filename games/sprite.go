package main

// test for moving object with arrow keys

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	fps = 50
)

func main() {
	p := tea.NewProgram(
		NewGrid(20, 20), tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Whops, there's been an error: %v", err)
		os.Exit(1)
	}
}

type grid struct {
	cells  []byte
	width  int
	height int

	playerX int
	playerY int
}

type spritedata [4][4]byte

var player = spritedata{{0, 1, 0, 0}, {1, 1, 1, 0}, {0, 1, 0, 0}, {1, 0, 1, 0}}

// 0 = empty space
// 1 = only bottom filled
// 2 = only top filled
// 3 = both filled
var valToRune = [4]rune{' ', '\u2584', '\u2580', '\u2588'}

func NewGrid(w, h int) grid {
	cells := make([]byte, w*h)
	g := grid{cells: cells, width: w, height: h}
	// draw borders
	for x := 0; x < w; x++ {
		g.set(x, 0, 1)
		g.set(x, h-1, 1)
	}
	for y := 0; y < h; y++ {
		g.set(0, y, 1)
		g.set(w-1, y, 1)
	}
	g.playerX = w / 2
	g.playerY = h / 2
	return g
}

func (g *grid) getPosition(x, y int) int {
	return y*g.width + x
}

func (g *grid) putSprite(x, y int, s spritedata) {
	i := g.getPosition(x, y)
	if i > len(g.cells)-1 || x < 0 || y < 0 || x >= g.width || y >= g.height {
		return
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			g.cells[i] = s[y][x]
			i++
		}
		i += g.width - 4
	}
}

type frameMsg struct{}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(_ time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (g grid) Init() tea.Cmd {
	return animate()
}

// value 1 = turn on pixel
// value 0 = turn off pixel
func (g *grid) set(x, y int, value byte) {
	i := g.getPosition(x, y)
	if i > len(g.cells)-1 || x < 0 || y < 0 || x >= g.width || y >= g.height {
		return
	}
	g.cells[i] = value
}

// returns a string representing our model
func (g grid) View() string {
	var sb strings.Builder
	y := 0
	for y < g.height {
		for x := 0; x < g.width; x++ {
			up := x + y*g.width
			down := up + g.width
			sb.WriteRune(valToRune[2*g.cells[up]+g.cells[down]])
		}
		y += 2
		if y < g.height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

// this function is automatically called by framework
// with a message as parameter
// TODO move sprite with arrows
func (g grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return g, tea.Quit
	case tea.WindowSizeMsg:
		return NewGrid(msg.Width, msg.Height*2), nil
	case frameMsg:
		// TODO at each frame fill grid with zero

		g.putSprite(g.playerX, g.playerY, player)
		return g, animate()
	default:
		return g, nil
	}
}
