package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"log"
)

type EdgeCollisionDetecton struct {
	// 1   0.99
	Channel          string
	cof, restitution float64

	min, max geom.Vector

	body bodies.Body

	targets   []bodies.Body
	checkAllC func(interface{})
	world     World
}

func NewEdgeCollisionDetection(edges geom.AABB, cof, restitution float64) Behavior {
	b := &EdgeCollisionDetecton{
		Channel:     "collisions:detected",
		cof:         cof,
		restitution: restitution,
	}
	b.SetAABB(edges)
	b.checkAllC = func(interface{}) { b.checkAll() }
	b.body = bodies.NewPoint()
	b.body.SetRestitution(b.restitution)
	b.body.SetCof(b.cof)
	b.body.SetTreatment(bodies.TREATMENT_STATIC)
	return b
}

func (b *EdgeCollisionDetecton) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *EdgeCollisionDetecton) Targets() []bodies.Body {
	if len(b.targets) == 0 {
		return b.world.Bodies()
	}
	return b.targets
}
func (b *EdgeCollisionDetecton) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		world.Off("integrate:velocities", &b.checkAllC)
	}
	if world != nil {
		// connect
		world.On("integrate:velocities", &b.checkAllC)
	}
	b.world = world
}

func (this *EdgeCollisionDetecton) SetAABB(aabb geom.AABB) {
	this.min = geom.Vector{
		X: aabb.X - aabb.HW,
		Y: aabb.Y - aabb.HH,
	}
	this.max = geom.Vector{
		X: aabb.X + aabb.HW,
		Y: aabb.Y + aabb.HH,
	}
}

func (this *EdgeCollisionDetecton) checkAll() {
	collisions := []Collision{}
	for _, body := range this.Targets() {
		if body.Treatment() == bodies.TREATMENT_DYNAMIC {
			ret := checkEdgeCollide(body, this.min, this.max, this.body)
			collisions = append(collisions, ret...)
		}
	}
	if len(collisions) != 0 {
		log.Print(this.Channel, collisions)
		this.world.Emit(this.Channel, collisions)
	}
}

func checkEdgeCollide(body bodies.Body, min, max geom.Vector, dummy bodies.Body) []Collision {
	return checkGeneral(body, min, max, dummy)
}

func checkGeneral(body bodies.Body, min, max geom.Vector, dummy bodies.Body) []Collision {
	collisions := []Collision{}

	aabb := body.AABB(0)
	trans := geom.NewTransformAngle(body.State().Angular.Pos)
	var overlap float64

	// right
	overlap = (aabb.X + aabb.HW) - max.X
	if overlap >= 0 {
		dir := trans.RotateInv(geom.Vector{X: 1})
		collisions = append(collisions, Collision{
			BodyA: body, BodyB: dummy,
			Overlap: overlap,
			Norm:    geom.Vector{X: 1},
			MTV:     geom.Vector{X: overlap},
			Pos:     trans.Rotate(body.Geometry().FarthestHullPoint(dir)),
		})
	}

	// bottom
	overlap = (aabb.Y + aabb.HH) - max.Y
	if overlap >= 0 {
		dir := trans.RotateInv(geom.Vector{Y: 1})
		collisions = append(collisions, Collision{
			BodyA: body, BodyB: dummy,
			Overlap: overlap,
			Norm:    geom.Vector{Y: 1},
			MTV:     geom.Vector{Y: overlap},
			Pos:     trans.Rotate(body.Geometry().FarthestHullPoint(dir)),
		})
	}

	// left
	overlap = min.X - (aabb.X - aabb.HW)
	if overlap >= 0 {
		dir := trans.RotateInv(geom.Vector{X: -1})
		collisions = append(collisions, Collision{
			BodyA: body, BodyB: dummy,
			Overlap: overlap,
			Norm:    geom.Vector{X: -1},
			MTV:     geom.Vector{X: -overlap},
			Pos:     trans.Rotate(body.Geometry().FarthestHullPoint(dir)),
		})
	}

	// top
	overlap = min.Y - (aabb.Y - aabb.HH)
	if overlap >= 0 {
		dir := trans.RotateInv(geom.Vector{Y: -1})
		collisions = append(collisions, Collision{
			BodyA: body, BodyB: dummy,
			Overlap: overlap,
			Norm:    geom.Vector{Y: -1},
			MTV:     geom.Vector{Y: -overlap},
			Pos:     trans.Rotate(body.Geometry().FarthestHullPoint(dir)),
		})
	}

	return collisions
}
