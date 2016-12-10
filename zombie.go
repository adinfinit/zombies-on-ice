package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombieroom/g"
)

type Zombie struct {
	Entity
}

func NewZombie(bounds g.Rect) *Zombie {
	zombie := &Zombie{}

	zombie.Position = g.RandomV2(bounds)

	zombie.Elasticity = 0.2
	zombie.Mass = 1.0
	zombie.Dampening = 0.999
	zombie.Radius = 0.5

	zombie.CollisionLayer = ZombieLayer
	zombie.CollisionMask = HammerLayer

	return zombie
}

func (zombie *Zombie) Update(dt float32) {
}

func (zombie *Zombie) Respawn(bounds g.Rect) {
	if len(zombie.Collision) == 0 {
		return
	}

	zombie.Position = g.RandomV2(bounds)
}

func (zombie *Zombie) Render(game *Game) {
	rotation := float32(0)
	gl.PushMatrix()
	{
		gl.Translatef(zombie.Position.X, zombie.Position.Y, 0)

		gl.Rotatef(g.RadToDeg(rotation), 0, 0, -1)

		tex := game.Assets.TextureRepeat("assets/zombie.png")
		tex.Draw(g.NewCircleRect(zombie.Radius))
	}
	gl.PopMatrix()
}
