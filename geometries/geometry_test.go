package geometries

import (
	"github.com/oniproject/physics.go/geom"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Geometries(t *testing.T) {
	Convey("test Geometries", t, func() {
		Convey("test if a polygon is convex", func() {
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

			So(IsPolygonConvex(convex), ShouldBeTrue)
			So(IsPolygonConvex(convexReverse), ShouldBeTrue)
			So(IsPolygonConvex(notConvex), ShouldBeFalse)
		})

		Convey("check moments of inertia of polygons", func() {
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

			So(PolygonMOI(point), ShouldEqual, 0)
			So(PolygonMOI(line), ShouldEqual, 4.0/12.0)
			So(PolygonMOI(square), ShouldEqual, 2.0*2.0/6.0)
		})

		Convey("check if points are inside a polygon", func() {
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

			So(IsPointInPolygon(inside, line), ShouldBeTrue)
			So(IsPointInPolygon(inside, square), ShouldBeTrue)
			So(IsPointInPolygon(outside, line), ShouldBeFalse)
			So(IsPointInPolygon(outside, square), ShouldBeFalse)
		})

		Convey("calculate polygon area", func() {
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

			So(PolygonArea(point), ShouldEqual, 0)
			So(PolygonArea(line), ShouldEqual, 0)
			So(PolygonArea(square), ShouldEqual, 25)
			So(PolygonArea(squareReverse), ShouldEqual, -25)
		})

		Convey("calculate polygon centroid", func() {
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

			So(PolygonCentroid(point), ShouldResemble, geom.Vector{9, 9})
			So(PolygonCentroid(line), ShouldResemble, geom.Vector{3, 0})
			So(PolygonCentroid(square), ShouldResemble, geom.Vector{1.5, 1.5})
		})

		Convey("check nearest point on a line", func() {
			var result geom.Vector
			line1 := geom.Vector{0, 2}
			line2 := geom.Vector{6, 8}

			result = NearestPointOnLine(geom.Vector{1, -5}, line1, line2)
			So(result, ShouldResemble, line1)

			result = NearestPointOnLine(geom.Vector{2, 2}, line1, line2)
			So(result, ShouldResemble, geom.Vector{1, 3})

			result = NearestPointOnLine(geom.Vector{10, 8}, line1, line2)
			So(result, ShouldResemble, line2)
		})
	})
}
