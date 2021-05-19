package collision

import (
	"errors"

	"github.com/oakmound/oak/v3/event"
)

// A Phase is a struct that other structs who want to use PhaseCollision
// should be composed of
type Phase struct {
	OnCollisionS *Space
	tree         *Tree
	// If allocating maps becomes an issue
	// we can have two constant maps that we
	// switch between on alternating frames
	Touching map[Label]bool
}

func (cp *Phase) getCollisionPhase() *Phase {
	return cp
}

type collisionPhase interface {
	getCollisionPhase() *Phase
}

// PhaseCollision binds to the entity behind the space's CID so that it will
// receive CollisionStart and CollisionStop events, appropriately when
// entities begin to collide or stop colliding with the space.
// If tree is nil, it uses DefTree
func PhaseCollision(s *Space, tree *Tree) error {
	en := s.CID.E()
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		oc.tree = tree
		if oc.tree == nil {
			oc.tree = DefaultTree
		}
		s.CID.Bind(event.Enter, phaseCollisionEnter)
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// CollisionStart/Stop: when a PhaseCollision entity starts/stops touching some label.
// Payload: (Label) the label the entity has started/stopped touching
const (
	Start = "CollisionStart"
	Stop  = "CollisionStop"
)

func phaseCollisionEnter(id event.CID, nothing interface{}) int {
	e := id.E().(collisionPhase)
	oc := e.getCollisionPhase()

	// check hits
	hits := oc.tree.Hits(oc.OnCollisionS)
	newTouching := map[Label]bool{}

	// if any are new, trigger on collision
	for _, h := range hits {
		l := h.Label
		if _, ok := oc.Touching[l]; !ok {
			event.CID(id).Trigger(Start, l)
		}
		newTouching[l] = true
	}

	// if we lost any, trigger off collision
	for l := range oc.Touching {
		if _, ok := newTouching[l]; !ok {
			event.CID(id).Trigger(Stop, l)
		}
	}

	oc.Touching = newTouching

	return 0
}
