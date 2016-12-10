package main

import "github.com/go-gl/gl/v2.1/gl"

func RenderAxis() {
	gl.Color4f(1, 0, 0, 1)
	gl.Begin(gl.QUADS)
	{
		gl.Vertex2f(0, 0)
		gl.Vertex2f(10, 0)
		gl.Vertex2f(10, 0.1)
		gl.Vertex2f(0, 0.1)
	}
	gl.End()

	gl.Color4f(0, 1, 0, 1)
	gl.Begin(gl.QUADS)
	{
		gl.Vertex2f(0, 0)
		gl.Vertex2f(0.1, 0)
		gl.Vertex2f(0.1, 10)
		gl.Vertex2f(0, 10)
	}
	gl.End()
}
