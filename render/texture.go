package render

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

// rendering

func (state *State) Texture(tex *g.Texture, dst g.Rect) {
	state.TextureTint(tex, dst, g.White)
}

func (state *State) TextureTint(tex *g.Texture, dst g.Rect, tint g.Color) {
	state.SpriteTint(tex, dst, tex.Bounds(), tint)
}

func (state *State) Sprite(tex *g.Texture, dst, src g.Rect) {
	state.SpriteTint(tex, dst, src, g.White)
}

func (state *State) SpriteTint(tex *g.Texture, dst, src g.Rect, tint g.Color) {
	uploaded := state.Textures.Upload(tex)
	uv := src.ScaleInv(tex.Size)

	uploaded.Bind()
	{
		gl.Color4ub(tint.RGBA())
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(uv.Min.X, uv.Max.Y)
			gl.Vertex2f(dst.Min.X, dst.Min.Y)

			gl.TexCoord2f(uv.Max.X, uv.Max.Y)
			gl.Vertex2f(dst.Max.X, dst.Min.Y)

			gl.TexCoord2f(uv.Max.X, uv.Min.Y)
			gl.Vertex2f(dst.Max.X, dst.Max.Y)

			gl.TexCoord2f(uv.Min.X, uv.Min.Y)
			gl.Vertex2f(dst.Min.X, dst.Max.Y)
		}
		gl.End()
	}
	uploaded.Unbind()
}

func (state *State) Line(tex *g.Texture, from, to g.V2, width float32) {
	state.LineTint(tex, from, to, width, g.White)
}

func (state *State) LineTint(tex *g.Texture, from, to g.V2, width float32, tint g.Color) {
	length := to.Sub(from).Length()
	normal := to.Sub(from).Rotate90().Normalize().Scale(width / 2)

	state.TextureQuadTint(tex, [4]g.V2{
		from.Sub(normal),
		from.Add(normal),
		to.Add(normal),
		to.Sub(normal),
	}, [4]g.V2{
		{0, length * tex.Size.Y},
		{1, length * tex.Size.Y},
		{1, 0},
		{0, 0},
	}, tint)
}

func (state *State) TextureQuadTint(tex *g.Texture, p [4]g.V2, uv [4]g.V2, tint g.Color) {
	uploaded := state.Textures.Upload(tex)

	uploaded.Bind()
	{
		gl.Color4ub(tint.RGBA())
		gl.Begin(gl.QUADS)
		{
			for i := range p {
				gl.TexCoord2f(uv[i].X, uv[i].Y)
				gl.Vertex2f(p[i].X, p[i].Y)
			}
		}
		gl.End()
	}
	uploaded.Unbind()
}
