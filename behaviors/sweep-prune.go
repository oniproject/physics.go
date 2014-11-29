package behaviors

import (
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
	"log"
	"sort"
)

const (
	maxDof        = 2
	collisionFlag = 8 // 16 if 3D
)

var dof = geom.Vector{0, 1}
var uid int = 0

func getUniqueId() int {
	uid++
	return uid
}

type tracker struct {
	id       int
	body     bodies.Body
	max, min geom.Vector
}
type pair struct {
	bodyA, bodyB bodies.Body
	axis         [maxDof]bool
}

type SweepPrune struct {
	Channel string
	tracked [][]*tracker

	trackBodyC, untrackBodyC, sweepC func(interface{})

	axis int

	world World
}

func NewSweepPrune() Behavior {
	b := &SweepPrune{
		Channel: "collisions:candidates",
	}
	b.trackBodyC = func(data interface{}) { b.trackBody(data.(bodies.Body)) }
	b.untrackBodyC = func(data interface{}) { b.untrackBody(data.(bodies.Body)) }
	b.sweepC = func(data interface{}) { b.sweep() }

	b.clear()

	return b
}

func (b *SweepPrune) ApplyTo(bodies []bodies.Body) {}
func (b *SweepPrune) Targets() []bodies.Body       { return nil }
func (b *SweepPrune) SetWorld(world World) {
	if b.world != nil {
		// disconnect
		world.Off("add:body", &b.trackBodyC)
		world.Off("remove:body", &b.untrackBodyC)
		world.Off("integrate:positions", &b.sweepC)
		b.clear()
	}
	if world != nil {
		// connect
		world.On("add:body", &b.trackBodyC)
		world.On("remove:body", &b.untrackBodyC)
		world.On("integrate:positions", &b.sweepC)
	}
	b.world = world
}

func (b *SweepPrune) clear() {
	b.tracked = [][]*tracker{}
	for xyz := 0; xyz < maxDof; xyz++ {
		b.tracked = append(b.tracked, []*tracker{})
	}
}

func (b *SweepPrune) trackBody(body bodies.Body) {
	// TODO
	tracker := &tracker{
		id:   getUniqueId(),
		body: body,
	}

	for xyz := 0; xyz < maxDof; xyz++ {
		b.tracked[xyz] = append(b.tracked[xyz], tracker)
	}
}

func (b *SweepPrune) untrackBody(body bodies.Body) {
	for xyz := range b.tracked {
		for i, tracker := range b.tracked[xyz] {
			if tracker.body == body {
				b.tracked[xyz] = append(b.tracked[xyz][i:], b.tracked[xyz][:i+1]...)
			}
		}
	}
}

func (b *SweepPrune) sweep() {
	candidates := b.broadPhase()
	if len(candidates) > 0 {
		b.world.Emit(b.Channel, candidates)
	}
}

func (b *SweepPrune) broadPhase() map[int]*pair {
	for xyz := range b.tracked {
		for _, tr := range b.tracked[xyz] {
			aabb := tr.body.AABB(0)
			span := geom.Vector{aabb.HW, aabb.HH}
			pos := tr.body.State().Pos
			tr.min = pos.Minus(span)
			tr.max = pos.Plus(span)
		}
	}

	for xyz := 0; xyz < maxDof; xyz++ {
		list := b.tracked[xyz]
		switch xyz {
		case 0:
			sort.Sort(byX(list))
		case 1:
			sort.Sort(byY(list))
		}
		b.tracked[xyz] = list
	}
	return b.checkOverlaps()
}

type byX []*tracker

func (a byX) Len() int      { return len(a) }
func (a byX) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byX) Less(i, j int) bool {
	//first, second := a[i].value().X, a[j].value().X
	first, second := a[i].min.X, a[j].min.X
	return first < second //|| first == second && !a[i].isMax && a[j].isMax
}

type byY []*tracker

func (a byY) Len() int      { return len(a) }
func (a byY) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byY) Less(i, j int) bool {
	first, second := a[i].min.Y, a[j].min.Y
	return first < second //|| first == second && !a[i].isMax && a[j].isMax
}

func (b *SweepPrune) checkOverlaps() map[int]*pair {
	candidates := make(map[int]*pair)

	for axis, list := range b.tracked {
		//axis, list := b.axis, b.tracked[b.axis]
		for i, tr1 := range list {
			for j := i + 1; j < len(list); j++ {
				tr2 := list[j]

				if tr1 == tr2 {
					continue
				}

				switch axis {
				case 0:
					if tr2.min.X > tr1.max.X {
						break
					}
				case 1:
					if tr2.min.Y > tr1.max.Y {
						break
					}
				}
				hash := pairHash(tr1.id, tr2.id)

				if hash == -1 {
					log.Panic("fail hash", tr1.id, tr2.id)
				}
				candidates[hash] = &pair{
					bodyA: tr1.body,
					bodyB: tr2.body,
				}
			}
		}
	}

	return candidates
}
