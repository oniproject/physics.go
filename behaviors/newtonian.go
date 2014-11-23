package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"math"
)

type Newtonian struct {
	Strength     float64
	maxSq, minSq float64

	behaveC func(interface{})

	targets []bodies.Body
	world   World
}

func NewNewtonian(strength float64) Behavior {
	b := &Newtonian{
		Strength: strength,
		maxSq:    math.Inf(1),
		minSq:    100.0 * strength,
	}
	b.behaveC = func(interface{}) { b.behave() }
	return b
}

func (b *Newtonian) ApplyTo(bodies []bodies.Body) {
	b.targets = bodies
}
func (b *Newtonian) Targets() []bodies.Body {
	if len(b.targets) == 0 {
		return b.world.Bodies()
	}
	return b.targets
}
func (b *Newtonian) SetWorld(world World) {
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

func (b *Newtonian) SetMinMax(min, max float64) {
	b.minSq = min * min
	b.maxSq = max * max
}

func (b *Newtonian) behave() {
	targets := b.Targets()
	for j, body := range targets {
		for i := j + 1; i < len(targets); i++ {
			other := targets[i]
			// clone position
			pos := other.State().Pos.Minus(body.State().Pos)
			normsq := pos.MagnitudeSquared()

			if normsq > b.minSq && normsq < b.maxSq {
				g := b.Strength / normsq
				pos = pos.Unit().Times(g * other.Mass())
				body.Accelerate(pos)
				other.Accelerate(pos.Times(-body.Mass() / other.Mass()))
			}
		}
	}
}
