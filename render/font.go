package render

import (
	"strings"

	"github.com/loov/zombies-on-ice/g"
)

func (state *State) TextLines(font *g.Font, lines []string, position g.V2, glyphHeight, lineHeight float32) {
	for _, text := range lines {
		state.Text(font, text, position, glyphHeight)
		position.Y -= lineHeight
	}
}

func (state *State) Text(font *g.Font, text string, position g.V2, glyphHeight float32) {
	state.TextTint(font, text, position, glyphHeight, g.White)
}

func (state *State) TextTint(font *g.Font, text string, position g.V2, glyphHeight float32, tint g.Color) {
	// todo move layouting code into font

	if font.GlyphsPerRow == 0 {
		font.GlyphsPerRow = int(font.Texture.Size.X / font.GlyphSize.X)
	}

	position = position.Add(g.V2{0.0, 0.0})
	glyphWidth := glyphHeight * font.GlyphSize.X / font.GlyphSize.Y

	zero := g.Rect{g.V2{}, font.GlyphSize}
	for _, r := range text {
		i := strings.IndexByte(font.Glyphs, byte(r))
		if i < 0 {
			position = position.Add(g.V2{glyphWidth, 0.0})
			continue
		}
		x, y := i%font.GlyphsPerRow, i/font.GlyphsPerRow

		glyph := zero.Offset(g.V2{
			X: float32(x) * font.GlyphSize.X,
			Y: float32(y) * font.GlyphSize.Y,
		})

		state.SpriteTint(
			font.Texture,
			g.Rect{
				position,
				position.Add(g.V2{
					X: glyphWidth,
					Y: glyphHeight,
				}),
			},
			glyph,
			tint,
		)

		position = position.Add(g.V2{glyphWidth, 0.0})
	}
}
