package fireworks

import (
	"math"
	"math/rand"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/maaslalani/confetty/array"
	"github.com/maaslalani/confetty/simulation"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	framesPerSecond = 60.0
	numParticles    = 50
)

var (
	colors     = []string{"#fdff6a", "#ff718d"}
	characters = []string{"+", "*", "•"}
)

type frameMsg time.Time
type fireworkMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

type model struct {
	system *simulation.System
}

func Spawn(width, height int) []simulation.Particle {
	color := lipgloss.Color(array.Sample(colors))
	v := float64(rand.Intn(10) + 20.0)

	particles := []simulation.Particle{}

	x := rand.Float64() * float64(width)
	y := rand.Float64() * float64(height)

	for i := 0; i < numParticles; i++ {
		p := simulation.Particle{
			Physics: harmonica.NewProjectile(
				harmonica.FPS(framesPerSecond),
				harmonica.Point{X: x, Y: y},
				harmonica.Vector{X: math.Cos(float64(i)) * v, Y: math.Sin(float64(i)) * v / 2},
				harmonica.Vector(harmonica.TerminalGravity),
			),
			Char: lipgloss.NewStyle().Foreground(color).Render(array.Sample(characters)),
		}
		particles = append(particles, p)
	}
	return particles
}

func InitialModel() model {
	width, height, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	return model{system: &simulation.System{
		Particles: Spawn(width, height),
		Frame: simulation.Frame{
			Width:  width,
			Height: height,
		},
	}}
}

// Init initializes the confetti after a small delay
func (m model) Init() tea.Cmd {
	return animate()
}

// Update updates the model every frame, it handles the animation loop and
// updates the particle physics every frame
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		m.system.Particles = Spawn(m.system.Frame.Width, m.system.Frame.Height)
		return m, nil
	case frameMsg:
		m.system.Update()
		return m, animate()
	case tea.WindowSizeMsg:
		m.system.Frame.Width = msg.Width
		m.system.Frame.Height = msg.Height
		return m, nil
	default:
		return m, nil
	}
}

// View displays all the particles on the screen
func (m model) View() string {
	return m.system.Render()
}
