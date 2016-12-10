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

// Clamp1 clamps value to 0..1
func Clamp1(a float32) float32 {
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
