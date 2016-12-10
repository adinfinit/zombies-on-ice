package g

type Rect struct{ Min, Max V2 }

func NewRect(w, h float32) Rect {
	return Rect{
		V2{-w / 2, -h / 2},
		V2{w / 2, h / 2},
	}
}

func NewCircleRect(r float32) Rect {
	return Rect{V2{-r, -r}, V2{r, r}}
}

func (r Rect) Size() V2 { return r.Max.Sub(r.Min) }

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
