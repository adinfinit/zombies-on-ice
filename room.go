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
	if false {
		ground := game.Assets.TextureRepeat("assets/ground.png")

		ground.DrawSub(
			room.Bounds,
			g.Rect{
				g.V2{0, 0},
				room.Bounds.Size().Scale(room.TextureScale),
			},
		)
	} else {
		ground := game.Assets.Texture("assets/room.png")
		ground.Draw(room.Bounds)
	}
}
