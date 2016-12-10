package g

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"time"

	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	path    string
	modtime time.Time
	lasterr error

	RGBA *image.RGBA
	Size V2
	ID   uint32
}

func (tex *Texture) check(err error) bool {
	if err != nil {
		if err != tex.lasterr {
			log.Println(err)
		}
		tex.lasterr = err
		return true
	}
	return false
}

func (tex *Texture) Reload() {
	stat, err := os.Stat(filepath.FromSlash(tex.path))
	if tex.check(err) {
		return
	}

	modtime := stat.ModTime()
	if modtime.Equal(tex.modtime) {
		return
	}
	tex.modtime = modtime

	m, err := loadImage(filepath.FromSlash(tex.path))
	if tex.check(err) {
		return
	}

	tex.lasterr = nil
	tex.RGBA = m
}

func (tex *Texture) Upload() {
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture.ID)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(texture.RGBA.Rect.Size().X),
		int32(texture.RGBA.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(texture.RGBA.Pix))

	gl.Disable(gl.TEXTURE_2D)

	tex.check(glerror())
}

func (tex *Texture) Bind() {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)

	tex.check(glerror())
}

func (tex *Texture) Unbind() {
	gl.Disable(gl.TEXTURE_2D)
	tex.check(glerror())
}

func loadImage(filepath string) (*image.RGBA, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(m.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), m, image.Point{0, 0}, draw.Src)

	return rgba, nil
}
