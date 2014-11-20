package bodies

import (
	//"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
)

type Circle struct {
	Point
}

func NewCircle(radius float64) Body {
	//c := &Circle{geometry: geometries.NewCircle(radius)}
	c := &Circle{Point: *NewPoint()}
	c.geometry = geometries.NewCircle(radius)
	c.Recalc()
	return c
}

func (this *Circle) Recalc() {
	c := this.geometry.(*geometries.Circle)
	// moment of inertia
	this.moi = this.mass * c.Radius * c.Radius / 2.0
}
