package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"github.com/oniproject/physics.go/geometries"
)

type BodyCollisionDetection struct {
	Check   string // chan to listen to for collision candidates
	Channel string // chan to publish events to

	targets []bodies.Body
	world   World

	checkAllC, checkC func(interface{})
}

func NewBodyCollisionDetection() Behavior {
	b := &BodyCollisionDetection{
		Check:   "collisions:candidates",
		Channel: "collisions:detected",
	}
	b.checkC = func(data interface{}) { b.check(data.(map[int]*pair)) }
	b.checkAllC = func(data interface{}) { b.checkAll(data.(map[int]*pair)) }
	return b
}

func (b *BodyCollisionDetection) ApplyTo(bodies []bodies.Body) { b.targets = bodies }
func (b *BodyCollisionDetection) Targets() []bodies.Body       { return b.targets }
func (b *BodyCollisionDetection) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		if b.Check == "forse" || b.Check == "" {
			world.Off("integrate:velocities", &b.checkAllC)
		} else {
			world.Off(b.Check, &b.checkC)
		}
	}
	if world != nil {
		// connect
		if b.Check == "forse" || b.Check == "" {
			world.On("integrate:velocities", &b.checkAllC)
		} else {
			world.On(b.Check, &b.checkC)
		}
	}
	b.world = world
}

func (b *BodyCollisionDetection) check(candidates map[int]*pair) {
	collisions := []Collision{}
	for _, pair := range candidates {
		// TODO check if in b.Targets()
		ret, ok := b.checkPair(pair.bodyA, pair.bodyB)
		if ok {
			collisions = append(collisions, ret)
		}
	}
	if len(collisions) > 0 {
		b.world.Emit(b.Channel, collisions)
	}
}
func (b *BodyCollisionDetection) checkAll(candidates map[int]*pair) {
	collisions := []Collision{}
	targets := b.Targets()
	for j, bodyA := range targets {
		for i := j + 1; i < len(targets); i++ {
			bodyB := targets[i]
			ret, ok := b.checkPair(bodyA, bodyB)
			if ok {
				collisions = append(collisions, ret)
			}
		}
	}
	if len(collisions) > 0 {
		b.world.Emit(b.Channel, collisions)
	}
}

func (b *BodyCollisionDetection) checkPair(bodyA, bodyB bodies.Body) (c Collision, ok bool) {
	// filter out bodies that dont collide with each other
	if bodyA.Treatment() != bodies.TREATMENT_DYNAMIC &&
		bodyA.Treatment() != bodies.TREATMENT_DYNAMIC {
		return c, false
	}

	gA, isA := bodyA.Geometry().(*geometries.Circle)
	gB, isB := bodyB.Geometry().(*geometries.Circle)
	if isA && isB {
		return checkCircles(bodyA, bodyB, gA, gB)
		// TODO
		//} else {
		//return checkGJK(bodyA, bodyB)
	}
	return c, false
}

func checkCircles(bodyA, bodyB bodies.Body, gA, gB *geometries.Circle) (c Collision, ok bool) {
	d := bodyB.State().Pos.Minus(bodyA.State().Pos)
	overlap := d.Magnitude() - (gA.Radius + gB.Radius)

	// hmm... they overlap exactly... choose a direction
	if d.Equals(geom.Vector{}) {
		d = geom.Vector{1, 0}
	}

	/*if overlap > 0 {
		// check the future
		d = bodyB.State().Vel.Times(dt).Minus(bodyA.State().Vel.Times(dt))
		overlap = d.Magnitude() - (gA.Radius + gB.Radius)
	}*/

	if overlap <= 0 {
		d = d.Unit()
		ok = true
		c = Collision{
			BodyA:   bodyA,
			BodyB:   bodyB,
			Norm:    d,
			MTV:     d.Times(-overlap),
			Pos:     d.Times(gA.Radius),
			Overlap: -overlap,
		}
	}

	return
}

func checkGJK(bodyA, bodyB bodies.Body) {
	// TODO
}

type fnT struct {
	bodyA, bodyB     bodies.Body
	tA, tB           geom.Transform
	marginA, marginB float64
	useCore          bool
	fn               func(geom.Vector) geom.VectorABP
}

var supportFnStack = map[int]*fnT{}

func getSupportFnStack(bodyA, bodyB bodies.Body) *fnT {
	hash := pairHash(int(bodyA.UID()), int(bodyB.UID()))
	if hash == -1 {
		panic("fail hash")
	}
	fn := supportFnStack[hash]

	if fn == nil {
		fn = &fnT{useCore: false}
		supportFnStack[hash] = fn
		fn.fn = func(dir geom.Vector) geom.VectorABP {
			var vA, vB geom.Vector
			var mA, mB float64
			if fn.useCore {
				mA, mB = fn.marginA, fn.marginB
			}

			vA = bodyA.Geometry().FarthestCorePoint(fn.tA.RotateInv(dir), mA)
			vA = fn.tA.Translate(vA)
			vA = fn.tA.Rotate(vA)

			v := fn.tA.Rotate(dir)
			v = fn.tB.RotateInv(dir)
			vB = bodyA.Geometry().FarthestCorePoint(v.Times(-1), mB)
			vB = fn.tB.Translate(vB)
			vB = fn.tB.Rotate(vB)

			// WTF?  fn.Rotate(dir.Times(-1))

			return geom.VectorABP{
				A:  vA,
				B:  vB,
				PT: vA.Minus(vB),
			}
		}
	}

	fn.useCore = false
	//fn.margin = 0
	fn.tA.SetTranslation(bodyA.State().Pos)
	fn.tA.SetRotation(bodyA.State().Angular.Pos, geom.Vector{})

	fn.tB.SetTranslation(bodyB.State().Pos)
	fn.tB.SetRotation(bodyB.State().Angular.Pos, geom.Vector{})

	fn.bodyA = bodyA
	fn.bodyB = bodyB

	return fn
}
