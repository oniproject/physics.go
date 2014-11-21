package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/util"
)

type Collision struct {
	BodyA, BodyB   bodies.Body
	Overlap        float64
	Norm, MTV, Pos geom.Vector
}

type EdgeCollisionDetecton struct {
	// 1   0.99
	Cof, Restitution float64
	f                *func(interface{})

	edges geom.AABB
	body  bodies.Body
	world util.EventTarget
}

func (this *EdgeCollisionDetecton) ApplyTo([]bodies.Body) {
	panic("not implemented")
}

//Behave(data XX)
func (this *EdgeCollisionDetecton) Connect(world util.EventTarget) {
	f := func(interface{}) { this.checkAll() }
	this.f = &f
	world.On("integrate:velocities", this.f)
}
func (this *EdgeCollisionDetecton) Disconnect(world util.EventTarget) {
	world.Off("integrate:velocities", this.f)
}
func (this *EdgeCollisionDetecton) Targets() []bodies.Body {
	panic("not implemented")
}

func (this *EdgeCollisionDetecton) checkAll() {
	collisions := []Collision{}
	for _, body := range this.Targets() {
		if body.Treatment() == bodies.TREATMENT_DYNAMIC {
			ret := checkEdgeCollide(body, this.edges, this.body)
			collisions = append(collisions, ret...)
		}
	}
	if len(collisions) != 0 {
		//this.world.Emit(this.options.cannel
		this.world.Emit("collisions:detected", collisions)
	}
}

//Init(options Options)
//Options() Options
func (this *EdgeCollisionDetecton) SetWorld(world util.EventTarget) {
	panic("not implemented")
}

func checkEdgeCollide(body bodies.Body, bounds geom.AABB, dummy bodies.Body) []Collision {
	return checkGeneral(body, bounds, dummy)
}

func checkGeneral(body bodies.Body, bounds geom.AABB, dummy bodies.Body) []Collision {
	collisions := []Collision{}

	/*

		aabb := body.AABB(0)
		var overlap float64

			// right
			overlap = (aabb.X + aabb.HW) - bounds.MaxX
			if overlap >= 0 {
				dir.set(1, 0).rotateInv(trans.SetRotation(body.State().Angular.Pos))
				collisions = append(collisions, &Collision{
					BodyA: body, BodyB: dummy,
					Overlap: overlap,
					Norm:    geom.Vector{X: 1},
					MTV:     geom.Vector{X: overlap},
					Pos:     body.Geometry().FarthestHullPoint(dir).Rotate(trans).Values(),
				})
			}

			// bottom
			overlap = (aabb.Y + aabb.HH) - bounds.MaxY
			if overlap >= 0 {
				dir.set(0, 1).rotateInv(trans.SetRotation(body.State().Angular.Pos))
				collisions = append(collisions, &Collision{
					BodyA: body, BodyB: dummy,
					Overlap: overlap,
					Norm:    geom.Vector{Y: 1},
					MTV:     geom.Vector{Y: overlap},
					Pos:     body.Geometry().FarthestHullPoint(dir).Rotate(trans).Values(),
				})
			}

			// left
			overlap = bounds.MinX - (aabb.X - aabb.HW)
			if overlap >= 0 {
				dir.set(-1, 0).rotateInv(trans.SetRotation(body.State().Angular.Pos))
				collisions = append(collisions, &Collision{
					BodyA: body, BodyB: dummy,
					Overlap: overlap,
					Norm:    geom.Vector{X: -1},
					MTV:     geom.Vector{X: -overlap},
					Pos:     body.Geometry().FarthestHullPoint(dir).Rotate(trans).Values(),
				})
			}

			// top
			overlap = bounds.MinY - (aabb.Y - aabb.HY)
			if overlap >= 0 {
				dir.set(0, -1).rotateInv(trans.SetRotation(body.State().Angular.Pos))
				collisions = append(collisions, &Collision{
					BodyA: body, BodyB: dummy,
					Overlap: overlap,
					Norm:    geom.Vector{Y: -1},
					MTV:     geom.Vector{Y: -overlap},
					Pos:     body.Geometry().FarthestHullPoint(dir).Rotate(trans).Values(),
				})
			}

	*/

	return collisions
}
