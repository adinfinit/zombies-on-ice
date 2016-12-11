package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

type Powerup struct {
	Entity

	Life     float32
	Rotation float32

	Health       float32
	HammerRadius float32
}

func NewPowerup(bounds g.Rect) *Powerup {
	powerup := &Powerup{}

	powerup.Position = g.RandomV2(bounds.ScaleInv(g.V2{2, 2}))
	powerup.Radius = 0.5
	powerup.Mass = 1
	powerup.CollisionLayer = PowerupLayer
	powerup.CollisionMask = PlayerLayer
	powerup.CollisionTrigger = true

	powerup.Life = 10.0

	powerup.Health = 0.5
	powerup.HammerRadius = 0.1

	return powerup
}

func (powerup *Powerup) Entites() []*Entity {
	return []*Entity{&powerup.Entity}
}

func (powerup *Powerup) Update(dt float32) {
	powerup.Life -= dt
	powerup.Rotation += dt
}

func (powerup *Powerup) Done() bool {
	return powerup.Life < 0.0 || len(powerup.Collision) > 0
}

func (powerup *Powerup) Apply(player *Player) {
	player.Health = g.Clamp01(player.Health + powerup.Health)
	player.Hammer.Radius += powerup.HammerRadius
}

func (powerup *Powerup) Render(game *Game) {
	if powerup.Life < 3.0 {
		if g.Mod(powerup.Life*1.5, 1) < 0.4 {
			return
		}
	}

	gl.PushMatrix()
	{
		gl.Translatef(powerup.Position.X, powerup.Position.Y, 0)
		gl.Rotatef(g.RadToDeg(powerup.Rotation), 0, 0, -1)

		tex := game.Assets.Texture("assets/healthpack.png")
		tex.Draw(g.NewCircleRect(powerup.Radius))
	}
	gl.PopMatrix()
}
