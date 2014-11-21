package behaviors

import (
	//"errors"
	"github.com/oniproject/physics.go/bodies"
	"github.com/oniproject/physics.go/geom"
)

const maxDof = 2

var dof = geom.Vector{0, 1}
var uid int64 = 0

func getUniqueId() int64 {
	uid++
	return uid
}

func pairHash(id1, id2 int64) int64 {
	switch {
	case id1 > id2:
		return id1<<16 | id2&0xFFFF
	case id1 < id2:
		return id2<<16 | id1&0xFFFF
	default:
		return -1
	}

}

type intr struct {
	t       bool
	val     geom.Vector
	tracker *tracker
}
type interval struct {
	min intr
	max intr
}
type tracker struct {
	id       int64
	body     bodies.Body
	interval interval
}

type pair struct {
	bodyA, bodyB bodies.Body
	flag         uint
}

type SweepPrune struct {
	//encounters []...
	//candidates []...
	tracked      []tracker
	pairs        map[int64]pair
	intervalList [][]intr

	trackBodyC, untrackBodyC, sweepC func(interface{})

	world World
}

func NewSweepPrune() Behavior {
	b := &SweepPrune{}
	//b.trackBodyC = func(data interface{}) { b.trackBody(data.(bodies.Body)) }
	//b.untrackBodyC = func(data interface{}) { b.untrackBody(data.(bodies.Body)) }
	//b.sweepC = func(data interface{}) { b.sweep() }
	return b
}

func (b *SweepPrune) ApplyTo(bodies []bodies.Body) {}
func (b *SweepPrune) Targets() []bodies.Body       { return nil }
func (b *SweepPrune) SetWorld(world World) {
	if b.world != nil {
		b.disconnect(b.world)
		b.world = nil
	}
	if world != nil {
		b.connect(b.world)
	}
}

func (b *SweepPrune) connect(world World) {
	world.On("add:body", &b.trackBodyC)
	world.On("remove:body", &b.untrackBodyC)
	world.On("integrate:velocities", &b.sweepC)

	//for _, body := range world.Bodies() {
	//b.trackBody(body)
	//}
}

func (b *SweepPrune) disconnect(world World) {
	world.Off("add:body", &b.trackBodyC)
	world.Off("remove:body", &b.untrackBodyC)
	world.Off("integrate:velocities", &b.sweepC)
	b.clear()
}

func (b *SweepPrune) clear() {
	b.tracked = []tracker{}
	b.pairs = make(map[int64]pair)
	b.intervalList = [][]intr{}

	for xyz := 0; xyz < maxDof; xyz++ {
		b.intervalList = append(b.intervalList, []intr{})
	}
}

/*func (b *SweepPrune) broadPhase() {
	b.updateIntervals()
	b.sortIntervals()
	return b.checkOverlaps()
}

func (b *SweepPrune) sortIntervals() {
	// TODO
	for xyz := 0; xyz < maxDof; xyz++ {
		list := b.intervalList[xyz]
		axis := xyz

		for i, bound := range list {
			hole := i

			left := list[hole-1]

			var boundVal, leftVal float64

			switch axis {
			case 0:
				boundVal = bound.val.X
				leftVal = left.val.X
			case 1:
				boundVal = bound.val.Y
				leftVal = left.val.Y
			}

			for hole > 0 &&
				(leftVal > boundVal ||
					leftVal == boundVal && (left.t && !bound.t)) {
				list[hole] = left
				hole--
				left = list[hole-1]
				switch axis {
				case 0:
					leftVal = left.val.X
				case 1:
					leftVal = left.val.Y
				}
			}

			list[hole] = bound
		}
	}
}

func (b *SweepPrune) getPair(tr1, tr2 tracker, doCreate bool) (pair, error) {
	// TODO

	hash := pairHash(tr1.id, tr2.id)

	if hash == -1 {
		return pair{}, errors.New("fail hash")
	}

	if c, ok := b.pairs[hash]; !ok {
		if !doCreate {
			return c, errors.New("!doCreate")
		}
		c = pair{
			bodyA: tr1.body,
			bodyB: tr2.body,
			flag:  1,
		}
		b.pairs[hash] = c
		return c, nil
	} else {
		if doCreate {
			c.flag = 1
		}
		return c, nil
	}
}

func (b *SweepPrune) checkOverlaps() int {
	// TODO

}

func (b *SweepPrune) updateIntervals() {
	// TODO
}
func (b *SweepPrune) trackBody(body bodies.Body) {
	// TODO
	tracker := tracker{
		id:   getUniqueId(),
		body: body,
	}
	intr := interval{
		min: intr{
			t:       false, //min
			tracker: &tracker,
		},
		max: intr{
			t:       true, //max
			tracker: &tracker,
		},
	}

	tracker.interval = intr
	b.tracked = append(b.tracked, tracker)

	for xyz := 0; xyz < maxDof; xyz++ {
		b.intervalList[xyz] = append(b.intervalList[xyz], intr.min, intr.max)
	}
}

func (b *SweepPrune) untrackBody(body bodies.Body) {
	// TODO
	for i, tracker := range b.tracked {
		if tracker.body == body {
			b.tracked = append(b.tracked[i:], b.tracked[:i+1]...)

			for xyz := 0; xyz < maxDof; xyz++ {
				count := 0
				list := b.intervalList[xyz]

			XX:
				for j, minmax := range list {
					if minmax == tracker.interval.min || minmax == tracker.interval.max {
						list = append(list[j:], list[:j+1]...)
						if count > 0 {
							break XX
						}
						count++
					}
				}
			}
		}
	}
}

func (b *SweepPrune) sweep() {
	// TODO
	candidates := b.broadPhase()
	if len(candidates) > 0 {
		b.world.Emit("collisions:candidates", candidates)
	}
}*/
