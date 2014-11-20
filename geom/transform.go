package geom

import (
	"math"
)

type Transform struct {
	Vect   Vector
	Angle  float64
	Origin Vector

	CosA, SinA float64
}

func NewTransformAngle(angle float64) *Transform {
	return NewTransform(Vector{0, 0}, angle, Vector{0, 0})
}
func NewTransform(vect Vector, angle float64, origin Vector) (t *Transform) {
	t = &Transform{}
	t.SetTranslation(vect)
	t.SetRotation(angle, origin)
	return
}

func (t *Transform) SetTranslation(vect Vector) {
	t.Vect = vect
}
func (t *Transform) SetRotation(angle float64, origin Vector) {
	t.CosA = math.Cos(angle)
	t.SinA = math.Sin(angle)
	t.Origin = origin
}

func (t *Transform) Translate(vect Vector) Vector {
	return vect.Plus(t.Vect)
}

func (t *Transform) RotateInv(vect Vector) Vector {
	return Vector{
		X: (vect.X-t.Origin.X)*t.CosA + (vect.Y-t.Origin.Y)*t.SinA + t.Origin.X,
		Y: -(vect.X-t.Origin.X)*t.SinA + (vect.Y-t.Origin.Y)*t.CosA + t.Origin.Y,
	}
}
