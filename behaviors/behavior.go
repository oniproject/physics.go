package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/util"
)

type World interface {
	util.EventTarget
	Bodies() []bodies.Body
}

type Behavior interface {
	ApplyTo([]bodies.Body)
	//Behave(data XX)
	Targets() []bodies.Body
	//Init(options Options)
	//Options() Options
	SetWorld(world World)
}

type behavior struct {
	targets []bodies.Body
	world   util.EventTarget
	behave  *func(interface{})
}

func (b *behavior) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *behavior) Targets() []bodies.Body       { return b.targets }

func (b *behavior) SetWorld(world util.EventTarget) {
	if b.world != nil {
		b.Disconnect(b.world)
		b.world = nil
	}
	if world != nil {
		b.Connect(b.world)
	}
}

func (b *behavior) Connect(world util.EventTarget) {
	if b.behave != nil {
		world.On("integrate:positions", b.behave)
	}
}
func (b *behavior) Disconnect(world util.EventTarget) {
	if b.behave != nil {
		world.Off("integrate:positions", b.behave)
	}
}
