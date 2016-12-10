package g

type Color struct{ R, G, B, A uint8 }

var (
	Red   = Color{0xFF, 0x00, 0x00, 0xFF}
	Green = Color{0x00, 0xFF, 0x00, 0xFF}
	Blue  = Color{0x00, 0x00, 0xFF, 0xFF}
)

func ColorHex(hex uint32) Color {
	return Color{
		R: uint8(hex >> 24),
		G: uint8(hex >> 16),
		B: uint8(hex >> 8),
		A: uint8(hex >> 0),
	}
}

func (c Color) Float() (r, g, b, a float32) {
	return float32(c.R) / 0xFF, float32(c.G) / 0xFF, float32(c.B) / 0xFF, float32(c.A) / 0xFF
}

func (c Color) WithAlpha(a uint8) Color {
	c.A = a
	return c
}

func (c Color) RGBA() (r, g, b, a uint8) { return c.R, c.G, c.B, c.A }
