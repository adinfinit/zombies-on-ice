package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

func RenderAxis() {
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

func RenderVector(at g.V2, dir g.V2, col g.Color) {
	gl.LineWidth(1)
	//gl.Color4f(0, 0, 1, 1)
	gl.Color4f(col.Float())
	gl.Begin(gl.LINES)
	{
		gl.Vertex2f(at.XY())
		gl.Vertex2f(at.Add(dir).XY())
	}
	gl.End()
}
