package g

import "strings"

type Font struct {
	Texture      *Texture
	Glyphs       string
	GlyphsPerRow int
	GlyphSize    V2
}

func (font *Font) Width(text string, glyphHeight float32) float32 {
	glyphWidth := glyphHeight * font.GlyphSize.X / font.GlyphSize.Y
	return glyphWidth * float32(len(text))
}

func (font *Font) Draw(text string, position V2, glyphHeight float32) {
	if font.GlyphsPerRow == 0 {
		font.GlyphsPerRow = int(font.Texture.Size.X / font.GlyphSize.X)
	}

	position = position.Add(V2{0.0, 0.0})
	glyphWidth := glyphHeight * font.GlyphSize.X / font.GlyphSize.Y

	zero := Rect{V2{}, font.GlyphSize}
	for _, r := range text {
		i := strings.IndexByte(font.Glyphs, byte(r))
		if i < 0 {
			position = position.Add(V2{glyphWidth, 0.0})
			continue
		}
		x, y := i%font.GlyphsPerRow, i/font.GlyphsPerRow

		glyph := zero.Offset(V2{
			X: float32(x) * font.GlyphSize.X,
			Y: float32(y) * font.GlyphSize.Y,
		})

		glyph = glyph.ScaleInv(font.Texture.Size)

		font.Texture.DrawSub(Rect{
			position,
			position.Add(V2{
				X: glyphWidth,
				Y: glyphHeight,
			}),
		}, glyph)

		position = position.Add(V2{glyphWidth, 0.0})
	}
}

func (font *Font) DrawColored(text string, position V2, glyphHeight float32, color Color) {
	if font.GlyphsPerRow == 0 {
		font.GlyphsPerRow = int(font.Texture.Size.X / font.GlyphSize.X)
	}

	position = position.Add(V2{0.0, 0.0})
	glyphWidth := glyphHeight * font.GlyphSize.X / font.GlyphSize.Y

	zero := Rect{V2{}, font.GlyphSize}
	for _, r := range text {
		i := strings.IndexByte(font.Glyphs, byte(r))
		if i < 0 {
			position = position.Add(V2{glyphWidth, 0.0})
			continue
		}
		x, y := i%font.GlyphsPerRow, i/font.GlyphsPerRow

		glyph := zero.Offset(V2{
			X: float32(x) * font.GlyphSize.X,
			Y: float32(y) * font.GlyphSize.Y,
		})

		glyph = glyph.ScaleInv(font.Texture.Size)

		font.Texture.DrawSubColored(Rect{
			position,
			position.Add(V2{
				X: glyphWidth,
				Y: glyphHeight,
			}),
		}, glyph, color)

		position = position.Add(V2{glyphWidth, 0.0})
	}
}
