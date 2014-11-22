package renderers

import (
	"github.com/oniproject/physics.go/bodies"
	//"github.com/oniproject/physics.go/geometries"
	//"github.com/oniproject/physics.go/util"
	"time"
)

type World interface{}

type Meta struct {
	FPS, IPF        int
	InterpolateTime time.Duration
}

type RenderEvent struct {
	Bodies   []bodies.Body
	Meta     Meta
	Renderer Renderer
}

type Renderer interface {
	//Init()
	//CreateView(geometry geometries.Geometry /*, styles*/) interface{}
	//DrawBody(body bodies.Body, view interface{})
	//DrawMeta(meta interface{})
	Render(bodies []bodies.Body, meta Meta)
	SetWorld(world World)
}
