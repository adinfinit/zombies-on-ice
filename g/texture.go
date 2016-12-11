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
	Path string

	modtime time.Time
	lasterr error

	Repeat bool
	RGBA   *image.RGBA
	Size   V2
	ID     uint32
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
	stat, err := os.Stat(filepath.FromSlash(tex.Path))
	if tex.check(err) {
		return
	}

	modtime := stat.ModTime()
	if modtime.Equal(tex.modtime) {
		return
	}
	tex.modtime = modtime

	m, err := loadImage(filepath.FromSlash(tex.Path))
	if tex.check(err) {
		return
	}

	tex.lasterr = nil
	tex.RGBA = m

	tex.Delete()
	tex.Upload()
}

func (tex *Texture) Draw(dst Rect) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	{
		gl.Color4f(1, 1, 1, 1)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, 1)
			gl.Vertex2f(dst.Min.X, dst.Min.Y)

			gl.TexCoord2f(1, 1)
			gl.Vertex2f(dst.Max.X, dst.Min.Y)

			gl.TexCoord2f(1, 0)
			gl.Vertex2f(dst.Max.X, dst.Max.Y)

			gl.TexCoord2f(0, 0)
			gl.Vertex2f(dst.Min.X, dst.Max.Y)
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}

func (tex *Texture) DrawColored(dst Rect, color Color) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	{
		gl.Color4f(color.Float())
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, 1)
			gl.Vertex2f(dst.Min.X, dst.Min.Y)

			gl.TexCoord2f(1, 1)
			gl.Vertex2f(dst.Max.X, dst.Min.Y)

			gl.TexCoord2f(1, 0)
			gl.Vertex2f(dst.Max.X, dst.Max.Y)

			gl.TexCoord2f(0, 0)
			gl.Vertex2f(dst.Min.X, dst.Max.Y)
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}

func (tex *Texture) DrawSub(dst Rect, src Rect) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	{
		gl.Color4f(1, 1, 1, 1)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(src.Min.X, src.Max.Y)
			gl.Vertex2f(dst.Min.X, dst.Min.Y)

			gl.TexCoord2f(src.Max.X, src.Max.Y)
			gl.Vertex2f(dst.Max.X, dst.Min.Y)

			gl.TexCoord2f(src.Max.X, src.Min.Y)
			gl.Vertex2f(dst.Max.X, dst.Max.Y)

			gl.TexCoord2f(src.Min.X, src.Min.Y)
			gl.Vertex2f(dst.Min.X, dst.Max.Y)
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}

func (tex *Texture) Line(from, to V2, width float32) {
	length := to.Sub(from).Length()
	normal := to.Sub(from).Rotate90().Normalize().Scale(width / 2)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	{
		gl.Color4f(1, 1, 1, 1)
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, length)
			gl.Vertex2f(from.Sub(normal).XY())

			gl.TexCoord2f(1, length)
			gl.Vertex2f(from.Add(normal).XY())

			gl.TexCoord2f(1, 0)
			gl.Vertex2f(to.Add(normal).XY())

			gl.TexCoord2f(0, 0)
			gl.Vertex2f(to.Sub(normal).XY())
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}

func (tex *Texture) LineColored(from, to V2, width float32, color Color) {
	length := to.Sub(from).Length()
	normal := to.Sub(from).Rotate90().Normalize().Scale(width / 2)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	{
		gl.Color4f(color.Float())
		gl.Begin(gl.QUADS)
		{
			gl.TexCoord2f(0, length)
			gl.Vertex2f(from.Sub(normal).XY())

			gl.TexCoord2f(1, length)
			gl.Vertex2f(from.Add(normal).XY())

			gl.TexCoord2f(1, 0)
			gl.Vertex2f(to.Add(normal).XY())

			gl.TexCoord2f(0, 0)
			gl.Vertex2f(to.Sub(normal).XY())
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}

func (tex *Texture) Upload() {
	log.Println("Upload texture", tex.Path)

	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &tex.ID)
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if tex.Repeat {
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
		int32(tex.RGBA.Rect.Size().X),
		int32(tex.RGBA.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(tex.RGBA.Pix))

	gl.Disable(gl.TEXTURE_2D)

	tex.check(glerror())
}

func (tex *Texture) Delete() {
	if tex.ID == 0 {
		return
	}
	gl.DeleteTextures(1, &tex.ID)
	tex.ID = 0
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

	if rgba, ok := m.(*image.RGBA); ok {
		return rgba, nil
	}

	rgba := image.NewRGBA(m.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), m, image.Point{0, 0}, draw.Src)

	return rgba, nil
}
