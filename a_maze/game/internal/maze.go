package maze

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps = 25
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
// 2 = treasure
// 3 = player
var valToString = [4]string{
	mazeStyle.Render("  "),
	mazeStyle.Render("\u2588\u2588"),
	treasureStyle.Render("🎂"),
	playerStyle.Render("🦔"),
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
	m.set(m.playerX, m.playerY, 3)
	m.treasureX = 3 + rand.Intn(w-4)
	m.treasureY = 3 + rand.Intn(h-4)
	m.set(m.treasureX, m.treasureY, 2)
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

func (m *maze) get(x, y int) byte {
	i := y*m.width + x
	if i > len(m.cells)-1 || x < 0 || y < 0 || x >= m.width || y >= m.height {
		return 1
	}
	return m.cells[i]
}

func (m *maze) checkForTreasure() (tea.Model, tea.Cmd) {
	if m.playerX == m.treasureX && m.playerY == m.treasureY {
		return m, tea.Quit
	}
	return m, nil
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
	return animate()
}

type frameMsg struct{}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(_ time.Time) tea.Msg {
		return frameMsg{}
	})
}

// this function is automatically called by framework
// with a message as parameter
func (m maze) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyUp {
			if m.get(m.playerX, m.playerY-1)%2 == 0 {
				m.set(m.playerX, m.playerY, 0)
				m.playerY -= 1
				m.set(m.playerX, m.playerY, 3)
				return m.checkForTreasure()
			}
		}
		if msg.Type == tea.KeyDown {
			if m.get(m.playerX, m.playerY+1)%2 == 0 {
				m.set(m.playerX, m.playerY, 0)
				m.playerY += 1
				m.set(m.playerX, m.playerY, 3)
				return m.checkForTreasure()
			}
		}
		if msg.Type == tea.KeyLeft {
			if m.get(m.playerX-1, m.playerY)%2 == 0 {
				m.set(m.playerX, m.playerY, 0)
				m.playerX -= 1
				m.set(m.playerX, m.playerY, 3)
				return m.checkForTreasure()
			}
		}
		if msg.Type == tea.KeyRight {
			if m.get(m.playerX+1, m.playerY)%2 == 0 {
				m.set(m.playerX, m.playerY, 0)
				m.playerX += 1
				m.set(m.playerX, m.playerY, 3)
				return m.checkForTreasure()
			}
		}
	case frameMsg:
		return m, animate()
	case tea.WindowSizeMsg:
		return NewMaze(msg.Width/2, msg.Height), nil
	default:
		return m, nil
	}
	return m, nil
}
