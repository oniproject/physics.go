package physics

import (
	//"github.com/oniproject/physics.go/behaviors"
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	//"github.com/oniproject/physics.go/util"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_World(t *testing.T) {
	Convey("World", t, func() {
		world := NewWorldImprovedEuler()

		circle := bodies.NewCircle(20)
		circle.SetPosition(10, 8)
		//return

		square := bodies.NewConvexPolygon([]geom.Vector{
			{0, 0},
			{0, 10},
			{10, 10},
			{10, 0},
		})
		square.SetPosition(5, 5)

		/* TODO
		world.AddBehavior(behaviors.NewSweepPrune())
		world.AddBehavior(behaviors.NewBodyCollisionDetection())
		*/

		world.Step(time.Time{})
		world.Step(time.Time{})

		Convey("should find a collision", func() {
			world.AddBody(circle)
			world.AddBody(square)

			collide := false
			callback := func(interface{}) { collide = true }

			world.On("collisions:detected", &callback)
			world.Step(time.Time{})
			world.Off("collisions:detected", &callback)

			So(collide, ShouldBeTrue)
		})

		Convey("should not find a collision after body removed", func() {
			world.RemoveBody(square)

			collide := false
			callback := func(interface{}) { collide = true }

			world.On("collisions:detected", &callback)
			world.Step(time.Time{})
			world.Off("collisions:detected", &callback)

			So(collide, ShouldBeTrue)
		})
	})
}
