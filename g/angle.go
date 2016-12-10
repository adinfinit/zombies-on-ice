package g

import "math"

const Pi = math.Pi
const Tau = Pi * 2

type Radian float32

func DegToRad(deg float32) float32 { return deg * Pi / 180 }
func RadToDeg(rad float32) float32 { return rad * 180 / Pi }

func Cos(v float32) float32      { return float32(math.Cos(float64(v))) }
func Sin(v float32) float32      { return float32(math.Sin(float64(v))) }
func Atan(v float32) float32     { return float32(math.Atan(float64(v))) }
func Atan2(y, x float32) float32 { return float32(math.Atan2(float64(y), float64(x))) }
