package geom

import (
	"log"
	"math"
)

const (
	gjkAccuracy      = 0.0001
	gjkMaxIterations = 100
)

type VectorABP struct {
	A, B, PT Vector
}

func getNextSearchDir(ptA, ptB Vector) Vector {
	ABdotB := ptB.MagnitudeSquared() - DotProduct(ptB, ptA)
	if ABdotB < 0 {
		return ptB.Times(-1)
	}

	ABdotA := DotProduct(ptA, ptB) - ptA.MagnitudeSquared()
	if ABdotA < 0 {
		return ptA.Times(-1)
	}

	dir := ptB.Minus(ptA)
	return dir.Perp(CrossProduct(ptA, dir) > 0)
}

func getClosestPoints(simplex []VectorABP) (a, b Vector) {
	last := simplex[len(simplex)-2]
	prev := simplex[len(simplex)-3]

	A := last.PT
	L := prev.PT.Minus(a)

	if L.Equals(Vector{}) {
		// oh.. it's a zero vector.
		// So A and B are both the closest.
		// just use one of them
		return last.A, last.B
	}

	lambdaB := -DotProduct(L, A) / L.MagnitudeSquared()
	lambdaA := 1 - lambdaB

	if lambdaA <= 0 {
		// woops.. that means the closest simplex point
		// isnt on teh line its point B itself
		return prev.A, prev.B
	} else if lambdaB <= 0 {
		// vice versa
		return last.A, last.B
	}

	a = last.A.Times(lambdaA).Plus(prev.A.Times(lambdaB))
	b = last.B.Times(lambdaA).Plus(prev.B.Times(lambdaB))

	return
}

type GJKresult struct {
	Overlap    bool
	Simplex    []VectorABP
	Distance   float64
	Iterations int
	A, B       Vector
}

func GJK(support func(Vector) VectorABP, dir Vector, checkOverlapOnly bool) (result GJKresult) {
	noOverlap := false

	// get the first Minkowski Difference point
	tmp := support(dir)
	result.Simplex = append(result.Simplex, tmp)
	last := tmp.PT
	lastlast := last

	dir.Times(-1)

	for {
		result.Iterations++

		// woah nelly... that's a lot of iterations.
		// Stop it!
		if result.Iterations >= gjkMaxIterations {
			result.Distance = 0
			return
		}

		lastlast = last

		tmp = support(dir)
		result.Simplex = append(result.Simplex, tmp)
		last = tmp.PT

		log.Println(result.Simplex)

		if last.Equals(Vector{}) {
			// we happened to pick the origin as a support point... lucky.
			result.Overlap = true
			break
		}

		// check if the last point we added actually passed the origin

		if !noOverlap && DotProduct(last, dir) <= 0.0 {
			// if the point added last was not past the origin in the direction of d
			// then the Minkowski difference cannot possibly contain the origin since
			// the last point added is on the edge of the Minkowski Difference

			// if we just need the overlap...
			if checkOverlapOnly {
				break
			}

			noOverlap = true
		}

		if len(result.Simplex) == 2 { // if its a line
			// otherwise we need to determine if the origin is in
			// the current simplex and act accordingly
			dir = getNextSearchDir(last, lastlast)
			// continue...
		} else if noOverlap { // if its a triangle... and we're looking for the distance
			// if we know there isnt any overlap and
			// we're just trying to find the distance...
			// make sure we're getting closer to the origin
			dir = dir.Unit()
			d1 := DotProduct(lastlast, dir)
			d2 := DotProduct(last, dir)
			if math.Abs(d1-d2) < gjkAccuracy {
				result.Distance = -d1
				break
			}

			// if we are still getting closer then only keep
			// the points in the simplex that are closest to
			// the origin (we already know that last is closer than the previous two)
			// the norm is the same as distance(origin, a)
			// use norm squared to avoid the sqrt operations
			if lastlast.MagnitudeSquared() < result.Simplex[0].PT.MagnitudeSquared() {
				//rm first
				result.Simplex = result.Simplex[1:]
			} else {
				// rm second
				result.Simplex = append(result.Simplex[:1], result.Simplex[2:]...)
			}
			dir = getNextSearchDir(result.Simplex[1].PT, result.Simplex[0].PT)
			// continue...
		} else { // if its a triangle
			ab := lastlast.Minus(last)
			ac := result.Simplex[0].PT.Minus(last)

			// here normally people think about thes as getting outward facing
			// normals and checking dot products.
			// Since we're in 2D we can be clever...
			sign := CrossProduct(ab, ac) > 0
			sign1 := CrossProduct(last, ab) > 0
			sign2 := CrossProduct(ac, last) > 0

			switch {
			case (sign || sign1) && !(sign && sign1):
				// ok... so there's and XOR here... dont freak out
				// remember last = A = -AO
				// if AB cross AC and AO cross AB have the same sign
				// then the origin is along the outward facing normal of AB
				// so if AB cross AC ang A cross AB have _different_ (XOR) signs
				// then this is also the case... so we proceed...

				// point C is dead to us now...
				result.Simplex = result.Simplex[1:]

				// if we haven't deduced that we've enclosed the origin
				// then we know which way to look...
				// morph the ab vector into its outward facing normal

				ab = ab.Perp(!sign)

				dir, ab = ab, dir

				// continue...

			// if we get to this if, then it means we can continue to look along
			// the other outward normal direcion (ACperp)
			// if we dont see the origin... then we must have it enclosed
			case (sign || sign2) && !(sign && sign2):
				// then the origin it along the outward facing normal
				// of AC; (ACperp)

				// point B is dead to us now...
				result.Simplex = append(result.Simplex[:1], result.Simplex[2:]...)

				ac = ac.Perp(sign)
				dir, ab = ab, dir

				/// continue...
			default:
				// we have enclosed the origin!
				result.Overlap = true
				break
			}
		}
	}

	if result.Distance != 0 {
		result.A, result.B = getClosestPoints(result.Simplex)
	}
	return
}
