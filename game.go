package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/loov/zombieroom/g"
)

type Player struct {
	Controller Controller

	Position g.V2
	Velocity g.V2
	Force    g.V2

	Direction g.V2

	Hammer Hammer
}

type Hammer struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Direction g.V2

	Tension float32
	Radius  float32
	Length  float32

	NormalLength      float32
	MaxLength         float32
	TensionMultiplier float32
	VelocityDampening float32
}

func (player *Player) Update(dt float32) {
	hammer := &player.Hammer

	const force = 30
	const maxspeed = 20

	{ // player force
		in := player.Controller.Inputs[0]
		if in.Direction.Length() > 0.001 {
			player.Force = in.Direction.Scale(force)
			// todo scale lateral movement

			// dir := player.Force
			// player.Direction = player.Direction.AddScale(dir, dt)
			// player.Direction = player.Direction.Normalize()
		} else {
			player.Force = player.Velocity.Normalize().Negate().Scale(force)
		}
	}

	{ // hammer forces
		hammer.Force = g.V2{}

		player.Hammer.Tension = 0
		d := hammer.Position.Sub(player.Position)
		length := d.Length()

		delta := g.Max(length-hammer.NormalLength, 0)
		hammer.Tension = delta * hammer.TensionMultiplier

		f := g.V2{
			d.X * hammer.Tension / (length + 1),
			d.Y * hammer.Tension / (length + 1),
		}

		hammer.Force = hammer.Force.Sub(f)
		player.Force = player.Force.Add(f)

		dist := hammer.Position.Sub(player.Position)
		if n := dist.Length(); n > hammer.MaxLength {
			hammer.Position = player.Position.AddScale(dist, hammer.MaxLength/n)
		}
	}

	{ // player physics
		player.Velocity = player.Velocity.AddScale(player.Force, dt)
		player.Velocity = g.ClampLength(player.Velocity, maxspeed)
		player.Position = player.Position.AddScale(player.Velocity, dt)

		dir := player.Position.Sub(hammer.Position)
		player.Direction = player.Direction.AddScale(dir, 100*dt).Normalize()
	}

	{ // hammer physics
		hammer.Velocity = hammer.Velocity.AddScale(hammer.Force, dt)
		hammer.Velocity = g.ClampLength(hammer.Velocity, maxspeed)
		hammer.Position = hammer.Position.AddScale(hammer.Velocity, dt)

		dir := hammer.Position.Sub(player.Position)
		hammer.Direction = hammer.Direction.AddScale(dir, 100*dt).Normalize()
	}
}

func (player *Player) Render(game *Game) {
	hammer := &player.Hammer
	{
		rope := game.Assets.TextureRepeat("assets/rope.png")
		rope.Line(
			hammer.Position,
			player.Position,
			hammer.Radius/2)
	}

	gl.PushMatrix()
	{
		gl.Translatef(player.Position.X, player.Position.Y, 0)

		rotation := -(player.Direction.Angle() - g.Tau/4)
		gl.Rotatef(g.RadToDeg(rotation), 0, 0, -1)

		tex := game.Assets.TextureRepeat("assets/player.png")
		tex.Draw(g.NewRect(1, 1))
	}
	gl.PopMatrix()

	gl.PushMatrix()
	{
		gl.Translatef(hammer.Position.X, hammer.Position.Y, 0)

		rotation := -(hammer.Direction.Angle() - g.Tau/4)
		gl.Rotatef(g.RadToDeg(rotation), 0, 0, -1)

		tex := game.Assets.TextureRepeat("assets/hammer.png")
		tex.Draw(g.NewRect(hammer.Radius*2, hammer.Radius*2))
	}
	gl.PopMatrix()
}

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

	Player Player
	Room   Room

	Clock float64
}

func NewGame() *Game {
	game := &Game{}

	game.Assets = NewAssets()

	game.Player.Hammer.Radius = 0.5
	game.Player.Hammer.NormalLength = 1.5
	game.Player.Hammer.MaxLength = 4
	game.Player.Hammer.TensionMultiplier = 20

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

	Keyboard_1.Update(&game.Player.Controller, window)
	game.Player.Update(dt)

	game.Room.Render(game)
	game.Player.Render(game)

	RenderAxis()
}

func (game *Game) Unload() {
	game.Assets.Unload()
}
