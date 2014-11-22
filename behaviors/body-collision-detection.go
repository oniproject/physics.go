package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	//"github.com/oniproject/physics.go/geom"
)

type BodyCollisionDetection struct {
	Check   string // chan to listen to for collision candidates
	Channel string // chan to publish events to

	targets []bodies.Body
	world   World

	checkAllC, checkC func(interface{})
}

func NewBodyCollisionDetection() Behavior {
	b := &BodyCollisionDetection{
		Check:   "collisions:candidates",
		Channel: "collisions:detected",
	}
	b.checkC = func(interface{}) {}
	b.checkAllC = func(interface{}) {}
	return b
}

func (b *BodyCollisionDetection) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *BodyCollisionDetection) Targets() []bodies.Body       { return b.targets }
func (b *BodyCollisionDetection) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		if b.Check == "forse" || b.Check == "" {
			world.Off("integrate:velocities", &b.checkAllC)
		} else {
			world.Off(b.Check, &b.checkC)
		}
	}
	if world != nil {
		// connect
		if b.Check == "forse" || b.Check == "" {
			world.On("integrate:velocities", &b.checkAllC)
		} else {
			world.On(b.Check, &b.checkC)
		}
	}
	b.world = world
}

func (b *BodyCollisionDetection) check()    {}
func (b *BodyCollisionDetection) checkAll() {}
