package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombieroom/g"
)

type Particles struct {
	List []*Particle
}

type Particle struct {
	Position        g.V2
	Velocity        g.V2
	Rotation        float32
	AngularVelocity float32
	Radius          float32
	Life            float32
}

func NewParticles() *Particles { return &Particles{} }

func (ps *Particles) Spawn(amount int, position g.V2, velocity g.V2, radius float32, spread float32) {
	for i := 0; i < amount; i++ {
		rotate := g.RandomBetween(-spread/2, spread/2)
		speed := g.RandomBetween(-spread, spread)
		ps.List = append(ps.List, &Particle{
			Position:        position,
			Velocity:        velocity.Rotate(rotate).Scale(1 + speed),
			Rotation:        g.RandomBetween(0, 7),
			AngularVelocity: g.RandomBetween(-spread, spread),
			Radius:          g.RandomBetween(radius, radius*2),
			Life:            1,
		})
	}
}

func (ps *Particles) Update(dt float32) {
	for _, p := range ps.List {
		p.Position = p.Position.AddScale(p.Velocity, dt)
		p.Velocity = p.Velocity.Scale(g.Pow(0.9, dt))
		p.Life -= dt
		p.Radius -= dt * 0.2
		p.Rotation += p.AngularVelocity
	}
}

func (ps *Particles) Kill(bounds g.Rect) {
	list := ps.List[:0:cap(ps.List)]
	for _, p := range ps.List {
		if p.Life > 0.0 && p.Radius > 0.0 && bounds.Contains(p.Position) {
			list = append(list, p)
		}
	}
	ps.List = list
}

func (ps *Particles) Render(game *Game) {
	tex := game.Assets.TextureRepeat("assets/blood.png")
	for _, p := range ps.List {
		gl.PushMatrix()
		gl.Translatef(p.Position.X, p.Position.Y, 0)
		gl.Rotatef(g.RadToDeg(p.Rotation), 0, 0, -1)
		tex.Draw(g.NewCircleRect(p.Radius))
		gl.PopMatrix()
	}
}
