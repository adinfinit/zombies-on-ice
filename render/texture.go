package render

import (
	"image"

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
	// todo
}

// handling

type Textures struct {
	Uploaded map[*g.Texture]*Texture
}

func NewTextures() *Textures {
	return &Textures{
		Uploaded: make(map[*g.Texture]*Texture),
	}
}

func (textures *Textures) Upload(tex *g.Texture) *Texture {
	uploaded, ok := textures.Uploaded[tex]
	if ok {
		uploaded.Update(tex)
		return uploaded
	}

	uploaded = NewTexture(tex)
	textures.Uploaded[tex] = uploaded

	return uploaded
}

func (textures *Textures) BeginFrame() {
	for _, gltex := range textures.Uploaded {
		gltex.UseCount = 0
	}
}

func (textures *Textures) EndFrame() {
	for _, gltex := range textures.Uploaded {
		if gltex.UseCount == 0 {
			textures.Delete(gltex)
		}
	}
}

func (textures *Textures) Delete(gltex *Texture) {
	delete(textures.Uploaded, gltex.Texture)
	gltex.Delete()
}

type Texture struct {
	UseCount int

	Texture *g.Texture
	RGBA    *image.RGBA
	ID      uint32
}

func NewTexture(tex *g.Texture) *Texture {
	gltex := &Texture{}
	gltex.Texture = tex
	gltex.RGBA = tex.RGBA
	gltex.Upload()

	return gltex
}

func (gltex *Texture) Upload() {
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &gltex.ID)
	gl.BindTexture(gl.TEXTURE_2D, gltex.ID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if gltex.Texture.Repeat {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(gltex.RGBA.Rect.Size().X),
		int32(gltex.RGBA.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(gltex.RGBA.Pix))

	gl.Disable(gl.TEXTURE_2D)
}

func (gltex *Texture) Update(tex *g.Texture) {
	if tex.RGBA != gltex.RGBA {
		gltex.Delete()
		gltex.Upload()
	}
}

func (gltex *Texture) Delete() {
	gl.DeleteTextures(1, &gltex.ID)
	gltex.ID = 0
}

func (gltex *Texture) Bind() {
	gltex.UseCount++

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, gltex.ID)
}

func (gltex *Texture) Unbind() {
	gl.Disable(gl.TEXTURE_2D)
}
