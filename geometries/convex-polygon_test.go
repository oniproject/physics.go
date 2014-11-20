package geometries

import (
	"github.com/oniproject/physics.go/geom"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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

	Convey("test ConvexPolygon", t, func() {

		// It(t, "polygons must have verti

		Convey("check centroid repositioning", func() {
			verts := square.Vertices
			trues := []geom.Vector{
				{2.5, 2.5},
				{2.5, -2.5},
				{-2.5, -2.5},
				{-2.5, 2.5},
			}
			for i := range verts {
				So(verts[i], ShouldResemble, trues[i])
			}
		})

		Convey("check farthest hull points pentagon", func() {
			var pt geom.Vector

			pt = poly.FarthestHullPoint(geom.Vector{1, 0.1})
			So(pt, ShouldResemble, poly.Vertices[0])

			pt = poly.FarthestHullPoint(geom.Vector{0.1, 1})
			So(pt, ShouldResemble, poly.Vertices[4])

			pt = poly.FarthestHullPoint(geom.Vector{0.1, -1})
			So(pt, ShouldResemble, poly.Vertices[1])

			pt = poly.FarthestHullPoint(geom.Vector{-0.1, -1})
			So(pt, ShouldResemble, poly.Vertices[2])
		})
		Convey("check farthest hull points triangle", func() {
			var pt geom.Vector

			pt = triangle.FarthestHullPoint(geom.Vector{1, -0.1})
			So(pt, ShouldResemble, triangle.Vertices[1])

			pt = triangle.FarthestHullPoint(geom.Vector{1, 1})
			So(pt, ShouldResemble, triangle.Vertices[0])

			pt = triangle.FarthestHullPoint(geom.Vector{-1, -0.1})
			So(pt, ShouldResemble, triangle.Vertices[2])
		})
		Convey("check aabb", func() {
			aabb := poly.AABB(0)
			So(aabb.HW, ShouldEqual, 1)
			So(aabb.HH, ShouldEqual, 1.5/2.0)
		})
		Convey("sould return correct vertices for FarthestHullPoint", func() {
			var v geom.Vector

			v = shape.FarthestHullPoint(geom.Vector{1, 0})
			So(v, ShouldResemble, shape.Vertices[1])

			v = shape.FarthestHullPoint(geom.Vector{-1, 0})
			So(v, ShouldResemble, shape.Vertices[0])

			v = shape.FarthestHullPoint(geom.Vector{0, 1})
			So(v, ShouldResemble, shape.Vertices[1])

			v = shape.FarthestHullPoint(geom.Vector{0, -1})
			So(v, ShouldResemble, shape.Vertices[2])
		})
	})
}
