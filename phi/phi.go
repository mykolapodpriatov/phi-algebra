// Package phi evaluates the Stoyan Φ-function of two convex polygons.
package phi

import (
	"math"

	"github.com/mykolapodpriatov/phi-algebra/geom"
)

type Result struct {
	Value    float64
	Gradient geom.Vec
}

// Of returns Φ(A, B + t). Sign: >0 disjoint, 0 touch, <0 overlap.
// Value is a SAT gap (separation or negated penetration). Gradient is the
// active axis: moving B along it increases Φ.
func Of(a, b geom.Polygon, t geom.Vec) Result {
	bt := b.Translated(t)
	axes := append(a.OutwardNormals(), bt.OutwardNormals()...)
	best := math.Inf(-1)
	bestN := geom.Vec{}
	for _, n := range axes {
		if n.Len() < geom.Eps {
			continue
		}
		n = n.Normalized()
		_, maxA := geom.Project(a, n)
		minB, _ := geom.Project(bt, n)
		gap := minB - maxA
		if gap > best {
			best = gap
			bestN = n
		}
	}
	if math.IsInf(best, -1) {
		return Result{}
	}
	return Result{Value: best, Gradient: bestN}
}
