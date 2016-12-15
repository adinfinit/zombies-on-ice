package g

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
