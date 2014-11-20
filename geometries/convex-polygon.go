package geometries

import (
	"errors"
	"github.com/oniproject/physics.go/geom"
)

var ERROR_NOT_CONVEX = errors.New("Error: The vertices specified do not mathc that of a _convex_ polygon.")

type ConvexPolygon struct {
	Vertices []geom.Vector
	area     float64
	aabb     *geom.AABB
}

func NewConvexPolygon(hull []geom.Vector) (cp *ConvexPolygon) {
	cp = &ConvexPolygon{aabb: nil}
	cp.setVertices(hull)
	return
}

func (this *ConvexPolygon) setVertices(hull []geom.Vector) {
	if !IsPolygonConvex(hull) {
		panic(ERROR_NOT_CONVEX)
	}

	center := PolygonCentroid(hull).Times(-1)
	trans := geom.NewTransform(center, 0, geom.Vector{})

	this.Vertices = []geom.Vector{}
	for _, vert := range hull {
		this.Vertices = append(this.Vertices, trans.Translate(vert))
	}

	this.area = PolygonArea(this.Vertices)

	this.aabb = nil
}

func (this *ConvexPolygon) AABB(angle float64) geom.AABB {
	if angle == 0 && this.aabb != nil {
		return *this.aabb
	}

	trans := geom.NewTransformAngle(angle)

	xaxis := trans.RotateInv(geom.Vector{1, 0})
	yaxis := trans.RotateInv(geom.Vector{0, 1})

	xmax := this.FarthestHullPoint(xaxis).Proj(xaxis)
	xmin := this.FarthestHullPoint(xaxis.Times(-1)).Proj(xaxis)
	ymax := this.FarthestHullPoint(yaxis).Proj(yaxis)
	ymin := this.FarthestHullPoint(yaxis.Times(-1)).Proj(yaxis)

	aabb := geom.NewAABB_byMM(xmin, ymin, xmax, ymax)

	if angle == 0 {
		this.aabb = &aabb
	}

	return aabb
}

func (this *ConvexPolygon) FarthestHullPoint(dir geom.Vector) geom.Vector {
	verts := this.Vertices

	if len(verts) < 2 {
		return verts[0]
	}

	prev := geom.DotProduct(verts[0], dir)
	val := geom.DotProduct(verts[1], dir)

	if len(verts) == 2 {
		if val >= prev {
			return verts[1]
		} else {
			return verts[0]
		}
	}

	i := 2
	if val >= prev {
		// go up
		// seatch until the next dot product is less than the prev
		for ; i < len(verts) && val >= prev; i++ {
			prev = val
			val = geom.DotProduct(verts[i], dir)
		}

		if val >= prev {
			i++
		}

		return verts[i-2]
	} else {
		// go down
		i = len(verts)
		for i > 1 && prev >= val {
			i--
			val = prev
			prev = geom.DotProduct(verts[i], dir)
		}

		return verts[(i+1)%len(verts)]
	}
}

func (this *ConvexPolygon) FarthestCorePoint(dir geom.Vector, margin float64) geom.Vector {
	panic("not implemented")
	return geom.Vector{0, 0}
}
