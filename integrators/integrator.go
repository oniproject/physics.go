package integrators

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/util"
	"time"
)

type World interface {
	util.EventTarget
}

type Integrator interface {
	//Init()
	// Options() Options
	//Connect(world util.EventTarget)
	//Disconnect(world util.EventTarget)
	//Integrate(bodies []Body, dt time.Duration)
	IntegratePositions(bodies []bodies.Body, dt time.Duration)
	IntegrateVelocities(bodies []bodies.Body, dt time.Duration)
	SetWorld(world World)
}

/*func IntegratorSetWorld(integrator Integrator, world World)                                {
	//if
}*/
