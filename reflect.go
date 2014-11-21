package physics

import (
	"github.com/oniproject/physics.go/behaviors"
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/integrators"
	"github.com/oniproject/physics.go/renderers"
	"reflect"
)

var typeBehavior = reflect.TypeOf((*behaviors.Behavior)(nil)).Elem()
var typeBody = reflect.TypeOf((*bodies.Body)(nil)).Elem()
var typeIntegrator = reflect.TypeOf((*integrators.Integrator)(nil)).Elem()
var typeRenderer = reflect.TypeOf((*renderers.Renderer)(nil)).Elem()

//var typeGeometry = reflect.TypeOf((*geometries.Geometry)(nil)).Elem()
//func IsGeometry(i interface{}) bool   { return reflect.TypeOf(i).Implements(typeGeometry) }

func IsBehavior(i interface{}) bool   { return reflect.TypeOf(i).Implements(typeBehavior) }
func IsBody(i interface{}) bool       { return reflect.TypeOf(i).Implements(typeBody) }
func IsIntegrator(i interface{}) bool { return reflect.TypeOf(i).Implements(typeIntegrator) }
func IsRenderer(i interface{}) bool   { return reflect.TypeOf(i).Implements(typeRenderer) }

func (w *world) Add(i interface{}) {
	switch {
	case IsBehavior(i):
		w.AddBehavior(i.(behaviors.Behavior))
	case IsBody(i):
		w.AddBody(i.(bodies.Body))
	default:
		panic("fail type")
	}
}

func (w *world) Remove(i interface{}) {
	switch {
	case IsBehavior(i):
		w.RemoveBehavior(i.(behaviors.Behavior))
	case IsBody(i):
		w.RemoveBody(i.(bodies.Body))
	default:
		panic("fail type")
	}
}
