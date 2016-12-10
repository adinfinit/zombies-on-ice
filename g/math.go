package g

import "math"

// Sqr returns the square of v
func Sqr(v float32) float32 { return v * v }

// Sqrt returns the square root of v
func Sqrt(v float32) float32 {
	//TODO: optimize
	return float32(math.Sqrt(float64(v)))
}

func ApplyDeadZone(v float32, deadZone float32) float32 {
	if v < -deadZone {
		return (v + deadZone) / (1 - deadZone)
	}
	if v > deadZone {
		return (v - deadZone) / (1 - deadZone)
	}
	return 0.0
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
