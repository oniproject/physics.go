package physics

import (
	"github.com/oniproject/physics.go/behaviors"
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

		square := bodies.NewConvexPolygon([]geom.Vector{
			{0, 0},
			{0, 10},
			{10, 10},
			{10, 0},
		})
		square.SetPosition(5, 5)

		sweepPrune := behaviors.NewSweepPrune()
		bodyColl := behaviors.NewBodyCollisionDetection()

		world.Add(sweepPrune, bodyColl)

		world.Step(time.Time{})
		world.Step(time.Time{})

		Convey("should find a collision", func() {
			world.Add(circle, square)

			collide := false
			callback := func(interface{}) { collide = true }

			candidatesLog := func(data interface{}) { Println("!!!!!!!!!! candidatesLog", data) }
			world.On("collisions:candidates", &candidatesLog)

			world.On("collisions:detected", &callback)
			world.Step(time.Time{})
			for _, body := range world.Bodies() {
				Println("body", body.State().Pos)
			}
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
