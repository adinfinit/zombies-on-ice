package g

func Lerp(a, b, p float32) float32      { return (b-a)*p + a }
func LerpClamp(a, b, p float32) float32 { return Lerp(a, b, Clamp01(p)) }
