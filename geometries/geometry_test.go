package geometries

import (
	"github.com/oniproject/physics.go/geom"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func Test_Geometries(t *testing.T) {

	Convey("objects", t, func() {
		Convey("Point", func() {
			p := NewPoint().(*Point)
			Convey("aabb", func() {
				So(p.AABB(0), ShouldResemble, geom.AABB{})
			})
			Convey("FarthestHullPoint", func() {
				So(p.FarthestHullPoint(geom.Vector{1, 0}), ShouldResemble, geom.Vector{0, 0})
				So(p.FarthestHullPoint(geom.Vector{0, 1}), ShouldResemble, geom.Vector{0, 0})
			})
			Convey("FarthestCorePoint", func() {
				So(p.FarthestCorePoint(geom.Vector{1, 0}, 10), ShouldResemble, geom.Vector{0, 0})
				So(p.FarthestCorePoint(geom.Vector{0, 1}, 10), ShouldResemble, geom.Vector{0, 0})
			})
		})

		Convey("Circle", func() {
			c := NewCircle(20).(*Circle)
			Convey("should with radius", func() {
				So(c.Radius, ShouldEqual, 20)
			})
			Convey("aabb", func() {
				So(c.AABB(0), ShouldResemble, geom.AABB{HW: 20, HH: 20})
			})
			Convey("FarthestHullPoint", func() {
				So(c.FarthestHullPoint(geom.Vector{1, 0}), ShouldResemble, geom.Vector{20, 0})
				So(c.FarthestHullPoint(geom.Vector{0, 1}), ShouldResemble, geom.Vector{0, 20})
			})
			Convey("FarthestCorePoint", func() {
				So(c.FarthestCorePoint(geom.Vector{1, 0}, 10), ShouldResemble, geom.Vector{10, 0})
				So(c.FarthestCorePoint(geom.Vector{0, 1}, 10), ShouldResemble, geom.Vector{0, 10})
			})
		})

		Convey("Rectangle", func() {
			r := NewRectangle(20, 40).(*Rectangle)
			Convey("should with width and height", func() {
				So(r.Width, ShouldEqual, 20)
				So(r.Height, ShouldEqual, 40)
			})
			Convey("aabb", func() {
				So(r.AABB(0), ShouldResemble, geom.AABB{HW: 10, HH: 20})
				r.AABB(math.Pi * 2)
				SkipSo(r.AABB(math.Pi*2), ShouldResemble, geom.AABB{Y: -10, HW: 10, HH: 10})
			})
			SkipConvey("FarthestHullPoint", func() {
				So(r.FarthestHullPoint(geom.Vector{1, 0}), ShouldResemble, geom.Vector{10, 20})
				So(r.FarthestHullPoint(geom.Vector{0, 1}), ShouldResemble, geom.Vector{0, 0})
			})
			SkipConvey("FarthestCorePoint", func() {
				So(r.FarthestCorePoint(geom.Vector{1, 0}, 10), ShouldResemble, geom.Vector{10, 0})
				So(r.FarthestCorePoint(geom.Vector{0, 1}, 10), ShouldResemble, geom.Vector{0, 10})
			})
		})
	})

	Convey("geometry utils", t, func() {
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

			So(IsPolygonConvex(nil), ShouldBeFalse)
			So(IsPolygonConvex([]geom.Vector{{}}), ShouldBeTrue)
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

			So(IsPointInPolygon(inside, []geom.Vector{inside}), ShouldBeTrue)
			So(IsPointInPolygon(inside, []geom.Vector{outside}), ShouldBeFalse)

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

			result = NearestPointOnLine(geom.Vector{0, 0}, line1, line1)
			So(result, ShouldResemble, line1)

			result = NearestPointOnLine(geom.Vector{1, -5}, line1, line2)
			So(result, ShouldResemble, line1)
			result = NearestPointOnLine(geom.Vector{2, 2}, line1, line2)
			So(result, ShouldResemble, geom.Vector{1, 3})
			result = NearestPointOnLine(geom.Vector{10, 8}, line1, line2)
			So(result, ShouldResemble, line2)
		})
	})
}
