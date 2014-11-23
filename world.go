package physics

import (
	"github.com/oniproject/physics.go/behaviors"
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/integrators"
	"github.com/oniproject/physics.go/renderers"
	"github.com/oniproject/physics.go/util"
	"log"
	"time"
)

type World interface {
	Add(things ...interface{})
	Remove(things ...interface{})

	AddBehavior(behaviors.Behavior)
	AddBody(bodies.Body)
	RemoveBehavior(behaviors.Behavior)
	RemoveBody(bodies.Body)

	Renderer() renderers.Renderer
	SetRenderer(renderers.Renderer)

	Integrator() integrators.Integrator
	SetIntegrator(integrators.Integrator)

	// destroy
	// init

	//Find(query Query) []interface{}
	//FindOne(query Query) interface{}

	Behaviors() []behaviors.Behavior
	Bodies() []bodies.Body

	//Has(thing interface{}) bool

	IsPaused() bool
	Pause()
	Unpause()

	// options

	Itertate(dt time.Duration)
	Render()

	Step(now time.Time)
	TimeStep() time.Duration
	SetTimeStep(now time.Duration)

	//Warp()

	util.EventTarget
}

/*func NewWorldDefault() World {
	return NewWorld(DefaultTimestep, DefaultMaxIPF, DefaultIntegrator)
}*/
func NewWorldImprovedEuler() (w World) {
	w = &world{
		maxIPF: 16,
		PubSub: util.NewPubSub(),

		warp: 1,
	}
	w.SetIntegrator(integrators.NewImprovedEuler())
	w.SetTimeStep(time.Second / 120)
	log.Println("init world", w)
	return
}

type world struct {
	maxIPF int

	meta renderers.Meta

	bodies     []bodies.Body
	behaviors  []behaviors.Behavior
	integrator integrators.Integrator
	renderer   renderers.Renderer

	paused bool
	//warp     time.Duration
	warp     float64
	time     time.Time
	lastTime time.Time

	animTime time.Time

	dt      time.Duration
	maxJump time.Duration

	util.PubSub
}

func (w *world) Integrator() integrators.Integrator { return w.integrator }
func (w *world) SetIntegrator(integrator integrators.Integrator) {
	if integrator == w.integrator {
		return
	}

	if w.integrator != nil {
		w.integrator.SetWorld(nil)
		w.Emit("remove:integrator", w.integrator)
		w.integrator = nil
	}

	if integrator != nil {
		w.integrator = integrator
		w.integrator.SetWorld(w)
		w.Emit("add:integrator", w.integrator)
	}
}

func (w *world) Renderer() renderers.Renderer { return w.renderer }
func (w *world) SetRenderer(renderer renderers.Renderer) {
	if renderer == w.renderer {
		return
	}

	if w.renderer != nil {
		w.renderer.SetWorld(nil)
		w.Emit("remove:renderer", w.renderer)
		w.renderer = nil
	}

	if renderer != nil {
		w.renderer = renderer
		w.renderer.SetWorld(w)
		w.Emit("add:renderer", w.renderer)
	}
}

func (w *world) TimeStep() time.Duration { return w.dt }
func (w *world) SetTimeStep(dt time.Duration) {
	if dt != 0 {
		w.dt = dt
		w.maxJump = dt * time.Duration(w.maxIPF)
	}
}

func (w *world) Behaviors() []behaviors.Behavior { return w.behaviors }
func (w *world) AddBehavior(behavior behaviors.Behavior) {
	w.RemoveBehavior(behavior)
	behavior.SetWorld(w)
	w.behaviors = append(w.behaviors, behavior)
	w.Emit("add:behavior", behavior)
}
func (w *world) RemoveBehavior(behavior behaviors.Behavior) {
	for i, b := range w.behaviors {
		if b == behavior {
			w.behaviors = append(w.behaviors[:i], w.behaviors[i+1:]...)
			w.Emit("remove:behavior", behavior)
			return
		}
	}
}

func (w *world) Bodies() []bodies.Body { return w.bodies }
func (w *world) AddBody(body bodies.Body) {
	w.RemoveBody(body)
	body.Recalc()
	//body.SetWorld(w)
	w.bodies = append(w.bodies, body)
	w.Emit("add:body", body)
}
func (w *world) RemoveBody(body bodies.Body) {
	for i, b := range w.bodies {
		if b == body {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			w.Emit("remove:body", body)
			return
		}
	}
}

type IntegrateEvent struct {
	Bodies []bodies.Body
	Dt     time.Duration
}

func (w *world) Itertate(dt time.Duration) {
	w.integrator.IntegrateVelocities(w.bodies, dt)
	w.Emit("integrate:velocities", IntegrateEvent{w.bodies, dt})
	w.integrator.IntegratePositions(w.bodies, dt)
	w.Emit("integrate:positions", IntegrateEvent{w.bodies, dt})
}

func (w *world) Render() {
	w.renderer.Render(w.bodies, w.meta)
	w.Emit("render", renderers.RenderEvent{w.bodies, w.meta, w.renderer})
}

func (w *world) IsPaused() bool { return w.paused }
func (w *world) Pause() {
	w.paused = true
	w.Emit("paused", nil)
}
func (w *world) Unpause() {
	w.paused = false
	w.Emit("unpaused", nil)
}

func (w *world) Step(now time.Time) {
	// if it's paused, don't step
	// or if it's the first step...
	if w.paused || w.animTime.IsZero() {
		switch {
		case !now.IsZero():
			w.animTime = now
		case w.animTime.IsZero():
			w.animTime = time.Now()
		}
		if !w.paused {
			w.Emit("step", w.meta)
		}
		return
	}

	// FIXME
	invWarp := 1.0 / w.warp
	//invWarp := 1.0 / float64(w.warp)

	animDt := time.Duration(float64(w.dt) * invWarp)
	// new time is specified, or just oni iteration ahead
	if now.IsZero() {
		now = w.animTime.Add(animDt)
	}
	// if the time between this step and the last
	animDiff := now.Sub(w.animTime)

	animMaxJump := time.Duration(float64(w.maxJump) * invWarp)

	// if the time difference is to big... ajust
	if animDiff > animMaxJump {
		w.animTime = now.Add(-animMaxJump)
		animDiff = animMaxJump
	}

	// the "world" time between this step and the last. Adjust for warp
	worldDiff := time.Duration(float64(animDiff) * w.warp)
	//worldDiff := animDiff * w.warp

	// the target time for the world time to step to
	target := w.time.Add(worldDiff - w.dt)

	//for w.time <= target {
	for w.time.Sub(target) <= 0 {
		w.time = w.time.Add(w.dt)
		w.animTime = w.animTime.Add(animDt)
		w.Itertate(w.dt * 1000)
	}

	w.meta.FPS = int(1000.0 / now.Sub(w.lastTime).Seconds())
	w.meta.IPF = int(worldDiff.Seconds() / w.dt.Seconds())
	w.meta.InterpolateTime = target.Sub(w.time) + w.dt

	w.lastTime = now

	w.Emit("step", w.meta)
}
