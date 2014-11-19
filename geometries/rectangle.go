package geometries

import (
	"github.com/oniproject/physics.go"
	"github.com/oniproject/physics.go/geom"
)

type Rectangle struct {
	Width, Height float64
}

func NewRectangle(w, h float64) physics.Geometry {
	return &Rectangle{w, h}
}

func (this *Rectangle) AABB(angle float64) physics.AABB {
	panic("not implemented")

	if angle == 0 {
		return physics.NewAABB_byWH(this.Width, this.Height)
	}

	trans := physics.NewTransformAngle(angle)

	xaxis := trans.RotateInv(geom.Vector{1, 0})
	yaxis := trans.RotateInv(geom.Vector{0, 1})
	xmax := this.FarthestHullPoint(xaxis).Proj(xaxis)
	xmin := this.FarthestHullPoint(xaxis.Times(-1)).Proj(xaxis)
	ymax := this.FarthestHullPoint(yaxis).Proj(yaxis)
	ymin := this.FarthestHullPoint(yaxis.Times(-1)).Proj(yaxis)

	return physics.NewAABB_byMM(xmin, ymin, xmax, ymax)
}

func (this *Rectangle) FarthestHullPoint(dir geom.Vector) geom.Vector {
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

	return geom.Vector{x, y}
}

func (this *Rectangle) FarthestCorePoint(dir geom.Vector, margin float64) geom.Vector {
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

	return geom.Vector{0, 0}
}
