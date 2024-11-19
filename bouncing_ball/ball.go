package main

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().Background(lipgloss.Color("#0c322c")).Foreground(lipgloss.Color("#30ba78")).Bold(true)

const (
	fps = 50
	//blockchar = "*"
	//blockchar = "\xE2\x9A\xAC"
	//blockchar = "\xF0\x9F\x9E\x84" // 🞄
	blockchar = "\xE2\xAD\x98"
	radius    = 4
)

type cellbuffer struct {
	cells  []string
	stride int
}

type model struct {
	cells                cellbuffer
	x, y                 float64
	xVelocity, yVelocity float64
	terminal_width       int
	terminal_height      int
	geekox               float64
}

func newBall() model {
	return model{
		x:         20,
		y:         10,
		xVelocity: 0.5,
		yVelocity: 0.5,
		geekox:    0,
	}
}

func writeString(cb *cellbuffer, x, y int, message string) {
	i := y*cb.stride + x
	if i > len(cb.cells)-1 || i < 0 {
		return
	}
	for _, ch := range strings.Split(message, "") {
		cb.cells[i] = ch
		i++
	}
}

func drawEllipse(cb *cellbuffer, xc, yc, rx, ry float64) {
	var (
		dx, dy, d1, d2 float64
		x              float64
		y              = ry
	)

	d1 = ry*ry - rx*rx*ry + 0.25*rx*rx
	dx = 2 * ry * ry * x
	dy = 2 * rx * rx * y

	for dx < dy {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d1 < 0 {
			x++
			dx = dx + (2 * ry * ry)
			d1 = d1 + dx + (ry * ry)
		} else {
			x++
			y--
			dx = dx + (2 * ry * ry)
			dy = dy - (2 * rx * rx)
			d1 = d1 + dx - dy + (ry * ry)
		}
	}

	d2 = ((ry * ry) * ((x + 0.5) * (x + 0.5))) + ((rx * rx) * ((y - 1) * (y - 1))) - (rx * rx * ry * ry)

	for y >= 0 {
		cb.set(int(x+xc), int(y+yc))
		cb.set(int(-x+xc), int(y+yc))
		cb.set(int(x+xc), int(-y+yc))
		cb.set(int(-x+xc), int(-y+yc))
		if d2 > 0 {
			y--
			dy = dy - (2 * rx * rx)
			d2 = d2 + (rx * rx) - dy
		} else {
			y--
			x++
			dx = dx + (2 * ry * ry)
			dy = dy - (2 * rx * rx)
			d2 = d2 + dx - dy + (rx * rx)
		}
	}
}

func (c *cellbuffer) init(w, h int) {
	if w == 0 {
		return
	}
	c.stride = w
	c.cells = make([]string, w*h)
	c.wipe()
}

func (c cellbuffer) set(x, y int) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = blockchar
}

func (c *cellbuffer) wipe() {
	for i := range c.cells {
		c.cells[i] = " "
	}
}

func (c cellbuffer) width() int {
	return c.stride
}

func (c cellbuffer) height() int {
	h := len(c.cells) / c.stride
	if len(c.cells)%c.stride != 0 {
		h++
	}
	return h
}

func (c cellbuffer) ready() bool {
	return len(c.cells) > 0
}

func (c cellbuffer) String() string {
	var b strings.Builder
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			b.WriteRune('\n')
		}
		b.WriteString(c.cells[i])
	}
	return b.String()
}

type frameMsg struct{}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(_ time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return animate()
}

// https://www.zq1.de/%7Ebernhard/images/share/geeko/logo.txt
func drawGeeko(cb *cellbuffer, x int) {
	writeString(cb, x, 10, "            ▁▃▄▅▆▇▇▇▇▇▇▇▆▅▅▄▃▂     ▕▇▆▄▂")
	writeString(cb, x, 11, "        ▂▄▆████████████████████▇▆▄▂▕█████▆▄▂")
	writeString(cb, x, 12, "     ▂▅█████████████████████████████████▛▀▔▔▀▖")
	writeString(cb, x, 13, "   ▗▆███████████████████████████████████▏  ▗▖▐▖")
	writeString(cb, x, 14, "  ▄████▀▀▀▔▔▀▀▀█████████████████████████▙▁  ▁▟▊")
	writeString(cb, x, 15, " ▐███▀   ▁▂▂▂   ▔▜████████████████████▙▃▀▀▀▀▀▀▀")
	writeString(cb, x, 16, " ███▍  ▗▇█▀▀▀█▆▖  ▀███████████████████████▇▇▀▀")
	writeString(cb, x, 17, " ▜██▎  ▐█▎    ▜█▖  ▐████▀▀▀▀▀▀████▊   ▔▔▔▔  ")
	writeString(cb, x, 18, " ▝██▙   ▀▀    ▐█▊   ▜███▃▁     ▔▜██▃▁")
	writeString(cb, x, 19, "  ▝███▄▂    ▂▅██▘    ▔▀▀▀▀▀      ▀▀▀▀▀")
	writeString(cb, x, 20, "   ▔▀█████████▛▘")
	writeString(cb, x, 21, "      ▔▀▀▀▀▀▀▔ ")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.cells.init(msg.Width, msg.Height)
		m.terminal_height = msg.Height
		m.terminal_width = msg.Width
		return m, nil
	case frameMsg:
		if !m.cells.ready() {
			return m, nil
		}
		newx := m.x + m.xVelocity
		newy := m.y + m.yVelocity
		if newx <= 2*radius || newx >= float64(m.terminal_width)-2*radius {
			m.xVelocity = -m.xVelocity
		}
		if newy >= float64(m.terminal_height)-radius {
			m.yVelocity = -m.yVelocity
			if rand.Intn(10) == 1 {
				m.yVelocity -= rand.Float64() / 2
			}

		}
		m.yVelocity += 0.03
		m.x += m.xVelocity
		m.y += m.yVelocity
		m.geekox += 0.03
		if int(m.geekox) > m.terminal_width-47 {
			m.geekox = 0
		}
		m.cells.wipe()
		drawEllipse(&m.cells, m.x, m.y, 2*radius, radius)
		writeString(&m.cells, int(m.x)-2, int(m.y)-1, "SUSE")
		writeString(&m.cells, int(m.x)-4, int(m.y), "HackWeek")
		writeString(&m.cells, int(m.x)-2, int(m.y)+1, "2024")
		drawGeeko(&m.cells, int(m.geekox))
		return m, animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	content := m.cells.String()
	return baseStyle.Render(content)
}
