package main

import "github.com/loov/zombieroom/g"

type Entity struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Mass       float32
	Radius     float32
	Elasticity float32
	Dampening  float32
}

func (en *Entity) Entities() []*Entity { return []*Entity{en} }

func (en *Entity) ResetForces() {
	en.Force = g.V2{}
	if en.Mass == 0 {
		en.Mass = 1.0
	}
}

func (en *Entity) AddForce(force g.V2) {
	en.Force = en.Force.Add(force)
}

func (en *Entity) IntegrateForces(dt float32) {
	const MaxEntitySpeed = 20

	en.Velocity = en.Velocity.AddScale(en.Force, dt/en.Mass)
	en.Velocity = g.ClampLength(en.Velocity, MaxEntitySpeed)
	en.Velocity = en.Velocity.Scale(en.Dampening)
	en.Position = en.Position.AddScale(en.Velocity, dt)
}

func (en *Entity) ConstrainInside(rect g.Rect) {
	g.EnforceInside(&en.Position, &en.Velocity, rect, en.Elasticity)
}
