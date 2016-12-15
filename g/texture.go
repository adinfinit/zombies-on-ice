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

func (tex *Texture) Bounds() Rect {
	b := tex.RGBA.Bounds()
	return Rect{
		V2{float32(b.Min.X), float32(b.Min.Y)},
		V2{float32(b.Max.X), float32(b.Max.Y)},
	}
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
	sz := m.Bounds().Size()
	tex.Size.X = float32(sz.X)
	tex.Size.Y = float32(sz.Y)
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
