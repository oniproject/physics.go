package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_PubSub(t *testing.T) {
	Convey("PubSub", t, func() {
		topic := "my-tomtopic"

		calledN, notCalledN := 0, 0
		data := int(5432)
		var result interface{}

		called := func(d interface{}) {
			calledN++
			result = d
		}
		notCalled := func(interface{}) {
			notCalledN++
		}

		ps := NewPubSub()
		ps.On(topic, &called)
		ps.On(topic, &notCalled)
		ps.Off(topic, &notCalled)

		Convey("should not call that are subscribed to different topics", func() {
			ps.Emit("not-my-topic", nil)

			So(calledN, ShouldEqual, 0)
			So(notCalledN, ShouldEqual, 0)
		})

		Convey("should call callback with emitted data", func() {
			ps.Emit(topic, data)

			So(calledN, ShouldEqual, 1)
			So(notCalledN, ShouldEqual, 0)

			So(data, ShouldEqual, result.(int))
		})

		Convey("should remove all callbacks", func() {
			ps.Off(topic, nil)

			So(ps.callbacks[topic], ShouldBeNil)
		})
	})
}
