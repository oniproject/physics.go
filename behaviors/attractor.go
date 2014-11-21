package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"math"
)

type Attractor struct {
	Pos      geom.Vector
	Strength float64
	Order    float64
	Max, Min float64

	targets []bodies.Body

	behaveC func(interface{})
	world   World
}

func NewAttractor() Behavior {
	b := &Attractor{
		Pos:      geom.Vector{0, 0},
		Strength: 1,
		Order:    2,
		Max:      math.Inf(+1),
		Min:      10,
	}

	b.behaveC = func(interface{}) { b.behave() }

	return b
}

func (b *Attractor) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *Attractor) Targets() []bodies.Body       { return b.targets }

func (b *Attractor) SetWorld(world World) {
	if b.world != nil {
		world.On("integrate:positions", &b.behaveC)
		b.world = nil
	}
	if world != nil {
		world.Off("integrate:positions", &b.behaveC)
	}
}

func (b *Attractor) behave() {
	for _, body := range b.Targets() {
		acc := b.Pos.Minus(body.State().Pos)
		norm := acc.Magnitude()

		if norm > b.Min && norm < b.Max {
			g := b.Strength / math.Pow(norm, b.Order)

			body.Accelerate(acc.Unit().Times(g))
		}
	}
}
