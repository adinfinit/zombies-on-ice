package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/loov/zombieroom/g"
)

type Player struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Hammer Hammer
}

type Hammer struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Tension float32
	Radius  float32
	Length  float32

	TensionMultiplier float32
	VelocityDampening float32
}

type Assets struct {
	Player g.Texture
	Rope   g.Texture
	Hammer g.Texture
}

type Game struct {
	Assets Assets
	Player Player

	Clock float64
}

func NewGame() *Game {
	game := &Game{}

	game.Assets.Hammer.Path = "assets/hammer.png"
	game.Assets.Player.Path = "assets/player.png"
	game.Assets.Rope.Path = "assets/rope.png"

	return game
}

func (game *Game) Update(window *glfw.Window, now float64) {
	// dt := float32(now - game.Clock)
	game.Clock = now

	game.Assets.Hammer.Reload()
	game.Assets.Player.Reload()
	game.Assets.Rope.Reload()

	// SCENE
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	ratio := float64(height) / float64(width)
	screenWidth := 30.0 // meters
	screenHeight := screenWidth * ratio

	gl.Ortho(-screenWidth/2, screenWidth/2, -screenHeight/2, screenHeight/2, 10, -10)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.ALPHA_TEST)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, game.Assets.Player.ID)

	{
		gl.Color4f(1, 1, 1, 1)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, 1)
			gl.Vertex2f(0, 0)
			gl.TexCoord2f(0, 0)
			gl.Vertex2f(0, 1)
			gl.TexCoord2f(1, 0)
			gl.Vertex2f(1, 1)
			gl.TexCoord2f(1, 1)
			gl.Vertex2f(1, 0)
		}
		gl.End()
	}

	gl.Disable(gl.TEXTURE_2D)
}

func (game *Game) Unload() {
	game.Assets.Hammer.Delete()
	game.Assets.Rope.Delete()
	game.Assets.Hammer.Delete()
}
