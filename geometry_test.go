package physics

import (
	"github.com/oniproject/physics.go/geom"
	"testing"
)

func Test_Geometries(t *testing.T) {
	It(t, "test if a polygon is convex", func() bool {
		convex := []geom.Vector{
			geom.Vector{X: 1, Y: 0},
			geom.Vector{X: 0.5, Y: -0.5},
			geom.Vector{X: -0.5, Y: -0.5},
			geom.Vector{X: -1, Y: 0},
			geom.Vector{X: 0, Y: 1},
		}
		convexReverse := []geom.Vector{
			geom.Vector{X: 0, Y: 1},
			geom.Vector{X: -1, Y: 0},
			geom.Vector{X: -0.5, Y: -0.5},
			geom.Vector{X: 0.5, Y: -0.5},
			geom.Vector{X: 1, Y: 0},
		}
		notConvex := []geom.Vector{
			geom.Vector{X: 1, Y: 0},
			geom.Vector{X: 0.5, Y: -0.5},
			geom.Vector{X: -0.5, Y: -0.5},
			geom.Vector{X: -1, Y: 0},
			geom.Vector{X: 2.3, Y: 1},
		}

		if !IsPolygonConvex(convex) {
			t.Error("convex")
			return false
		}
		if !IsPolygonConvex(convexReverse) {
			t.Error("convexReverse")
			return false
		}
		if IsPolygonConvex(notConvex) {
			t.Error("notConvex")
			return false
		}
		return true
	})

	It(t, "check moments of inertia of polygons", func() bool {
		point := []geom.Vector{
			geom.Vector{X: 0, Y: 0},
		}
		line := []geom.Vector{
			geom.Vector{X: -1, Y: 0},
			geom.Vector{X: 1, Y: 0},
		}
		square := []geom.Vector{
			geom.Vector{X: 1, Y: 1},
			geom.Vector{X: 1, Y: -1},
			geom.Vector{X: -1, Y: 1},
			geom.Vector{X: -1, Y: -1},
		}

		if PolygonMOI(point) != 0 {
			t.Error(point)
			return false
		}
		if PolygonMOI(line) != 4.0/12.0 {
			t.Error(line, PolygonMOI(line), 4.0/12.0)
			return false
		}
		if PolygonMOI(square) != 2.0*2.0/6.0 {
			t.Error(square, PolygonMOI(square), "!=", 2.0*2.0/6.0)
			return false
		}
		return true
	})

	It(t, "check if points are inside a polygon", func() bool {
		inside := geom.Vector{1, 0}
		outside := geom.Vector{5, 5.1}
		line := []geom.Vector{
			geom.Vector{X: -5, Y: 0},
			geom.Vector{X: 5, Y: 0},
		}
		square := []geom.Vector{
			geom.Vector{X: 5, Y: 5},
			geom.Vector{X: 5, Y: -5},
			geom.Vector{X: -5, Y: 5},
			geom.Vector{X: -5, Y: -5},
		}

		if !IsPointInPolygon(inside, line) {
			t.Error("inside line", line)
			return false
		}
		if !IsPointInPolygon(inside, square) {
			t.Error("inside square", square)
			return false
		}
		if IsPointInPolygon(outside, line) {
			t.Error("outside line", line)
			return false
		}
		if IsPointInPolygon(outside, square) {
			t.Error("outside square", square)
			return false
		}
		return true
	})

	It(t, "calculate polygon area", func() bool {
		point := []geom.Vector{
			geom.Vector{X: 9, Y: 9},
		}
		line := []geom.Vector{
			geom.Vector{X: 1, Y: 0},
			geom.Vector{X: 5, Y: 0},
		}
		square := []geom.Vector{
			geom.Vector{X: 5, Y: 5},
			geom.Vector{X: 5, Y: 0},
			geom.Vector{X: 0, Y: 0},
			geom.Vector{X: 0, Y: 5},
		}
		squareReverse := []geom.Vector{
			geom.Vector{X: 0, Y: 5},
			geom.Vector{X: 0, Y: 0},
			geom.Vector{X: 5, Y: 0},
			geom.Vector{X: 5, Y: 5},
		}

		if PolygonArea(point) != 0 {
			t.Error("point")
			return false
		}
		if PolygonArea(line) != 0 {
			t.Error("line")
			return false
		}
		if PolygonArea(square) != 25 {
			t.Error("square", PolygonArea(square))
			return false
		}
		if PolygonArea(squareReverse) != -25 {
			t.Error("squareReverse", PolygonArea(squareReverse))
			return false
		}
		return true
	})

	It(t, "calculate polygon centroid", func() bool {
		point := []geom.Vector{
			geom.Vector{X: 9, Y: 9},
		}
		line := []geom.Vector{
			geom.Vector{X: 1, Y: 0},
			geom.Vector{X: 5, Y: 0},
		}
		square := []geom.Vector{
			geom.Vector{X: 3, Y: 3},
			geom.Vector{X: 3, Y: 0},
			geom.Vector{X: 0, Y: 0},
			geom.Vector{X: 0, Y: 3},
		}

		var centroid geom.Vector

		centroid = PolygonCentroid(point)
		if !centroid.Equals(geom.Vector{9, 9}) {
			t.Error("point", centroid)
			return false
		}

		centroid = PolygonCentroid(line)
		if !centroid.Equals(geom.Vector{3, 0}) {
			t.Error("line", centroid)
			return false
		}

		centroid = PolygonCentroid(square)
		if !centroid.Equals(geom.Vector{1.5, 1.5}) {
			t.Error("square", centroid)
			return false
		}

		return true
	})

	It(t, "check nearest point on a line", func() bool {
		var point, result geom.Vector
		line1 := geom.Vector{0, 2}
		line2 := geom.Vector{6, 8}

		point.X, point.Y = 1, -5
		result = NearestPointOnLine(point, line1, line2)
		if result.X != line1.X || result.Y != line1.Y {
			t.Error(point, result)
			return false
		}

		point.X, point.Y = 2, 2
		result = NearestPointOnLine(point, line1, line2)
		if result.X != 1 || result.Y != 3 {
			t.Error(point, result)
			return false
		}

		point.X, point.Y = 10, 8
		result = NearestPointOnLine(point, line1, line2)
		if result.X != line2.X || result.Y != line2.Y {
			t.Error(point, result)
			return false
		}

		return true
	})
}
