package g

type Color struct{ R, G, B, A uint8 }

func LerpColor(a, b Color, p float32) Color {
	ar, ag, ab, aa := a.Float()
	br, bg, bb, ba := b.Float()

	return ColorFloat(
		LerpClamp(ar, br, p),
		LerpClamp(ag, bg, p),
		LerpClamp(ab, bb, p),
		LerpClamp(aa, ba, p),
	)
}

var (
	White = Color{0xFF, 0xFF, 0xFF, 0xFF}
	Black = Color{0x00, 0x00, 0x00, 0xFF}
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

func ColorFloat(r, g, b, a float32) Color {
	return Color{Sat8(r), Sat8(g), Sat8(b), Sat8(a)}
}

func ColorHSLA(h, s, l, a float32) Color { return ColorFloat(hsla(h, s, l, a)) }
func ColorHSL(h, s, l float32) Color     { return ColorHSLA(h, s, l, 1) }

func (c Color) Float() (r, g, b, a float32) {
	return float32(c.R) / 0xFF, float32(c.G) / 0xFF, float32(c.B) / 0xFF, float32(c.A) / 0xFF
}

func (c Color) WithAlpha(a uint8) Color {
	c.A = a
	return c
}

func (c Color) RGBA() (r, g, b, a uint8) { return c.R, c.G, c.B, c.A }

func hue(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	} else if 2*h < 1 {
		return v2
	} else if 3*h < 2 {
		return v1 + (v2-v1)*(2.0/3.0-h)*6
	}

	return v1
}

func hsla(h, s, l, a float32) (r, g, b, ra float32) {
	if s == 0 {
		return l, l, l, a
	}

	h = Mod(h, 1)

	var v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - s*l
	}

	v1 := 2*l - v2
	r = hue(v1, v2, h+1.0/3.0)
	g = hue(v1, v2, h)
	b = hue(v1, v2, h-1.0/3.0)
	ra = a

	return
}

func Sat8(v float32) uint8 {
	v *= 255.0
	if v >= 255 {
		return 255
	} else if v <= 0 {
		return 0
	}
	return uint8(v)
}
