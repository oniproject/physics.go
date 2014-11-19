package geometries

import (
	"github.com/oniproject/physics.go/geom"
	"testing"
)

func It(t *testing.T, str string, fn func() bool) {
	if !fn() {
		t.Error("It", str)
	}
}

func Test_ConvexPolygon(t *testing.T) {
	square := NewConvexPolygon([]geom.Vector{
		{5, 5},
		{5, 0},
		{0, 0},
		{0, 5},
	})
	triangle := NewConvexPolygon([]geom.Vector{
		{5, 5},
		{5, 0},
		{0, 0},
	})
	poly := NewConvexPolygon([]geom.Vector{
		{1, 0},
		{0.5, -0.5},
		{-0.5, -0.5},
		{-1, 0},
		{0, 1},
	})

	shape := NewConvexPolygon([]geom.Vector{
		{0, 242},
		{300, 242.01},
		{150, 45},
	})

	// It(t, "polygons must have verti

	It(t, "check centroid repositioning", func() bool {
		verts := square.Vertices
		trues := []geom.Vector{
			{2.5, 2.5},
			{2.5, -2.5},
			{-2.5, -2.5},
			{-2.5, 2.5},
		}
		for i := range verts {
			if verts[i].X != trues[i].X || verts[i].Y != trues[i].Y {
				t.Error(i, verts[i])
				return false
			}
		}
		return true
	})

	It(t, "check farthest hull points pentagon", func() bool {
		var dir, pt geom.Vector
		var i int

		dir = geom.Vector{1, 0.1}
		i, pt = 0, poly.FarthestHullPoint(dir)
		if pt.X != poly.Vertices[i].X || pt.Y != poly.Vertices[i].Y {
			t.Error(i, pt, poly.Vertices[i])
			return false
		}
		dir = geom.Vector{0.1, 1}
		i, pt = 4, poly.FarthestHullPoint(dir)
		if pt.X != poly.Vertices[i].X || pt.Y != poly.Vertices[i].Y {
			t.Error(i, pt, poly.Vertices[i])
			return false
		}
		dir = geom.Vector{0.1, -1}
		i, pt = 1, poly.FarthestHullPoint(dir)
		if pt.X != poly.Vertices[i].X || pt.Y != poly.Vertices[i].Y {
			t.Error(i, pt, poly.Vertices[i])
			return false
		}
		dir = geom.Vector{-0.1, -1}
		i, pt = 2, poly.FarthestHullPoint(dir)
		if pt.X != poly.Vertices[i].X || pt.Y != poly.Vertices[i].Y {
			t.Error(i, pt, poly.Vertices[i])
			return false
		}

		return true
	})
	It(t, "check farthest hull points triangle", func() bool {
		var dir, pt geom.Vector
		var i int

		dir = geom.Vector{1, -0.1}
		i, pt = 1, triangle.FarthestHullPoint(dir)
		if pt.X != triangle.Vertices[i].X || pt.Y != triangle.Vertices[i].Y {
			t.Error(i, pt, triangle.Vertices[i])
			return false
		}
		dir = geom.Vector{1, 1}
		i, pt = 0, triangle.FarthestHullPoint(dir)
		if pt.X != triangle.Vertices[i].X || pt.Y != triangle.Vertices[i].Y {
			t.Error(i, pt, triangle.Vertices[i])
			return false
		}
		dir = geom.Vector{-1, -0.1}
		i, pt = 2, triangle.FarthestHullPoint(dir)
		if pt.X != triangle.Vertices[i].X || pt.Y != triangle.Vertices[i].Y {
			t.Error(i, pt, triangle.Vertices[i])
			return false
		}

		return true
	})
	It(t, "check aabb", func() bool {
		aabb := poly.AABB(0)
		if aabb.HW != 1 || aabb.HH != 1.5/2.0 {
			t.Error(1, 1.5/2.0, aabb)
			return false
		}
		return true
	})
	It(t, "sould return correct vertices for FarthestHullPoint", func() bool {
		var v geom.Vector
		var i int

		i, v = 1, shape.FarthestHullPoint(geom.Vector{1, 0})
		if !v.Equals(shape.Vertices[i]) {
			t.Error(i, v, shape.Vertices[i])
			return false
		}
		i, v = 0, shape.FarthestHullPoint(geom.Vector{-1, 0})
		if !v.Equals(shape.Vertices[i]) {
			t.Error(i, v, shape.Vertices[i])
			return false
		}
		i, v = 1, shape.FarthestHullPoint(geom.Vector{0, 1})
		if !v.Equals(shape.Vertices[i]) {
			t.Error(i, v, shape.Vertices[i])
			return false
		}
		i, v = 2, shape.FarthestHullPoint(geom.Vector{0, -1})
		if !v.Equals(shape.Vertices[i]) {
			t.Error(i, v, shape.Vertices[i])
			return false
		}

		return true
	})
}
