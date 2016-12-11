package main

import "github.com/loov/zombies-on-ice/g"

type Assets struct {
	Textures map[string]*g.Texture
}

func NewAssets() *Assets {
	return &Assets{
		Textures: make(map[string]*g.Texture),
	}
}

func (assets *Assets) Reload() {
	for _, tex := range assets.Textures {
		tex.Reload()
	}
}

func (assets *Assets) Texture(path string) *g.Texture { return assets.texture(path, false) }

func (assets *Assets) SpriteFont(path string, glyphSize g.V2, glyphs string) *g.Font {
	tex := assets.Texture(path)

	return &g.Font{
		Texture:   tex,
		Glyphs:    glyphs,
		GlyphSize: glyphSize,
	}
}
func (assets *Assets) TextureRepeat(path string) *g.Texture { return assets.texture(path, true) }

func (assets *Assets) texture(path string, repeat bool) *g.Texture {
	npath := path
	if repeat {
		npath = "@" + path
	}

	tex, ok := assets.Textures[npath]
	if !ok {
		tex = &g.Texture{}
		tex.Path = path
		tex.Repeat = repeat
		tex.Reload()

		assets.Textures[npath] = tex
	}

	return tex
}

func (assets *Assets) Unload() {
	for _, tex := range assets.Textures {
		tex.Delete()
	}
	*assets = *NewAssets()
}
