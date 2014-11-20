package bodies

import (
	//"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
)

type Rectangle struct {
	Point
}

func NewRectangle(w, h float64) Body {
	r := &Rectangle{}
	r.geometry = geometries.NewRectangle(w, h)
	r.Recalc()
	return r
}

func (this *Rectangle) Recalc() {
	r := this.geometry.(*geometries.Rectangle)
	w := r.Width
	h := r.Height

	// moment of inertia
	this.moi = (w*w + h*h) * this.mass / 12.0
}
