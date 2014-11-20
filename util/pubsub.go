package util

type EventTarget interface {
	On(name string, fn *func(interface{}))
	Off(name string, fn *func(interface{})) bool
	Emit(name string, data interface{})
}

type PubSub struct {
	callbacks map[string][]*func(interface{})
}

func NewPubSub() PubSub {
	return PubSub{make(map[string][]*func(interface{}))}
}

func (ps *PubSub) On(name string, fn *func(interface{})) {
	ps.Off(name, fn)

	/*if _, ok := ps.callbacks[name]; !ok {
		ps.callbacks[name] = []*func(interface{}){}
	}*/

	ps.callbacks[name] = append(ps.callbacks[name], fn)
}

func (ps *PubSub) Off(name string, fn *func(interface{})) (found bool) {
	if fn == nil {
		delete(ps.callbacks, name)
		return
	}

	arr := ps.callbacks[name]

	for i, cb := range arr {
		if cb == fn {
			arr = append(arr[:i], arr[i+1:]...)
			found = true
		}
	}

	ps.callbacks[name] = arr
	return
}

func (ps *PubSub) Emit(name string, data interface{}) {
	for _, fn := range ps.callbacks[name] {
		(*fn)(data)
	}
}
