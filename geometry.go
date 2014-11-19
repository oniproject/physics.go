package physics

import (
	"github.com/oniproject/physics.go/geom"
	"math"
)

type Geometry interface {
	//Init(Options)
	//Options() Options
	AABB(angle float64) AABB
	FarthestCorePoint(dir geom.Vector, margin float64) geom.Vector
	FarthestHullPoint(dir geom.Vector) geom.Vector
}

func IsPolygonConvex(hull []geom.Vector) bool {
	if hull == nil || len(hull) == 0 {
		return false
	}

	if len(hull) < 3 {
		// it must be a point or a line...
		// which are convex
		return true
	}

	prev := hull[0].Minus(hull[len(hull)-1])

	sign := false
	signV := 0.0
	for i := 1; i < len(hull); i++ {
		next := hull[i%len(hull)].Minus(hull[(i-1)%len(hull)])

		if !sign {
			sign = true
			signV = geom.CrossProduct(prev, next)
		} else {

			a := signV
			b := geom.CrossProduct(prev, next)

			if a < 0 && b > 0 || a > 0 && b < 0 {
				// if the cross products are different signs it's not convex
				return false
			}
		}

		prev = next
	}
	return true
}

func PolygonMOI(hull []geom.Vector) float64 {
	if hull == nil || len(hull) < 2 {
		// it must be a point
		// moi = 0
		return 0
	}

	if len(hull) == 2 {
		// it's a line
		// gets length squared
		return hull[1].DistanceFromSquared(hull[0]) / 12.0
	}
	//panic("not implemented")

	num, denom := 0.0, 0.0
	prev, next := hull[0], hull[1]
	for i := 1; i < len(hull); i++ {
		next = hull[i]
		tmp := math.Abs(geom.CrossProduct(next, prev))
		nsNext := next.MagnitudeSquared()
		nsPrev := prev.MagnitudeSquared()
		num += tmp * (nsNext + geom.DotProduct(prev, next) + nsPrev)
		denom += tmp
		prev = next
	}

	return num / (6.0 * denom)
}

func IsPointInPolygon(pt geom.Vector, hull []geom.Vector) bool {
	if len(hull) < 2 {
		// it's a point...
		return pt.Equals(hull[0])
	}

	if len(hull) == 2 {
		// it's a line
		ang := pt.Angle(&hull[0])
		ang += pt.Angle(&hull[1])
		return math.Abs(ang) == math.Pi
	}

	prev := hull[0].Minus(pt)

	ang := 0.0
	for i := 1; i <= len(hull); i++ {
		next := hull[i%len(hull)].Minus(pt)
		ang += next.Angle(&prev)
		prev = next
	}

	return math.Abs(ang) > 1e-6
}

func PolygonArea(hull []geom.Vector) float64 {
	if len(hull) < 3 {
		// it must be a point or line
		// area = 0
		return 0
	}

	ret := 0.0
	prev := hull[len(hull)-1]
	for _, next := range hull {
		ret += geom.CrossProduct(next, prev)
		prev = next
	}
	return ret / 2.0
}

func PolygonCentroid(hull []geom.Vector) (ret geom.Vector) {
	if len(hull) < 2 {
		// it must be a point
		return hull[0]
	}

	if len(hull) == 2 {
		// it's a line
		// get the midpoint
		x := (hull[1].X + hull[0].X) / 2
		y := (hull[1].Y + hull[0].Y) / 2
		return geom.Vector{x, y}
	}

	prev := hull[len(hull)-1]

	for _, next := range hull {
		tmp := geom.CrossProduct(next, prev)
		n := prev.Plus(next).Times(tmp)
		ret.Vadd(&n)
		prev = next
	}

	tmp := 1.0 / (6.0 * PolygonArea(hull))

	return ret.Times(tmp)
}

func NearestPointOnLine(pt, linePt1, linePt2 geom.Vector) geom.Vector {
	A := linePt1.Minus(pt)
	L := linePt2.Minus(pt).Minus(A)

	if L.Equals(geom.Vector{0, 0}) {
		return linePt1
	}

	lamdaB := -geom.DotProduct(L, A) / L.MagnitudeSquared()
	lamdaA := 1 - lamdaB

	if lamdaA <= 0 {
		return linePt2
	} else if lamdaB <= 0 {
		return linePt1
	}

	return linePt2.Times(lamdaB).Plus(linePt1.Times(lamdaA))
}
