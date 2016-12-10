package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/loov/zombieroom/g"
)

type Room struct {
	Bounds g.Rect

	TextureScale float32
}

func (room *Room) Render(game *Game) {
	ground := game.Assets.TextureRepeat("assets/ground.png")

	ground.DrawSub(
		room.Bounds,
		g.Rect{
			g.V2{0, 0},
			room.Bounds.Size().Scale(room.TextureScale),
		},
	)
}

type Game struct {
	Assets *Assets

	Room    Room
	Players []*Player

	Clock float64
}

func NewGame() *Game {
	game := &Game{}

	game.Assets = NewAssets()

	game.Players = append(game.Players, NewPlayer(&Keyboard_1))
	game.Players = append(game.Players, NewPlayer(&Keyboard_0))

	game.Room.Bounds.Min = g.V2{-14, -8}
	game.Room.Bounds.Max = g.V2{14, 8}
	game.Room.TextureScale = 0.5

	return game
}

func (game *Game) Update(window *glfw.Window, now float64) {
	dt := float32(now - game.Clock)
	game.Clock = now

	game.Assets.Reload()

	// SCENE
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	screenRatio := float32(height) / float32(width)
	roomSize := game.Room.Bounds.Size()

	var screenSize g.V2
	roomRatio := roomSize.Y / roomSize.X

	if screenRatio < roomRatio {
		screenSize.Y = roomSize.Y + 2
		screenSize.X = screenSize.Y / screenRatio
	} else {
		screenSize.X = roomSize.X + 2
		screenSize.Y = screenSize.X * screenRatio
	}

	gl.Ortho(
		float64(-screenSize.X/2),
		float64(screenSize.X/2),
		float64(-screenSize.Y/2),
		float64(screenSize.Y/2),
		10, -10)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.ALPHA_TEST)

	for _, player := range game.Players {
		player.Updater.Update(&player.Controller, window)
	}

	for _, player := range game.Players {
		player.Update(game, dt)
	}

	game.Room.Render(game)
	for _, player := range game.Players {
		player.Render(game)
	}

	RenderAxis()
}

func (game *Game) Unload() {
	game.Assets.Unload()
}
