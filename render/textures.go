package render

import (
	"image"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/loov/zombies-on-ice/g"
)

type Textures struct {
	Uploaded    map[*g.Texture]*Texture
	MaxTextures int
}

func NewTextures() *Textures {
	return &Textures{
		Uploaded:    make(map[*g.Texture]*Texture),
		MaxTextures: 16,
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
	if len(textures.Uploaded) < textures.MaxTextures {
		return
	}

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
	log.Printf("Uploading %v", gltex.Texture.Path)

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
