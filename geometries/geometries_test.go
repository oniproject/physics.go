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
func Test_Geometry(t *testing.T) {
	/*It("test if a polygon is convex", func() bool {
		convex := [
			{X: 1, Y: 0},
			{X:0.5, Y: -0.5},
			{X:-0.5, Y: 0.5},
			{X: -1, Y: 0},
			{X: 0, Y: 1},
		]
		notConvex := [
			{X: 1, Y: 0},
			{X:0.5, Y: -0.5},
			{X:-0.5, Y: 0.5},
			{X: -1, Y: 0},
			{X: 2.3, Y: 1},
		]
		if !IsPolygonConvex(convex) {
			return false
		}
		if IsPolygonConvex(notConvex) {
			return false
		}
		convex.Reverse()
		if !IsPolygonConvex(convex) {
			return false
		}
	})*/

	It("check moments of inertia of polygons", func() bool {
		point := []Vector{
			Vector{X: 0, Y: 0},
		}
		line := []Vector{
			Vector{X: -1, Y: 0},
			Vector{X: 1, Y: 0},
		}
		square := []Vector{
			Vector{X: 1, Y: 1},
			Vector{X: 1, Y: -1},
			Vector{X: -1, Y: 1},
			Vector{X: -1, Y: -1},
		}

		if PolygonMOI(point) != 0 {
			return false
		}
		if PolygonMOI(line) != 4/12 {
			return false
		}
		if PolygonMOI(square) != 2*2/6 {
			return false
		}
	})
}
