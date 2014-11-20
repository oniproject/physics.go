package geometries

import (
	"github.com/oniproject/physics.go/geom"
)

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) Geometry {
	return &Circle{radius}
}

func (this *Circle) AABB(angle float64) geom.AABB {
	r := this.Radius
	return geom.NewAABB_byMM(-r, -r, r, r)
}
func (this *Circle) FarthestHullPoint(dir geom.Vector) geom.Vector {
	n := dir.Unit()
	return n.Times(this.Radius)
}
func (this *Circle) FarthestCorePoint(dir geom.Vector, margin float64) geom.Vector {
	n := dir.Unit()
	return n.Times(this.Radius - margin)
}
