package maze

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// use official SUSE colors for default background and foreground
// Midnight background, Waterhole foreground
// https://brand.suse.com/design-language#color
var commonBG = lipgloss.Color("#192072")
var mazeStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#2453ff"))
var playerStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#efefef"))
var treasureStyle = lipgloss.NewStyle().Background(commonBG).Foreground(lipgloss.Color("#fe7c3f"))

type cellContent byte

const (
	EmptyCell = iota
	WallCell
	TreasureCell
	PlayerCell
	DoorCell
	nDoors = 7
)

// 0 = empty space
// 1 = wall
// 2 = treasure
// 3 = player
// 4 = door
var valToString = [5]string{
	mazeStyle.Render("  "),
	mazeStyle.Render("\u2588\u2588"),
	mazeStyle.Render("🪪"),
	mazeStyle.Render("🐧"),
	mazeStyle.Render("🚪"),
}

// this is our data model
// every MazeModel cell is a string of two runes
type MazeModel struct {
	cells     []cellContent
	width     int
	height    int
	playerX   int
	playerY   int
	treasureX int
	treasureY int
	StepsDone int // exported step counter
	doorsX    [nDoors]int
	doorsY    [nDoors]int
}

func NewMaze(w, h int) MazeModel {
	cells := make([]cellContent, w*h)
	m := MazeModel{cells: cells, width: w, height: h, StepsDone: 0}
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
	//place player (+/- in the center)
	m.playerX = w / 2
	m.playerY = h / 2
	m.set(m.playerX, m.playerY, PlayerCell)
	//some doors (random)
	for i := 0; i < nDoors; i++ {
		m.doorsX[i] = 3 + rand.Intn(w-4)
		m.doorsY[i] = 3 + rand.Intn(h-4)
		m.set(m.doorsX[i], m.doorsY[i], DoorCell)
	}
	//treasure (random)
	m.treasureX = 3 + rand.Intn(w-4)
	m.treasureY = 3 + rand.Intn(h-4)
	m.set(m.treasureX, m.treasureY, TreasureCell)

	return m
}

func (m *MazeModel) set(x, y int, value cellContent) {
	i := y*m.width + x
	if i > len(m.cells)-1 || x < 0 || y < 0 || x >= m.width || y >= m.height {
		return
	}
	m.cells[i] = value
}

func (m *MazeModel) get(x, y int) cellContent {
	i := y*m.width + x
	if i > len(m.cells)-1 || x < 0 || y < 0 || x >= m.width || y >= m.height {
		return 1
	}
	return m.cells[i]
}

func (m *MazeModel) checkCollisions() (tea.Model, tea.Cmd) {
	m.StepsDone += 1
	for i := 0; i < nDoors; i++ {
		if m.playerX == m.doorsX[i] && m.playerY == m.doorsY[i] {
			m.set(m.playerX, m.playerY, DoorCell)
			m.playerX = m.width / 2
			m.playerY = m.height / 2
			m.set(m.playerX, m.playerY, PlayerCell)
		}
	}
	if m.playerX == m.treasureX && m.playerY == m.treasureY {
		return m, tea.Quit
	}
	return m, nil
}

// returns a string representing our model
func (m MazeModel) View() string {
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
func (m MazeModel) Init() tea.Cmd {
	return nil
}

// this function is automatically called by framework
// with a message as parameter
func (m MazeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return &m, tea.Quit
		}
		if msg.Type == tea.KeyUp && m.get(m.playerX, m.playerY-1)%2 == 0 {
			m.set(m.playerX, m.playerY, 0)
			m.playerY -= 1
			m.set(m.playerX, m.playerY, 3)
			return m.checkCollisions()
		}
		if msg.Type == tea.KeyDown && m.get(m.playerX, m.playerY+1)%2 == 0 {
			m.set(m.playerX, m.playerY, 0)
			m.playerY += 1
			m.set(m.playerX, m.playerY, 3)
			return m.checkCollisions()
		}
		if msg.Type == tea.KeyLeft && m.get(m.playerX-1, m.playerY)%2 == 0 {
			m.set(m.playerX, m.playerY, 0)
			m.playerX -= 1
			m.set(m.playerX, m.playerY, 3)
			return m.checkCollisions()
		}
		if msg.Type == tea.KeyRight && m.get(m.playerX+1, m.playerY)%2 == 0 {
			m.set(m.playerX, m.playerY, 0)
			m.playerX += 1
			m.set(m.playerX, m.playerY, 3)
			return m.checkCollisions()
		}
	case tea.WindowSizeMsg:
		// on resize, generate a new Maze
		// half width because every maze cell is 2 chars
		return NewMaze(msg.Width/2, msg.Height), nil
	}
	return m, nil
}
