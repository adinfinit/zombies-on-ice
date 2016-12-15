package main

import "github.com/loov/zombies-on-ice/g"

type Room struct {
	Bounds g.Rect

	TextureScale float32
}

func NewRoom() *Room {
	room := &Room{}

	//1600, 1000
	room.Bounds.Min = g.V2{-16, -10}
	room.Bounds.Max = g.V2{16, 10}
	room.TextureScale = 0.5

	return room
}

func (room *Room) Render(game *Game) {
	ground := game.Assets.Texture("assets/room.png")
	game.Renderer.Texture(ground, room.Bounds)
}
