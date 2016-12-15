package main

import (
	"fmt"

	"github.com/loov/zombies-on-ice/g"
)

type Spawner struct {
	NoPlayers bool

	DeathVisuals bool

	Wave                int
	DisplayWaveMessage  float32
	DisplayDeathMessage float32

	PowerupCooldown float32
}

func NewSpawner() *Spawner {
	return &Spawner{}
}

func (spawner *Spawner) startNextWave(game *Game, dt float32) {
	if spawner.Wave < 3 {
		if len(game.Zombies) >= 5 {
			return
		}
	} else {
		if len(game.Zombies) >= 15 {
			return
		}
	}

	spawner.DisplayWaveMessage = 3.0
	spawner.Wave += 1

	zombieCount := spawner.Wave * 5
	if spawner.Wave > 3 {
		zombieCount += spawner.Wave * 10
	}
	if spawner.Wave > 6 {
		zombieCount += spawner.Wave * 10
	}

	for i := 0; i < zombieCount; i++ {
		game.Zombies = append(game.Zombies,
			NewZombie(game.Room.Bounds))
	}
}

func (spawner *Spawner) resetGame(game *Game) {
	spawner.PowerupCooldown = 0.0
	spawner.DisplayWaveMessage = 0.0
	spawner.DeathVisuals = false

	spawner.Wave = 0
	for _, player := range game.Players {
		player.Respawn()
	}

	game.Zombies = []*Zombie{}
}

func (spawner *Spawner) Update(game *Game, dt float32) {
	spawner.DisplayWaveMessage -= dt
	spawner.DisplayDeathMessage -= dt

	spawner.NoPlayers = len(game.Players) == 0
	if spawner.NoPlayers {
		return
	}

	anyPlayerAlive := false
	for _, player := range game.Players {
		if !player.Dead {
			anyPlayerAlive = true
			break
		}
	}
	if !anyPlayerAlive {
		if !spawner.DeathVisuals {
			spawner.DeathVisuals = true
			spawner.DisplayDeathMessage = 5.0
		}
		if spawner.DisplayDeathMessage < 0 {
			spawner.resetGame(game)
		}
		return
	}

	spawner.startNextWave(game, dt)

	spawner.PowerupCooldown -= dt
	if spawner.PowerupCooldown < 0.0 {
		spawner.PowerupCooldown = 15.0 + g.RandomBetween(0.0, 5.0)
		// spawner.PowerupCooldown = 3.0
		if len(game.Powerups) < 3 {
			powerup := NewPowerup(game.Room.Bounds)
			game.Powerups = append(game.Powerups, powerup)
		}
	}
}

func (spawner *Spawner) Render(game *Game) {
	if spawner.NoPlayers {
		game.Renderer.TextLines(
			game.Font,
			[]string{
				"Press WASD",
				"Press Arrows",
				"Press Start on Gamepad",
			}, g.V2{-12.0, 1.0}, 2.0, 1.5,
		)
		return
	}

	if spawner.DeathVisuals {
		text := fmt.Sprintf("Death at Wave %v", spawner.Wave)
		height := float32(3.0)
		width := game.Font.Width(text, height)
		pos := g.V2{-width / 2, height / 2}
		game.Renderer.Text(game.Font, text, pos, height)

		pos.Y -= height * 0.7

		height = 2.0
		for i, player := range game.Players {
			text := fmt.Sprintf("%d# %d points", i, int(player.Points))
			game.Renderer.TextTint(game.Font, text, pos, height, player.Color)

			pos.Y -= height * 0.7
		}
		return
	}

	if spawner.DisplayWaveMessage < 0 {
		return
	}

	text := fmt.Sprintf("Wave %v", spawner.Wave)
	height := float32(3.0)
	width := game.Font.Width(text, height)
	game.Renderer.Text(game.Font, text, g.V2{-width / 2, height / 2}, height)
}
