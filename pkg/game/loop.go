package game

import (
	"log"
	"math"
	"runtime"

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
	XAccel, YAccel   float64
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
	p.XAccel += (p.XTarget - p.X) / 100
	p.YAccel += (p.YTarget - p.Y) / 100

	if p.XAccel < 0 {
		p.XAccel = math.Max(-25, p.XAccel)
	} else {
		p.XAccel = math.Min(25, p.XAccel)
	}

	if p.YAccel < 0 {
		p.YAccel = math.Max(-25, p.YAccel)
	} else {
		p.YAccel = math.Min(25, p.YAccel)
	}

	p.X += p.XAccel
	p.Y += p.YAccel
	/*
		p.XVelo = 0
		p.YVelo = 0
	*/
}

func (game *Game) render(p *Player) {
	rect := sdl.Rect{X: int32(p.X), Y: int32(p.Y), W: 50, H: 50}
	r, g, b, a := uint8(p.Color.R), uint8(p.Color.G), uint8(p.Color.B), uint8(p.Color.A)
	game.Renderer.SetDrawColor(r, g, b, a)
	game.Renderer.FillRect(&rect)
}

// Run blocks while running the game.
func (g *Game) Run() {
	log.Println("Hello from the game loop")

	var quit bool

	player := NewPlayer()

	for !quit {
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
		player.Update()

		if err := g.Renderer.SetDrawColorArray(200, 140, 200); err != nil {
			panic(err)
		}

		if err := g.Renderer.Clear(); err != nil {
			panic(err)
		}

		g.render(player)

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
