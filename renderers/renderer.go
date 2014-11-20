package renderers

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geometries"
	"github.com/oniproject/physics.go/util"
)

type Renderer interface {
	//Init()
	CreateView(geometry geometries.Geometry /*, styles*/) interface{}
	DrawBody(body bodies.Body, view interface{})
	DrawMeta(meta interface{})
	Render(bodies []bodies.Body, meta interface{})
	SetWorld(world util.EventTarget)
}
