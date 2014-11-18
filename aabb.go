package physics

import (
	"github.com/oniproject/physics.go/geom"
	"math"
)

type AABB struct {
	X  float64 // the x coord of the center point
	Y  float64 // the y coord of the center point
	HW float64 // the half-width
	HH float64 // the half-height
}

/* Physics.aabb( minX, minY, maxX, maxY ) -> Object
 * - minX (Number): The x coord of the "top left" point
 * - minY (Number): The y coord of the "top left" point
 * - maxX (Number): The x coord of the "bottom right" point
 * - maxY (Number): The y coord of the "bottom right" point
 */
func NewAABB_byMM(minX, minY, maxX, maxY float64) (aabb AABB) {
	// here, we should have all the arguments as numbers
	aabb.HW = math.Abs(maxX-minX) * 0.5
	aabb.HH = math.Abs(maxY-minY) * 0.5
	aabb.X = (maxX + minX) * 0.5
	aabb.Y = (maxY + minY) * 0.5
	return
}

/** - pt1 (Vectorish): The first corner
 * - pt2 (Vectorish): The opposite corner
 */
func NewAABB_byPoints(pt1, pt2 geom.Vector) (aabb AABB) {
	return NewAABB_byMM(pt1.X, pt1.Y, pt2.X, pt2.Y)
}

func NewAABB_byWH(w, h float64) (aabb AABB) {
	aabb.HW, aabb.HH = w*0.5, h*0.5
	return
}
func NewAABB_byCenter(w, h float64, c geom.Vector) (aabb AABB) {
	aabb.X, aabb.Y = c.X, c.Y
	aabb.HW, aabb.HH = w*0.5, h*0.5
	return
}

/**
 * Physics.aabb.contains( aabb, pt ) -> Boolean
 * - aabb (Object): The aabb
 * - pt (Vector): The point
 * + (Boolean): `true` if `pt` is inside `aabb`, `false` otherwise
 *
 * Check if a point is inside an aabb.
 **/
func AABBcontains(aabb AABB, pt geom.Vector) bool {
	return (pt.X > (aabb.X - aabb.HW)) &&
		(pt.X < (aabb.X + aabb.HW)) &&
		(pt.Y > (aabb.Y - aabb.HH)) &&
		(pt.Y < (aabb.Y + aabb.HH))
}

/**
 * Physics.aabb.clone( aabb ) -> Object
 * - aabb (Object): The aabb to clone
 * + (Object): The clone
 *
 * Clone an aabb.
 **/
/*Physics.aabb.clone = function( aabb ){
    return {
        x: aabb.x,
        y: aabb.y,
        hw: aabb.hw,
        hh: aabb.hh
    };
};*/

/**
 * Physics.aabb.overlap( aabb1, aabb2 ) -> Boolean
 * - aabb1 (Object): The first aabb
 * - aabb2 (Object): The second aabb
 * + (Boolean): `true` if they overlap, `false` otherwise
 *
 * Check if two AABBs overlap.
 **/
func AABBoverlap(aabb1, aabb2 AABB) bool {
	min1 := aabb1.X - aabb1.HW
	min2 := aabb2.X - aabb2.HW
	max1 := aabb1.X + aabb1.HW
	max2 := aabb2.X + aabb2.HW

	// first check x-axis

	if (min2 <= max1 && max1 <= max2) || (min1 <= max2 && max2 <= max1) {
		// overlap in x-axis
		// check y...
		min1 = aabb1.Y - aabb1.HH
		min2 = aabb2.Y - aabb2.HH
		max1 = aabb1.Y + aabb1.HH
		max2 = aabb2.Y + aabb2.HH

		return (min2 <= max1 && max1 <= max2) || (min1 <= max2 && max2 <= max1)
	}

	// they don't overlap
	return false
}
