package geometries

import (
	"github.com/oniproject/physics.go"
	"github.com/oniproject/physics.go/geom"
)

type Circle struct {
	Radius float64
	aabb   physics.AABB
}

func NewCircle(radius float64) physics.Geometry {
	return &Circle{radius, physics.AABB{}}
}

func (this *Circle) AABB(angle float64) physics.AABB {
	if this.aabb.HW != r {
		this.aabb = NewAABB_byMM(-r, -r, r, r)
	}
	return this.aabb
}
func (this *Circle) FarthestCorePoint(dir geom.Vector) geom.Vector {
	return dir.Normalize().Mult(this.Radius)
}
func (this *Circle) FarthestHullPoint(dir geom.Vector, margin float64) geom.Vector {
	return dir.Normalize().Mult(this.Radius - margin)
}
