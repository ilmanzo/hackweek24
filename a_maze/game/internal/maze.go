package maze

import (
	"strings"

	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// this is our data model
// every maze cell is a string of two runes
type maze struct {
	cells     []byte
	width     int
	height    int
	playerX   int
	playerY   int
	treasureX int
	treasureY int
}

// use official SUSE colors for default background and foreground
// Midnight background, Waterhole foreground
// https://brand.suse.com/design-language#color
var commonBG = lipgloss.Color("#192072")
var mazeStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#2453ff"))
var playerStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#efefef"))
var treasureStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#fe7c3f"))

// 0 = empty space
// 1 = wall
// 2 = player
// 3 = treasure
var valToString = [4]string{
	mazeStyle.Render("  "),
	mazeStyle.Render("\u2588\u2588"),
	playerStyle.Render("🦔"),
	treasureStyle.Render("🎂"),
}

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
	// carve some extra random spots (25%)
	for i := 0; i < (w*h)/4; i++ {
		m.set(2+rand.Intn(w-3), 2+rand.Intn(h-3), 0)
	}
	//place player (in the center)
	//and treasure (random)
	m.playerX = w / 2
	m.playerY = h / 2
	m.set(m.playerX, m.playerY, 2)
	m.treasureX = 2 + rand.Intn(w-3)
	m.treasureY = 2 + rand.Intn(h-3)
	m.set(m.treasureX, m.treasureY, 3)
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
func (m maze) View() string {
	var sb strings.Builder
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			i := x + y*m.width
			sb.WriteString(valToString[m.cells[i]])
		}
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
		return NewMaze(msg.Width/2, msg.Height), nil
	default:
		return m, nil
	}
}
