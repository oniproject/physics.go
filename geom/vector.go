package geom

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func (p *Vector) Hashcode() (hash uint64) {
	x, y := uint64(p.X), uint64(p.Y)
	hash = x + y
	return
}

func (p *Vector) Equals(oi interface{}) (equals bool) {
	o, equals := oi.(*Vector)
	if !equals {
		var op Vector
		op, equals = oi.(Vector)
		equals = equals && p.EqualsVector(op)
		return
	}
	equals = p.EqualsVector(*o)
	return
}

func (p *Vector) Translate(offset Vector) { *p = p.Plus(offset) }

func (p *Vector) Rotate(rad float64) {
	p.X = p.X*math.Cos(rad) - p.Y*math.Sin(rad)
	p.Y = p.X*math.Sin(rad) + p.Y*math.Cos(rad)
}

func (p *Vector) RotateLeft()  { p.X, p.Y = -p.Y, p.X }
func (p *Vector) RotateRight() { p.X, p.Y = p.Y, -p.X }

func (p Vector) Unit() (u Vector) {
	m := p.Magnitude()
	u.X = p.X / m
	u.Y = p.Y / m
	return
}

func (p *Vector) Scale(xfactor, yfactor float64) {
	p.X *= xfactor
	p.Y *= yfactor
}

func (p Vector) EqualsVector(q Vector) bool { return p.X == q.X && p.Y == q.Y }

func (p Vector) DistanceFrom(q Vector) (d float64)         { return p.Minus(q).Magnitude() }
func (p Vector) DistanceFromSquared(q Vector) (ds float64) { return p.Minus(q).MagnitudeSquared() }

func (p Vector) Magnitude() (m float64)         { return math.Sqrt(p.MagnitudeSquared()) }
func (p Vector) MagnitudeSquared() (ms float64) { return p.X*p.X + p.Y*p.Y }

func (p Vector) Minus(q Vector) (r Vector) {
	r.X = p.X - q.X
	r.Y = p.Y - q.Y
	return
}

func (p Vector) Plus(q Vector) (r Vector) {
	r.X = p.X + q.X
	r.Y = p.Y + q.Y
	return
}

func (p Vector) Times(s float64) (r Vector) {
	r.X = p.X * s
	r.Y = p.Y * s
	return
}

func (p Vector) QuadPP(q Vector) bool { return q.X >= p.X && q.Y >= p.Y }
func (p Vector) QuadPM(q Vector) bool { return q.X >= p.X && q.Y <= p.Y }
func (p Vector) QuadMP(q Vector) bool { return q.X <= p.X && q.Y >= p.Y }
func (p Vector) QuadMM(q Vector) bool { return q.X <= p.X && q.Y <= p.Y }

func DotProduct(p, q Vector) (r float64)   { return p.X*q.X + p.Y*q.Y }
func CrossProduct(p, q Vector) (z float64) { return p.X*q.Y - p.Y*q.X }

func VectorAngle(X, Y Vector) (r float64) {
	XdotY := DotProduct(X, Y)
	mXmY := X.Magnitude() * Y.Magnitude()
	r = math.Acos(XdotY / mXmY)
	z := CrossProduct(X, Y)
	if z < 0 {
		r *= -1
	}
	return
}

func VertexAngle(A, B, C Vector) (r float64) {
	X := A.Minus(B)
	Y := C.Minus(B)
	r = VectorAngle(X, Y)
	/*if r < 0 {
		r *= -1
	}*/
	return
}

func VectorChan(points []Vector) (ch <-chan Vector) {
	tch := make(chan Vector, len(points))
	go func(points []Vector, ch chan<- Vector) {
		for _, p := range points {
			ch <- p
		}
		close(ch)
	}(points, tch)
	ch = tch
	return
}

/// ................

/*func (p *Vector) Clone(v Vector) (r Vector) {
	r.X = p.X
	r.Y = p.Y
	return
}*/

func (p *Vector) Vadd(v *Vector) *Vector {
	p.X += v.X
	p.Y += v.Y
	return p
}
func (p *Vector) Vsub(v *Vector) *Vector {
	p.X += v.X
	p.Y += v.Y
	return p
}
func (p *Vector) Mult(n float64) *Vector {
	p.X *= n
	p.Y *= n
	return p
}

func (p Vector) Angle(v *Vector) (ang float64) {
	if p.Equals(Vector{0, 0}) {
		return v.Angle(nil)
	} else {
		if v != nil && !v.Equals(Vector{0, 0}) {
			ang = math.Atan2(p.Y*v.X-p.X*v.Y, p.X*v.X+p.Y*v.Y)
		} else {
			ang = math.Atan2(p.Y, p.X)
		}
	}

	for ang > math.Pi {
		ang -= math.Pi * 2
	}
	for ang < -math.Pi {
		ang += math.Pi * 2
	}

	return
}
