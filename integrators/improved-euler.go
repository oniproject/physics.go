package integrators

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/util"
	"time"
)

type ImprovedEuler struct {
	Drag float64
}

func NewImprovedEuler() Integrator {
	return &ImprovedEuler{}
}

func (this *ImprovedEuler) SetWorld(world util.EventTarget) {}

func (this *ImprovedEuler) IntegratePositions(things []bodies.Body, dt time.Duration) {
	drag := 1 - this.Drag

	for _, body := range things {
		if body.Treatment() == bodies.TREATMENT_STATIC {
			// set the velocity and acceleration to zero
			body.State().Vel = geom.Vector{}
			body.State().Acc = geom.Vector{}
			body.State().Angular.Vel = 0
			body.State().Angular.Acc = 0
			continue
		}

		// Inspired from https://github.com/soulwire/Coffee-Physics
		// @licence MIT
		//
		// x += (v * dt) + (a * 0.5 * dt * dt)
		// v += a * dt

		// Scale force to mass.
		// state.acc.mult( body.massInv );

		state := body.State()

		state.Old.Vel = state.Vel
		state.Old.Acc = state.Acc

		state.Vel = state.Vel.Plus(state.Acc.Times(dt.Seconds()))

		if drag != 0 {
			state.Vel = state.Vel.Times(drag)
		}

		state.Acc = geom.Vector{}

		state.Old.Angular.Vel = state.Angular.Vel
		state.Angular.Vel += state.Angular.Acc * dt.Seconds()
		state.Angular.Acc = 0
	}
}
func (this *ImprovedEuler) IntegrateVelocities(things []bodies.Body, dt time.Duration) {
	halfdtdt := 0.5 * dt.Seconds() * dt.Seconds()

	for _, body := range things {
		if body.Treatment() == bodies.TREATMENT_STATIC {
			continue
		}

		state := body.State()

		state.Old.Pos = state.Pos
		vel := state.Old.Vel
		state.Pos = state.Pos.Plus(vel.Times(dt.Seconds())).Plus(state.Old.Acc.Times(halfdtdt))

		state.Old.Angular.Pos = state.Angular.Pos
		state.Angular.Pos += state.Old.Angular.Vel*float64(dt) + state.Old.Angular.Acc*halfdtdt
		state.Old.Angular.Acc = 0
	}
}
