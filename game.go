package main

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/loov/zombies-on-ice/g"
)

const (
	ZeroLayer = CollisionLayer(1 << iota)
	PlayerLayer
	HammerLayer
	ZombieLayer
	PowerupLayer
)

type Game struct {
	Assets      *Assets
	Font        *g.Font
	Controllers *Controllers

	Spawner *Spawner

	Room      *Room
	Players   []*Player
	Zombies   []*Zombie
	Powerups  []*Powerup
	Particles *Particles

	CameraShake float32

	Clock float64

	lastPlayerID int
}

func NewGame() *Game {
	game := &Game{}

	game.Assets = NewAssets()
	game.Font = game.Assets.SpriteFont(
		"assets/arcade_43x74.png",
		g.V2{43, 74},
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.,;:?!-_~#\"'&()[]|`\\/@°+=*$£€<>%")

	game.Controllers = NewControllers()
	game.Room = NewRoom()
	game.Particles = NewParticles()

	game.Spawner = NewSpawner()
	/*
		for i := 0; i < 20; i++ {
			game.Zombies = append(game.Zombies, NewZombie(game.Room.Bounds))
		}
		for i := 0; i < 10; i++ {
			game.Powerups = append(game.Powerups, NewPowerup(game.Room.Bounds))
		}
	*/

	return game
}

func (game *Game) Update(window *glfw.Window, now float64) {
	game.Assets.Reload()

	dt := float32(now - game.Clock)
	game.Clock = now

	{ // update inputs
		game.Controllers.Update(window)

		active := []*Player{}
		for _, player := range game.Players {
			if !game.Controllers.Removed[&player.Controller] {
				active = append(active, player)
			}
		}

		for _, plugged := range game.Controllers.Plugged {
			if plugged.Controller == nil {
				game.lastPlayerID++
				player := NewPlayer(game.lastPlayerID)
				plugged.Controller = &player.Controller
				active = append(active, player)
			}
		}

		game.Players = active
	}

	game.Spawner.Update(game, dt)

	{
		// list all entities
		entities := []*Entity{}
		for _, player := range game.Players {
			entities = append(entities, player.Entities()...)
		}
		for _, zombie := range game.Zombies {
			entities = append(entities, zombie.Entities()...)
		}
		for _, powerup := range game.Powerups {
			entities = append(entities, powerup.Entities()...)
		}

		playerBySurvivor := make(map[*Entity]*Player)
		playerByHammer := make(map[*Entity]*Player)
		for _, player := range game.Players {
			playerBySurvivor[&player.Survivor] = player
			playerByHammer[&player.Hammer.Entity] = player
		}

		// reset entities
		for _, entity := range entities {
			entity.ResetForces()
		}

		// update survivors and hammers
		for _, player := range game.Players {
			player.Update(dt)
		}

		// update zombies
		for _, zombie := range game.Zombies {
			zombie.Update(game, dt)
		}

		// update powerups
		for _, powerup := range game.Powerups {
			powerup.Update(dt)
		}

		// update collision info
		HandleCollisions(entities, dt)

		// update camera shake
		amount := float32(0.0)
		for _, zombie := range game.Zombies {
			if strength, dead := zombie.DeathStrength(); dead {
				amount += strength
				for _, collision := range zombie.Collision {
					game.Particles.Spawn(32,
						collision.A.Position, collision.B.Velocity, 0.1, 0.4)

					playerByHammer[collision.B].Points += 1
				}
			}
		}
		game.CameraShake += amount * 0.05
		game.CameraShake -= dt
		game.CameraShake *= g.Pow(0.1, dt)
		if game.CameraShake < 0 {
			game.CameraShake = 0
		}
		if game.CameraShake > 0.5 {
			game.CameraShake = 0.5
		}

		// integrate forces
		for _, entity := range entities {
			entity.IntegrateForces(dt)
		}

		// update particles
		game.Particles.Update(dt)

		// apply constraints
		for _, player := range game.Players {
			player.ApplyConstraints(game.Room.Bounds)
		}

		// count points
		for _, player := range game.Players {
			for _, collision := range player.Survivor.Collision {
				if collision.B == &player.Hammer.Entity {
					continue
				}
				if collision.B.CollisionLayer == ZombieLayer {
					player.Health -= 0.01
				}
				if collision.B.CollisionLayer == HammerLayer {
					player.Health -= g.Clamp01(collision.B.Velocity.Length()) * 0.1
				}
			}
		}

		// update powerups
		powerups := []*Powerup{}
		for _, powerup := range game.Powerups {
			if powerup.Done() {
				if len(powerup.Collision) > 0 {
					player := playerBySurvivor[powerup.Collision[0].B]
					powerup.Apply(player)
				}
			} else {
				powerups = append(powerups, powerup)
			}
		}
		game.Powerups = powerups

		// respawn dead players
		for _, player := range game.Players {
			if player.Dead() {
				game.Particles.Spawn(64, player.Survivor.Position, g.V2{5, 0}, 0.2, g.Tau)
				player.Respawn()
			}
		}

		// remove dead zombies
		zombies := []*Zombie{}
		for _, zombie := range game.Zombies {
			if _, dead := zombie.DeathStrength(); !dead {
				zombies = append(zombies, zombie)
			}
		}
		game.Zombies = zombies

		// kill particles
		game.Particles.Kill(game.Room.Bounds)
	}
}

func (game *Game) Render(window *glfw.Window) {
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
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

	gl.Translatef(g.RandomV2Circle(game.CameraShake).XYZ())

	gl.Enable(gl.MULTISAMPLE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.ALPHA_TEST)

	// actual rendering
	game.Room.Render(game)

	game.Particles.RenderDecals(game)

	for _, zombie := range game.Zombies {
		zombie.Render(game)
	}

	for _, player := range game.Players {
		player.Render(game)
	}

	for _, powerup := range game.Powerups {
		powerup.Render(game)
	}

	game.Particles.Render(game)

	zero := g.V2{
		-screenSize.X/2 + 1,
		screenSize.Y/2 - 2,
	}

	for _, player := range game.Players {
		game.Font.DrawColored(fmt.Sprintf("%v", player.Points), zero, 1, player.Color)
		zero.Y -= 0.7
	}

	game.Spawner.Render(game)
}

func (game *Game) Unload() {
	game.Assets.Unload()
}
