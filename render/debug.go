package render

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

func (state *State) Gizmo() {
	gl.LineWidth(1)

	gl.Color4f(1, 0, 0, 1)
	gl.Begin(gl.LINES)
	{
		gl.Vertex2f(0, 0)
		gl.Vertex2f(10, 0)
	}
	gl.End()

	gl.Color4f(0, 1, 0, 1)
	gl.Begin(gl.LINES)
	{
		gl.Vertex2f(0, 0)
		gl.Vertex2f(0, 10)
	}
	gl.End()
}

func (state *State) Ray(position, direction g.V2, color g.Color) {
	gl.LineWidth(1)
	gl.Color4ub(color.RGBA())
	gl.Begin(gl.LINES)
	{
		gl.Vertex2f(position.XY())
		gl.Vertex2f(position.Add(direction).XY())
	}
	gl.End()
}

func glerror() error {
	if errcode := gl.GetError(); errcode != 0 {
		return fmt.Errorf("bind error: %d", errcode)
	}
	return nil
}
