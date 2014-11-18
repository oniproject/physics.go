package geometries

import (
	_ "github.com/oniproject/physics.go"
)

type Rectangle struct {
	Width, Height float64
}

func NewRectangle(w, h float64) Geometry {
	return &Rectangle{w, h}
}

func (this *Rectangle) AABB(angle float64) AABB {
	panic("not implemented")

	if angle == 0 {
		return NewAABB_byWH(this.width, this.height)
	}

	//var scratch = Physics.scratchpad()
	//,p = scratch.vector()

	/* FIXME trans := scratch.transform().setRotation(angle || 0)
	xaxis := scratch.vector().set(1, 0).rotateInv(trans)
	yaxis := scratch.vector().set(0, 1).rotateInv(trans)
	xmax := this.FarthestHullPoint(xaxis, p).proj(xaxis)
	xmin := -this.FarthestHullPoint(xaxis.negate(), p).proj(xaxis)
	ymax := this.FarthestHullPoint(yaxis, p).proj(yaxis)
	ymin := -this.FarthestHullPoint(yaxis.negate(), p).proj(yaxis)
	*/
	// scratch.done();*/

	return NewAABB_byMM(xmin, ymin, xmax, ymax)
}

func (this *Rectangle) FarthestCorePoint(dir Vector) Vector {
	x, y := dir.X, dir.Y

	switch {
	case x < 0:
		x = -this.Width * 0.5
	case x > 0:
		x = this.Width * 0.5
	default:
		x = 0
	}

	switch {
	case y < 0:
		y = -this.Height * 0.5
	case x > 0:
		y = this.Height * 0.5
	default:
		y = 0
	}

	return Vector{x, y}
}

func (this *Rectangle) FarthestHullPoint(dir Vector, margin float64) Vector {
	result := this.FarthestHullPoint(dir)
	x, y := result.X, result.Y

	switch {
	case x < 0:
		x = -margin * 0.5
	case x > 0:
		x = margin * 0.5
	default:
		x = 0
	}

	switch {
	case y < 0:
		y = -margin * 0.5
	case x > 0:
		y = margin * 0.5
	default:
		y = 0
	}

	return Vector{0, 0}
}
