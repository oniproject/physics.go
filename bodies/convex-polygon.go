package bodies

import (
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
)

type ConvexPolygon struct {
	Point
}

func NewConvexPolygon(vertices []geom.Vector) Body {
	c := &ConvexPolygon{Point: *NewPoint()}
	c.geometry = geometries.NewConvexPolygon(vertices)
	c.Recalc()
	return c
}

func (this *ConvexPolygon) Recalc() {
	v := this.geometry.(*geometries.ConvexPolygon)
	// moment of inertia
	this.moi = geometries.PolygonMOI(v.Vertices)
}
