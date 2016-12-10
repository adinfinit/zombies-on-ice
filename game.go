package main

import (
	"github.com/goxjs/glfw"
	"github.com/loov/zombieroom/g"
)

type Player struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Hammer Hammer
}

type Hammer struct {
	Position g.V2
	Velocity g.V2
	Force    g.V2

	Tension float32
	Radius  float32
	Length  float32

	TensionMultiplier float32
	VelocityDampening float32
}

type Assets struct {
	PlayerTexture g.Texture
	RopeTexture   g.Texture
	HammerTexture g.Texture
}

type Game struct {
	Assets Assets
	Player Player
}

func NewGame() *Game {
	return &Game{}
}

func (game *Game) Update(window glfw.Window, dt float32) {
}

func (game *Game) Unload() {
	//todo
}
