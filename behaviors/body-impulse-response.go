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

func (b *BodyImpulseResponse) collideBodies(bodyA, bodyB bodies.Body, n, point, mtv geom.Vector, contact bool) {
	fixedA := bodyA.Treatment() == bodies.TREATMENT_STATIC || bodyA.Treatment() == bodies.TREATMENT_KINEMATIC
	fixedB := bodyB.Treatment() == bodies.TREATMENT_STATIC || bodyB.Treatment() == bodies.TREATMENT_KINEMATIC

	// do nothing if both are fixed
	if fixedA && fixedB {
		return
	}

	// extract bodies
	switch {
	case fixedA:
		bodyB.State().Pos = bodyB.State().Pos.Plus(mtv)
	case fixedB:
		bodyA.State().Pos = bodyA.State().Pos.Minus(mtv)
	default:
		mtv := mtv.Times(0.5)
		bodyA.State().Pos = bodyA.State().Pos.Minus(mtv)
		bodyB.State().Pos = bodyB.State().Pos.Plus(mtv)
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
	perp := n.Perp(false)

	// collision point from A's center
	rA := point
	// collision point from B's center
	rB := point.Plus(bodyA.State().Pos).Minus(bodyB.State().Pos)

	//tmp

	angVelA := bodyA.State().Angular.Vel
	angVelB := bodyB.State().Angular.Vel

	// relative velocity towards B at collision point
	vAB := bodyB.State().Vel.
		Plus(rB.Perp(false).Times(angVelB)).
		Minus(bodyA.State().Vel).
		Minus(rA.Perp(false).Times(angVelA))

	// break up components along normal and perp-normal directions
	rAproj := rA.Proj(n)
	rAreg := rA.Proj(perp)
	rBproj := rB.Proj(n)
	rBreg := rB.Proj(perp)
	vproj := vAB.Proj(n)
	vreg := vAB.Proj(perp)
	// impulse
	// sign
	//max

	// if moving away from each other... dont' bother.
	if vproj >= 0 {
		return
	}

	impulse := -((1 + cor) * vproj) / (invMassA + invMassB + invMoiA*rAreg*rAreg + invMoiB*rBreg*rBreg)

	// apply impulse
	switch {
	case fixedA:
		bodyB.State().Vel = bodyB.State().Vel.Plus(n.Times(impulse * invMassB))
		bodyB.State().Angular.Vel -= impulse * invMoiB * rBreg
	case fixedB:
		bodyA.State().Vel = bodyA.State().Vel.Minus(n.Times(impulse * invMassA))
		bodyA.State().Angular.Vel += impulse * invMoiA * rAreg
	default:
		bodyB.State().Vel = bodyB.State().Vel.Plus(n.Times(impulse * invMassB))
		bodyB.State().Angular.Vel -= impulse * invMoiB * rBreg
		bodyA.State().Vel = bodyA.State().Vel.Minus(n.Times(invMassA * bodyB.Mass())) // XXX
		bodyA.State().Angular.Vel += impulse * invMoiA * rAreg
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
			bodyB.State().Vel = bodyB.State().Vel.Minus(perp.Times(impulse * invMassB))
			bodyB.State().Angular.Vel -= impulse * invMoiB * rBproj
		case fixedB:
			bodyA.State().Vel = bodyA.State().Vel.Plus(perp.Times(impulse * invMassA))
			bodyA.State().Angular.Vel += impulse * invMoiA * rAproj
		default:
			bodyB.State().Vel = bodyB.State().Vel.Minus(perp.Times(impulse * invMassB))
			bodyB.State().Angular.Vel -= impulse * invMoiB * rBproj
			bodyA.State().Vel = bodyA.State().Vel.Plus(perp.Times(invMassA * bodyB.Mass())) // XXX
			bodyA.State().Angular.Vel += impulse * invMoiA * rAproj
		}
	}
}
