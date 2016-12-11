package g

import (
	"math"
	"math/rand"
)

// Sqr returns the square of v
func Sqr(v float32) float32 { return v * v }

// Sqrt returns the square root of v
func Sqrt(v float32) float32 {
	//TODO: optimize
	return float32(math.Sqrt(float64(v)))
}

func Pow(v, e float32) float32 {
	return float32(math.Pow(float64(v), float64(e)))
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

func RandomBetween(a, b float32) float32 {
	return rand.Float32()*(b-a) + a
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

func MinMax(a, b float32) (float32, float32) {
	if a < b {
		return a, b
	}
	return b, a
}

func Abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}
