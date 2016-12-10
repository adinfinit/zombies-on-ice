package g

var Zero2 = V2{}

type V2 struct{ X, Y float32 }

// Add adds two vectors and returns the result
func (a V2) Add(b V2) V2               { return V2{a.X + b.X, a.Y + b.Y} }
func (a V2) AddMul(b V2, s float32) V2 { return V2{a.X + b.X*s, a.Y + b.Y*s} }

// Sub subtracts two vectors and returns the result
func (a V2) Sub(b V2) V2 { return V2{a.X - b.X, a.Y - b.Y} }

// Scale scales each component and returns the result
func (a V2) Scale(s float32) V2 { return V2{a.X * s, a.Y * s} }

// Length returns the length of the vector
func (a V2) Length() float32 { return Sqrt(a.X*a.X + a.Y*a.Y) }

// Length2 returns the squared length of the vector
func (a V2) Length2() float32 { return a.X*a.X + a.Y*a.Y }

// Distance returns the distance to vector b
func (a V2) Distance(b V2) float32 {
	dx, dy := a.X-b.X, a.Y-b.Y
	return Sqrt(dx*dx + dy*dy)
}

// Distance2 returns the squared distance to vector b
func (a V2) Distance2(b V2) float32 {
	dx, dy := a.X-b.X, a.Y-b.Y
	return dx*dx + dy*dy
}

func (a V2) UnitNormalize() V2 {
	m := a.Length()
	if m < 1 {
		m = 1
	}
	return V2{a.X / m, a.Y / m}
}

// Cross product of a and b
func (a V2) Cross(b V2) float32 { return a.X*b.Y - a.Y*b.X }

func (a V2) NearZero() bool { return a.Length2() < 0.0001 }

func (a V2) Rotate(angle float32) V2 {
	cs, sn := Cos(angle), Sin(angle)
	return V2{a.X*cs - a.Y*sn, a.X*sn + a.Y*cs}
}
