package g

type Rect struct{ Min, Max V2 }

func NewRect(w, h float32) Rect {
	return Rect{
		V2{-w / 2, -h / 2},
		V2{w / 2, h / 2},
	}
}

// scales inner such that it matches bounds aspect ratio
func NewCenteredRect(inner, bounds V2, margin float32) Rect {
	boundsRatio := bounds.Y / bounds.X
	innerRatio := inner.Y / inner.X

	var size V2
	if boundsRatio < innerRatio {
		size.Y = inner.Y + margin
		size.X = size.Y / boundsRatio
	} else {
		size.X = inner.X + margin
		size.Y = size.X * boundsRatio
	}

	return NewRect(size.X, size.Y)
}

func NewCircleRect(r float32) Rect {
	return Rect{V2{-r, -r}, V2{r, r}}
}

func (r Rect) Size() V2 { return r.Max.Sub(r.Min) }

func (r Rect) Offset(delta V2) Rect {
	return Rect{
		r.Min.Add(delta),
		r.Max.Add(delta),
	}
}
func (r Rect) ScaleInv(v V2) Rect {
	return Rect{
		V2{
			r.Min.X / v.X,
			r.Min.Y / v.Y,
		},
		V2{
			r.Max.X / v.X,
			r.Max.Y / v.Y,
		},
	}
}

func (r Rect) Contains(p V2) bool {
	return (r.Min.X <= p.X) && (p.X <= r.Max.X) &&
		(r.Min.Y <= p.Y) && (p.Y <= r.Max.Y)
}

func EnforceInside(pos, vel *V2, bounds Rect, dampening float32) {
	minx, maxx := MinMax(bounds.Min.X, bounds.Max.X)
	if pos.X < minx {
		pos.X = minx
		vel.X = +Abs(vel.X) * dampening
	}
	if pos.X > maxx {
		pos.X = maxx
		vel.X = -Abs(vel.X) * dampening
	}

	miny, maxy := MinMax(bounds.Min.Y, bounds.Max.Y)
	if pos.Y < miny {
		pos.Y = miny
		vel.Y = +Abs(vel.Y) * dampening
	}
	if pos.Y > maxy {
		pos.Y = maxy
		vel.Y = -Abs(vel.Y) * dampening
	}
}
