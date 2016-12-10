package main

import (
	"github.com/go-gl/gl/v2.1/gl"

	"github.com/loov/zombieroom/g"
)

type Player struct {
	Controller Controller
	Updater    ControllerUpdater

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

func NewPlayer(updater ControllerUpdater) *Player {
	player := &Player{}
	player.Updater = updater

	player.Hammer.Radius = 0.5
	player.Hammer.NormalLength = 1.5
	player.Hammer.MaxLength = 4
	player.Hammer.TensionMultiplier = 20

	return player
}

func (player *Player) AddForces(game *Game, dt float32) {
	hammer := &player.Hammer

	const force = 30
	const maxspeed = 20

	{ // player force
		in := player.Controller.Inputs[0]
		if in.Direction.Length() > 0.001 {
			player.Force = in.Direction.Scale(force)
			// todo scale lateral movement

			lateral := in.Direction.Rotate90c().Normalize()
			scale := lateral.Dot(player.Velocity)
			player.Force = player.Force.Add(lateral.Scale(-scale * 2.0))
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
}

func (player *Player) IntegrateForces(game *Game, dt float32) {
	hammer := &player.Hammer

	const maxspeed = 20

	{ // player physics
		player.Velocity = player.Velocity.AddScale(player.Force, dt)
		player.Velocity = g.ClampLength(player.Velocity, maxspeed)
		player.Position = player.Position.AddScale(player.Velocity, dt)

		g.EnforceInside(&player.Position, &player.Velocity, game.Room.Bounds, 0.2)

		dir := player.Position.Sub(hammer.Position)
		player.Direction = player.Direction.AddScale(dir, 100*dt).Normalize()
	}

	{ // hammer physics
		hammer.Velocity = hammer.Velocity.AddScale(hammer.Force, dt)
		hammer.Velocity = g.ClampLength(hammer.Velocity, maxspeed)
		hammer.Position = hammer.Position.AddScale(hammer.Velocity, dt)

		g.EnforceInside(&hammer.Position, &hammer.Velocity, game.Room.Bounds, 0.3)

		dir := hammer.Position.Sub(player.Position)
		hammer.Direction = hammer.Direction.AddScale(dir, 100*dt).Normalize()
	}
}

func (player *Player) Update(game *Game, dt float32) {
	player.AddForces(game, dt)
	player.IntegrateForces(game, dt)
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
