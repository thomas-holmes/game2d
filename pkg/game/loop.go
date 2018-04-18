package game

import (
	"log"
	"math"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Renderer *sdl.Renderer
	Window   *sdl.Window
}

// NewGame returns an entire game object. Only one should exist.
func NewGame() *Game {
	return &Game{}
}

// Init initializes the display and locks the goroutine to the executing thread.
func (g *Game) Init() error {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	window, err := sdl.CreateWindow("Game2D", 0, 0, 1920, 1080, sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	renderer, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return err
	}

	g.Window = window
	g.Renderer = renderer

	return nil
}

type Player struct {
	X, Y             float64
	XVelo, YVelo     float64
	XTarget, YTarget float64
	Color            sdl.Color
}

func NewPlayer() *Player {
	return &Player{X: 0, Y: 0, XVelo: 0, YVelo: 0, Color: sdl.Color{R: 0, G: 200, B: 60, A: 255}}
}

func (p *Player) HandleInput(input sdl.Event) {
	switch e := input.(type) {
	case *sdl.KeyboardEvent:
		if e.State != sdl.PRESSED {
			break
		}
		switch e.Keysym.Sym {
		case sdl.K_w:
			p.YVelo += -1
		case sdl.K_a:
			p.XVelo += -1
		case sdl.K_s:
			p.YVelo += 1
		case sdl.K_d:
			p.XVelo += 1
		}
	case *sdl.MouseMotionEvent:
		x, y := float64(e.X), float64(e.Y)
		p.XTarget = x
		p.YTarget = y
	}
}

func (p *Player) Update() {
	p.XVelo += (p.XTarget - p.X) / 100
	p.YVelo += (p.YTarget - p.Y) / 100

	if p.XVelo < 0 {
		p.XVelo = math.Max(-25, p.XVelo)
	} else {
		p.XVelo = math.Min(25, p.XVelo)
	}

	if p.YVelo < 0 {
		p.YVelo = math.Max(-25, p.YVelo)
	} else {
		p.YVelo = math.Min(25, p.YVelo)
	}

	p.X += p.XVelo
	p.Y += p.YVelo
}

func (game *Game) render(p *Player, delta float64) {
	vX, vY := p.XVelo+p.XVelo*delta, p.YVelo+p.YVelo*delta
	// vX, vY := float64(0), float64(0)
	x, y := p.X+vX, p.Y+vY

	rect := sdl.Rect{X: int32(x), Y: int32(y), W: 50, H: 50}
	r, g, b, a := uint8(p.Color.R), uint8(p.Color.G), uint8(p.Color.B), uint8(p.Color.A)
	game.Renderer.SetDrawColor(r, g, b, a)
	game.Renderer.FillRect(&rect)
}

const TimeStep int64 = 31250000 // 64 FPS

// Run blocks while running the game.
func (g *Game) Run() {
	log.Println("Hello from the game loop")

	var quit bool

	player := NewPlayer()

	start := time.Now()
	ticks := int64(0)
	lag := int64(0)

	for !quit {
		nowTicks := time.Now().Sub(start).Nanoseconds()
		lag += (nowTicks - ticks)
		ticks = nowTicks

		fps()
		for {
			e := sdl.PollEvent()
			if e == nil {
				break
			}
			if g.checkForQuit(e) {
				quit = true
			}
			player.HandleInput(e)
		}

		for lag > TimeStep {
			player.Update()
			lag -= TimeStep
		}

		if err := g.Renderer.SetDrawColorArray(200, 140, 200); err != nil {
			panic(err)
		}

		if err := g.Renderer.Clear(); err != nil {
			panic(err)
		}

		delta := float64(lag) / float64(TimeStep)
		g.render(player, delta)

		g.Renderer.Present()
	}
}

func (g *Game) checkForQuit(input sdl.Event) bool {
	switch e := input.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if e.Keysym.Sym == sdl.K_ESCAPE {
			return true
		}
	default:
		return false
	}
	return false
}

var frames uint
var last time.Time

func fps() {
	frames++
	now := time.Now()
	if now.Sub(last).Nanoseconds() > 1e9 {
		log.Println("FPS:", frames)
		frames = 0
		last = now
	}
}
