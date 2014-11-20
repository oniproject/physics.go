package geom

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_AABB(t *testing.T) {
	Convey("test AABB", t, func() {
		Convey("should initialize to zero", func() {
			aabb := AABB{}
			So(aabb, ShouldResemble, AABB{X: 0, Y: 0, HW: 0, HH: 0})
		})

		Convey("should initialize provided a width/height", func() {
			aabb := NewAABB_byWH(4, 5)
			So(aabb, ShouldResemble, AABB{X: 0, Y: 0, HW: 2, HH: 2.5})
		})

		Convey("should initialize provided a width/height and point", func() {
			aabb := NewAABB_byCenter(4, 5, Vector{X: 20, Y: 9})
			So(aabb, ShouldResemble, AABB{X: 20, Y: 9, HW: 2, HH: 2.5})
		})

		Convey("should initialize provided two points", func() {
			aabb := NewAABB_byPoints(Vector{X: 13, Y: 21}, Vector{X: 20, Y: 9})
			So(aabb, ShouldResemble, AABB{X: 16.5, Y: 15, HW: 3.5, HH: 6})
		})

		Convey("should initialize provided minima and maxima", func() {
			aabb := NewAABB_byMM(13, 9, 20, 21)
			So(aabb, ShouldResemble, AABB{X: 16.5, Y: 15, HW: 3.5, HH: 6})
		})
	})
}
