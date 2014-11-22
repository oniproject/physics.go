package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
)

type ConstantAcceleration struct {
	// default {x: 0, y: 0.0004
	Acc geom.Vector

	behaveC func(interface{})

	targets []bodies.Body
	world   World
}

func NewConstantAcceleration(x, y float64) Behavior {
	b := &ConstantAcceleration{Acc: geom.Vector{X: x, Y: y}}
	b.behaveC = func(interface{}) { b.behave() }

	return b
}

func (b *ConstantAcceleration) behave() {
	for _, body := range b.Targets() {
		body.Accelerate(b.Acc)
	}
}

func (b *ConstantAcceleration) ApplyTo(bodies []bodies.Body) {
	b.targets = bodies
}
func (b *ConstantAcceleration) Targets() []bodies.Body {
	if len(b.targets) == 0 {
		return b.world.Bodies()
	}
	return b.targets
}
func (b *ConstantAcceleration) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		world.Off("integrate:positions", &b.behaveC)
	}
	if world != nil {
		// connect
		world.On("integrate:positions", &b.behaveC)
	}
	b.world = world
}
