package bodies

import (
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
	"math"
	"time"
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

	asleep              bool
	sleepAngPosMean     float64
	sleepAngPosVariance float64
	sleepPosMean        geom.Vector
	sleepPosVariance    geom.Vector
	sleepMeanK          float64

	sleepIdleTime time.Duration

	SleepTimeLimit     time.Duration
	SleepSpeedLimit    float64
	SleepVarianceLimit float64

	world World
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

func (p *Point) MOI() float64 { return p.moi }

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

func (p *Point) SetWorld(world World) {
	if p.world != nil {
		// disconnect
	}
	if world != nil {
		// connect
	}
	p.world = world
}

func (p *Point) IsSleep() bool { return p.asleep }
func (p *Point) Sleep()        { p.asleep = true }
func (p *Point) WakeUp() {
	p.asleep = false
	p.sleepAngPosMean = 0
	p.sleepAngPosVariance = 0
	p.sleepPosMean = geom.Vector{}
	p.sleepPosVariance = geom.Vector{}
	p.sleepMeanK = 0
}
func (p *Point) SleepCheck(dt time.Duration) {
	aabb := p.geometry.AABB(0)
	r := math.Max(aabb.HW, aabb.HH)

	if p.asleep {
		// check vel
		v := p.state.Vel.Magnitude() + math.Abs(r*p.state.Angular.Vel)
		limit := p.SleepSpeedLimit
		if limit == 0 && p.world != nil {
			limit = p.world.SleepSpeedLimit()
		}

		if v >= limit {
			p.WakeUp()
			return
		}
	}

	p.sleepMeanK++

	p.sleepPosMean, p.sleepPosVariance =
		p.pushRunningVectorAvg(p.sleepMeanK,
			p.sleepPosMean,
			p.sleepPosVariance,
			p.state.Pos)
	p.sleepAngPosMean, p.sleepAngPosVariance =
		p.pushRunningAvg(p.sleepMeanK,
			p.sleepAngPosMean,
			p.sleepAngPosVariance,
			p.state.Angular.Pos)

	v := p.sleepPosVariance.Magnitude() + math.Abs(r*p.sleepAngPosVariance)
	limit := p.SleepVarianceLimit
	if limit == 0 && p.world != nil {
		limit = p.world.SleepVarianceLimit()
	}

	if v <= limit {
		// check idle time
		limit := p.SleepTimeLimit
		if limit == 0 && p.world != nil {
			limit = p.world.SleepTimeLimit()
		}

		p.sleepIdleTime += dt

		if p.sleepIdleTime > limit {
			p.Sleep()
		}
	} else {
		p.WakeUp()
	}

}

// Running average
// http://www.johndcook.com/blog/standard_deviation
// k is num elements
// m is current mean
// s is current std deviation
// v is value to push
func (p *Point) pushRunningAvg(k, m, s, v float64) (float64, float64) {
	x := v - m
	// Mk = Mk-1+ (xk – Mk-1)/k
	// Sk = Sk-1 + (xk – Mk-1)*(xk – Mk).
	m += x / k
	s += x * (v - m)
	return m, s
}

// Running vector average
// http://www.johndcook.com/blog/standard_deviation
// k is num elements
// m is current mean (vector)
// s is current std deviation (vector)
// v is vector to push
func (p *Point) pushRunningVectorAvg(k float64, m, s, v geom.Vector) (geom.Vector, geom.Vector) {
	invK := 1 / k
	x := v.X - m.X
	y := v.Y - m.Y

	// Mk = Mk-1+ (xk – Mk-1)/k
	// Sk = Sk-1 + (xk – Mk-1)*(xk – Mk).
	m = m.Plus(geom.Vector{x * invK, y * invK})
	x *= v.X - m.Y
	y *= v.Y - m.Y
	s = s.Plus(geom.Vector{x, y})
	return m, s
}
