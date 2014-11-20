package bodies

import (
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
)

type Point struct {
	hidden      bool
	treatment   uint
	mass        float64
	restitution float64
	cof         float64
	view        interface{}

	moi float64

	state *BodyState
	uid   int64

	geometry geometries.Geometry
}

var uidGen int64 = 0

//func NewPoint(treatment uint, x, y, angle, mass float64) *Point {
func NewPoint() (p *Point) {
	uidGen++
	p = &Point{
		uid:         uidGen,
		hidden:      false,
		treatment:   TREATMENT_DYNAMIC,
		mass:        1.0,
		restitution: 1.0,
		cof:         0.8,
		geometry:    geometries.NewPoint(),
		view:        nil,
	}
	p.state = &BodyState{
		Pos: geom.Vector{},
		Vel: geom.Vector{},
		Acc: geom.Vector{},
		Angular: Angular{
			Pos: 0,
			Vel: 0,
			Acc: 0,
		},
		Old: &BodyState{},
	}

	if p.mass == 0 {
		panic("Error: Bodies must have non-zero mass")
	}

	return
}

func (p *Point) Geometry() geometries.Geometry { return p.geometry }

func (p *Point) Cof() float64             { return p.cof }
func (p *Point) SetCof(v float64)         { p.cof = v }
func (p *Point) Hidden() bool             { return p.hidden }
func (p *Point) SetHidden(v bool)         { p.hidden = v }
func (p *Point) Mass() float64            { return p.mass }
func (p *Point) SetMass(v float64)        { p.mass = v }
func (p *Point) Restitution() float64     { return p.restitution }
func (p *Point) SetRestitution(v float64) { p.restitution = v }
func (p *Point) Treatment() uint          { return p.treatment }
func (p *Point) SetTreatment(v uint)      { p.treatment = v }
func (p *Point) View() interface{}        { return p.view }
func (p *Point) SetView(v interface{})    { p.view = v }

func (p *Point) SetPosition(x, y float64) { p.state.Pos = geom.Vector{x, y} }
func (p *Point) SetVelocity(x, y float64) { p.state.Vel = geom.Vector{x, y} }

func (p *Point) State() *BodyState { return p.state }
func (p *Point) UID() int64        { return p.uid }

func (p *Point) Accelerate(acc geom.Vector) {
	if p.treatment == TREATMENT_DYNAMIC {
		p.state.Acc = p.state.Acc.Plus(acc)
	}
}

func (p *Point) ApplyForce(force, pp geom.Vector) {
	if p.treatment != TREATMENT_DYNAMIC {
		return
	}

	if /*pp &&*/ p.moi != 0 {
		p.state.Angular.Acc -= geom.CrossProduct(pp, force) / p.moi
	}

	p.Accelerate(force.Times(1.0 / p.mass))
}

func (p *Point) AABB(angle float64) (aabb geom.AABB) {
	aabb = p.geometry.AABB(angle)
	aabb.X += p.state.Pos.X
	aabb.Y += p.state.Pos.Y
	return
}

func (p *Point) Recalc() {}
