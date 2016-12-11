package main

import "github.com/loov/zombieroom/g"

type Collision struct {
	A, B          *Entity
	Normal        g.V2
	VelocityDelta g.V2
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

	CollisionTrigger bool
	CollisionLayer   CollisionLayer
	CollisionMask    CollisionLayer
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

func HandleCollisions(entities []*Entity, dt float32) {
	const SafeZone = 0.75

	for i, a := range entities {
		for _, b := range entities[i+1:] {
			delta := a.Position.Sub(b.Position)
			if delta.Length() >= (a.Radius+b.Radius)*SafeZone {
				continue
			}

			penetration := delta.Normalize().Scale((a.Radius+b.Radius)*SafeZone - delta.Length())

			// relative size
			ra := a.Mass / (a.Mass + b.Mass)
			rb := b.Mass / (a.Mass + b.Mass)

			if !a.CollisionTrigger {
				a.Position = a.Position.Add(penetration.Scale(rb))
			}
			if !b.CollisionTrigger {
				b.Position = b.Position.Sub(penetration.Scale(ra))
			}

			normal := delta.Normalize()
			p := 2.0 * (a.Velocity.Dot(normal) - b.Velocity.Dot(normal)) / (a.Mass + b.Mass)

			if !a.CollisionTrigger {
				a.AddForce(normal.Scale(-p * a.Mass / dt))
			}

			if !b.CollisionTrigger {
				b.AddForce(normal.Scale(p * b.Mass / dt))
			}

			velocityDelta := a.Velocity.Sub(b.Velocity)
			if a.CollisionMask&b.CollisionLayer != 0 {
				a.Collision = append(a.Collision, Collision{
					A:             a,
					B:             b,
					Normal:        normal,
					VelocityDelta: velocityDelta,
				})
			}

			if b.CollisionMask&a.CollisionLayer != 0 {
				b.Collision = append(b.Collision, Collision{
					A:             b,
					B:             a,
					Normal:        normal.Negate(),
					VelocityDelta: velocityDelta.Negate(),
				})
			}
		}
	}
}
