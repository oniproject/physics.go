package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"math"
)

type BodyImpulseResponse struct {
	Channel string

	respondC func(interface{})

	targets []bodies.Body
	world   World
}

func NewBodyImpulseResponse() Behavior {
	b := &BodyImpulseResponse{
		Channel: "collisions:detected",
	}
	b.respondC = func(data interface{}) { b.respond(data.([]Collision)) }
	return b
}

func (b *BodyImpulseResponse) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *BodyImpulseResponse) Targets() []bodies.Body       { return b.targets }
func (b *BodyImpulseResponse) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		world.Off(b.Channel, &b.respondC)
	}
	if world != nil {
		// connect
		world.On(b.Channel, &b.respondC)
	}
	b.world = world
}

func (b *BodyImpulseResponse) respond(collisions []Collision) {
	// TODO shuffle
	for _, c := range collisions {
		b.collideBodies(c.BodyA, c.BodyB, c.Norm, c.Pos, c.MTV, false)
	}
}

func (b *BodyImpulseResponse) collideBodies(bodyA, bodyB bodies.Body, norm, point, mtv geom.Vector, contact bool) {
	fixedA := bodyA.Treatment() != bodies.TREATMENT_DYNAMIC
	fixedB := bodyB.Treatment() != bodies.TREATMENT_DYNAMIC

	// do nothing if both are fixed
	if fixedA && fixedB {
		return
	}

	stateA, stateB := bodyA.State(), bodyB.State()

	// extract bodies
	switch {
	case fixedA:
		stateB.Pos = stateB.Pos.Plus(mtv)
	case fixedB:
		stateA.Pos = stateA.Pos.Minus(mtv)
	default:
		mtv = mtv.Times(0.5)
		stateA.Pos = stateA.Pos.Minus(mtv)
		stateB.Pos = stateB.Pos.Plus(mtv)
	}

	// inverse masses and moments of inertia.
	// give fixed bodies infinite mass and moi
	invMoiA, invMoiB, invMassA, invMassB := float64(0), float64(0), float64(0), float64(0)
	if !fixedA {
		invMoiA = 1.0 / bodyA.MOI()
		if math.IsInf(invMoiA, 0) {
			invMoiA = 0
		}
		invMassA = 1.0 / bodyA.Mass()
	}
	if !fixedB {
		invMoiB = 1.0 / bodyB.MOI()
		if math.IsInf(invMoiB, 0) {
			invMoiB = 0
		}
		invMassB = 1.0 / bodyB.Mass()
	}

	// coefficient of restitution between bodies
	var cor float64
	if !contact {
		cor = bodyA.Restitution() * bodyB.Restitution()
	}

	// coefficient of friction between bodies
	cof := bodyA.Cof() * bodyB.Cof()

	// vector perpendicular to n
	perp := norm.Perp(false)

	// collision point from A's center
	rA := point
	// collision point from B's center
	rB := point.Plus(stateA.Pos).Minus(stateB.Pos)

	//tmp

	angVelA := stateA.Angular.Vel
	angVelB := stateB.Angular.Vel

	// relative velocity towards B at collision point
	vAB := stateB.Vel.
		Plus(rB.Perp(false).Times(angVelB)).
		Minus(stateA.Vel).
		Minus(rA.Perp(false).Times(angVelA))

	// break up components along normal and perp-normal directions
	rAproj, rAreg := rA.Proj(norm), rA.Proj(perp)
	rBproj, rBreg := rB.Proj(norm), rB.Proj(perp)
	vproj, vreg := vAB.Proj(norm), vAB.Proj(perp)

	// if moving away from each other... dont' bother.
	if vproj >= 0 {
		return
	}

	impulse := -((1 + cor) * vproj) / (invMassA + invMassB + invMoiA*rAreg*rAreg + invMoiB*rBreg*rBreg)

	// apply impulse
	switch {
	case fixedA:
		norm = norm.Times(impulse * invMassB)
		stateB.Vel = stateB.Vel.Plus(norm)
		stateB.Angular.Vel -= impulse * invMoiB * rBreg
	case fixedB:
		norm = norm.Times(impulse * invMassA)
		stateA.Vel = stateA.Vel.Minus(norm)
		stateA.Angular.Vel += impulse * invMoiA * rAreg
	default:
		norm = norm.Times(impulse * invMassB)
		stateB.Vel = stateB.Vel.Plus(norm)
		stateB.Angular.Vel -= impulse * invMoiB * rBreg
		norm = norm.Times(invMassA * bodyB.Mass())
		stateA.Vel = stateA.Vel.Minus(norm) // XXX
		stateA.Angular.Vel += impulse * invMoiA * rAreg
	}

	//inContact := (impulse < 0.004)
	inContact := false

	// if we have friction and a relative velocity perpendicular to the normal
	if cof != 0 && vreg != 0 {
		// TODO: here, we could first assume static friction applies
		// and that the tangential relative velocity is zero.
		// Then we could calculate the impulse and check if the
		// tangential impulse is less than that allowed by static
		// friction. If not, _then_ apply kinetic friction.

		// instead we're just applying kinetic friction and making
		// sure the impulse we apply is less than the maximum
		// allowed amount

		// maximum impulse allowed by kinetic friction
		max := vreg / (invMassA + invMassB + invMoiA*rAproj*rAproj + invMoiB*rBproj*rBproj)

		if !inContact {
			// the sign of vreg (plus or minus 1)
			sign := 1
			if vreg < 0 {
				sign = -1
			}

			// get impulse due to friction
			impulse *= float64(sign) * cof
			// make sure the impulse isnt giving the system energy
			if sign == 1 {
				impulse = math.Min(impulse, max)
			} else {
				impulse = math.Max(impulse, max)
			}
		} else {
			impulse = max
		}

		// apply frictional impulse
		switch {
		case fixedA:
			perp = perp.Times(impulse * invMassB)
			stateB.Vel = stateB.Vel.Minus(perp)
			stateB.Angular.Vel -= impulse * invMoiB * rBproj
		case fixedB:
			perp = perp.Times(impulse * invMassA)
			stateA.Vel = stateA.Vel.Plus(perp)
			stateA.Angular.Vel += impulse * invMoiA * rAproj
		default:
			perp = perp.Times(impulse * invMassB)
			stateB.Vel = stateB.Vel.Minus(perp)
			stateB.Angular.Vel -= impulse * invMoiB * rBproj
			perp = perp.Times(invMassA * bodyB.Mass())
			stateA.Vel = stateA.Vel.Plus(perp) // XXX
			stateA.Angular.Vel += impulse * invMoiA * rAproj
		}
	}
}
