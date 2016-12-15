package g

import "github.com/loov/zombies-on-ice/g/internal/breaklines"

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

func (font *Font) BreakLines(text string, glyphHeight, frameWidth float32) []string {
	glyphWidth := glyphHeight * font.GlyphSize.X / font.GlyphSize.Y
	return breaklines.SMAWK(text, frameWidth, func(text string) float32 {
		return float32(len(text)) * glyphWidth
	})
}

/*
type Typesetter struct {
	Font *Font

	Start V2
	Dot   V2

	GlyphHeight float32
	LineHeight float32
}

func (t *Typesetter) Advance(r rune)
*/
