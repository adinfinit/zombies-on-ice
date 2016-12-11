package main

import (
	"fmt"

	"github.com/loov/zombies-on-ice/g"
)

type Spawner struct {
	Wave               int
	DisplayWaveMessage float32

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

func (spawner *Spawner) Update(game *Game, dt float32) {
	spawner.startNextWave(game, dt)

	spawner.PowerupCooldown -= dt
	if spawner.PowerupCooldown < 0.0 {
		spawner.PowerupCooldown = 15.0 + g.RandomBetween(0.0, 5.0)
		if len(game.Powerups) < 3 {
			powerup := NewPowerup(game.Room.Bounds)
			game.Powerups = append(game.Powerups, powerup)
		}
	}

	spawner.DisplayWaveMessage -= dt
}

func (spawner *Spawner) Render(game *Game) {
	if spawner.DisplayWaveMessage < 0 {
		return
	}

	text := fmt.Sprintf("Wave %v", spawner.Wave)
	height := float32(3.0)
	width := game.Font.Width(text, height)
	game.Font.Draw(text, g.V2{-width / 2, height / 2}, height)
}
