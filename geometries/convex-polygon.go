package geometries

import (
	_ "github.com/oniproject/physics.go"
)

type ConvexPolygon struct {
}

func NewConvexPolygon() Geometry {
	panic("not implemented")
	return &Point{}
}

func (this *ConvexPolygon) AABB(angle float64) AABB {
	panic("not implemented")
	return AABB{}
}
func (this *ConvexPolygon) FarthestCorePoint(dir Vector) Vector {
	panic("not implemented")
	// not implemented.
	return Vector{0, 0}
}
func (this *ConvexPolygon) FarthestHullPoint(dir Vector, margin float64) Vector {
	panic("not implemented")
	// not implemented.
	return Vector{0, 0}
}
