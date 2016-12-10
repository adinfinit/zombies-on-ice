package main

import "github.com/loov/zombieroom/g"

type Collision struct {
	A, B   *Entity
	Normal g.V2
}

type CollisionLayer uint8

type Entity struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Mass       float32
	Radius     float32
	Elasticity float32
	Dampening  float32

	Collision []Collision

	CollisionLayer CollisionLayer
	CollisionMask  CollisionLayer
}

func (en *Entity) Entities() []*Entity { return []*Entity{en} }

func (en *Entity) ResetForces() {
	en.Force = g.V2{}
	if en.Mass == 0 {
		en.Mass = 1.0
	}

	en.Collision = nil
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

func HandleCollisions(entities []*Entity) {
	const SafeZone = 0.9

	for i, a := range entities {
		for _, b := range entities[i+1:] {
			if a.CollisionLayer&b.CollisionMask == 0 && b.CollisionLayer&a.CollisionMask == 0 {
				continue
			}

			dist := a.Position.Sub(b.Position)
			if dist.Length() < (a.Radius+b.Radius)*SafeZone {
				if a.CollisionMask&b.CollisionLayer != 0 {
					a.Collision = append(a.Collision, Collision{
						A:      a,
						B:      b,
						Normal: g.V2{},
					})
				}

				if b.CollisionMask&a.CollisionLayer != 0 {
					b.Collision = append(b.Collision, Collision{
						A:      b,
						B:      a,
						Normal: g.V2{},
					})
				}
			}
		}
	}
}
