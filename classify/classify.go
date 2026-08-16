package classify

import (
	"github.com/mykolapodpriatov/phi-algebra/geom"
	"github.com/mykolapodpriatov/phi-algebra/phi"
)

type Relation string

const (
	Disjoint Relation = "disjoint"
	Touch    Relation = "touch"
	Overlap  Relation = "overlap"
)

func Of(a, b geom.Polygon, t geom.Vec) (Relation, phi.Result) {
	r := phi.Of(a, b, t)
	switch {
	case r.Value > geom.Eps:
		return Disjoint, r
	case r.Value < -geom.Eps:
		return Overlap, r
	default:
		return Touch, r
	}
}
