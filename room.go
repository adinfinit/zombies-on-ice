package main

import "github.com/loov/zombieroom/g"

type Room struct {
	Bounds g.Rect

	TextureScale float32
}

func NewRoom() *Room {
	room := &Room{}

	room.Bounds.Min = g.V2{-14, -8}
	room.Bounds.Max = g.V2{14, 8}
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
