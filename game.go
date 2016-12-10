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
	Ground g.Texture
	Test   g.Texture
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

	game.Assets.Ground.Repeat = true
	game.Assets.Ground.Path = "assets/ground.png"

	game.Assets.Test.Repeat = true
	game.Assets.Test.Path = "assets/test.png"

	return game
}

func (game *Game) Update(window *glfw.Window, now float64) {
	// dt := float32(now - game.Clock)
	game.Clock = now

	game.Assets.Player.Reload()
	game.Assets.Hammer.Reload()
	game.Assets.Rope.Reload()

	game.Assets.Ground.Reload()
	game.Assets.Test.Reload()

	// SCENE
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	screenRatio := float32(height) / float32(width)
	roomSize := g.V2{40.0, 30.0}

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

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, game.Assets.Ground.ID)
	{
		gl.Color4f(1, 1, 1, 1)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, roomSize.Y/2)
			gl.Vertex2f(-roomSize.X/2, roomSize.Y/2)

			gl.TexCoord2f(0, 0)
			gl.Vertex2f(-roomSize.X/2, -roomSize.Y/2)

			gl.TexCoord2f(roomSize.X/2, 0)
			gl.Vertex2f(roomSize.X/2, -roomSize.Y/2)

			gl.TexCoord2f(roomSize.X/2, roomSize.Y/2)
			gl.Vertex2f(roomSize.X/2, roomSize.Y/2)
		}
		gl.End()
	}

	gl.Disable(gl.TEXTURE_2D)
}

func (game *Game) Unload() {
	game.Assets.Player.Delete()
	game.Assets.Hammer.Delete()
	game.Assets.Rope.Delete()
	game.Assets.Ground.Delete()
	game.Assets.Test.Delete()
}
