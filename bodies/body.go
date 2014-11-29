package bodies

import (
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
	"time"
)

type World interface {
	SleepVarianceLimit() float64
	SleepSpeedLimit() float64
	SleepTimeLimit() time.Duration
}

type Angular struct {
	Pos, Vel, Acc float64
}

type BodyState struct {
	Pos, Vel, Acc geom.Vector
	Angular       Angular
	Old           *BodyState
}

const (
	TREATMENT_DYNAMIC   = 1
	TREATMENT_KINEMATIC = 2
	TREATMENT_STATIC    = 3
)

type Body interface {
	AABB(angle float64) geom.AABB
	Accelerate(acc geom.Vector)
	ApplyForce(force, p geom.Vector)

	Cof() float64
	SetCof(float64)

	Geometry() geometries.Geometry
	Hidden() bool
	SetHidden(bool)

	Mass() float64
	SetMass(float64)
	// init
	//Options() Options

	Recalc()

	Restitution() float64
	SetRestitution(float64)

	State() *BodyState

	//Stype for render
	Treatment() uint
	SetTreatment(uint)

	UID() int64
	View() interface{}
	SetView(interface{})

	SetPosition(float64, float64)
	SetVelocity(float64, float64)

	MOI() float64

	SetWorld(World)

	IsSleep() bool
	Sleep()
	WakeUp()
	SleepCheck(dt time.Duration)
}
