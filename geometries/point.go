package geometries

import (
	"github.com/oniproject/physics.go/geom"
)

type Point struct{}

func NewPoint() Geometry {
	return &Point{}
}

func (this *Point) AABB(angle float64) geom.AABB {
	return geom.AABB{}
}
func (this *Point) FarthestHullPoint(dir geom.Vector) geom.Vector {
	// not implemented.
	return geom.Vector{0, 0}
}
func (this *Point) FarthestCorePoint(dir geom.Vector, margin float64) geom.Vector {
	// not implemented.
	return geom.Vector{0, 0}
}
