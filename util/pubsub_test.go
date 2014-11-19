package util

import "testing"

func It(t *testing.T, str string, fn func() bool) {
	if !fn() {
		t.Error("It", str)
	}
}

func Test_PubSub(t *testing.T) {
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

	It(t, "should not call that are subscribed to different topics", func() bool {
		ps.Emit("not-my-topic", nil)
		if calledN != 0 || notCalledN != 0 {
			return false
		}
		return true
	})

	It(t, "should call callback with emitted data", func() bool {
		ps.Emit(topic, data)
		if data != result.(int) {
			t.Error("fail data", data, result)
			return false
		}
		if calledN != 1 || notCalledN != 0 {
			t.Error("fail callN", calledN, notCalledN)
			return false
		}
		return true
	})
}
