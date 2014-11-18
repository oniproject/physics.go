package physics

import (
	"github.com/oniproject/physics.go/geom"
	"testing"
)

func It(t *testing.T, str string, fn func() bool) {
	if !fn() {
		t.Error("It", str)
	}
}

func Test_AABB(t *testing.T) {
	var matches = func(t *testing.T, aabb, other AABB) bool {
		if aabb.X != other.X {
			t.Error("fail x", aabb, other)
			return false
		}
		if aabb.Y != other.Y {
			t.Error("fail y", aabb, other)
			return false
		}
		if aabb.HW != other.HW {
			t.Error("fail hw", aabb, other)
			return false
		}
		if aabb.HH != other.HH {
			t.Error("fail hh", aabb, other)
			return false
		}
		return true
	}

	It(t, "should initialize to zero", func() bool {
		aabb := AABB{}
		return matches(t, aabb, AABB{X: 0, Y: 0, HW: 0, HH: 0})
	})

	It(t, "should initialize provided a width/height", func() bool {
		aabb := NewAABB_byWH(4, 5)
		return matches(t, aabb, AABB{X: 0, Y: 0, HW: 2, HH: 2.5})
	})

	It(t, "should initialize provided a width/height and point", func() bool {
		aabb := NewAABB_byCenter(4, 5, geom.Vector{X: 20, Y: 9})
		return matches(t, aabb, AABB{X: 20, Y: 9, HW: 2, HH: 2.5})
	})

	It(t, "should initialize provided two points", func() bool {
		aabb := NewAABB_byPoints(geom.Vector{X: 13, Y: 21}, geom.Vector{X: 20, Y: 9})
		return matches(t, aabb, AABB{X: 16.5, Y: 15, HW: 3.5, HH: 6})
	})

	It(t, "should initialize provided minima and maxima", func() bool {
		aabb := NewAABB_byMM(13, 9, 20, 21)
		return matches(t, aabb, AABB{X: 16.5, Y: 15, HW: 3.5, HH: 6})
	})
}
