package g

// Clamp clamps value to low..high
func Clamp(a, low, high float32) float32 {
	if a < low {
		return low
	} else if a > high {
		return high
	}
	return a
}

// Clamp01 clamps value to 0..1
func Clamp01(a float32) float32 {
	if a < 0 {
		return 0
	} else if a > 1 {
		return 1
	}
	return a
}

// ClampUnit clamps value to -1..1
//TODO: what's the conventional name for this?
func ClampUnit(a float32) float32 {
	if a < -1 {
		return -1
	} else if a > 1 {
		return 1
	}
	return a
}

func ClampLength(a V2, maxsize float32) V2 {
	n := a.Length()
	if n > maxsize {
		return a.Scale(maxsize / n)
	}
	return a
}
