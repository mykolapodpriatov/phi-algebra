// Package boolop implements convex boolean operations and connectedness kind.
package boolop

import (
	"fmt"
	"math"

	"github.com/mykolapodpriatov/phi-algebra/classify"
	"github.com/mykolapodpriatov/phi-algebra/geom"
)

type Kind string

const (
	Empty           Kind = "empty"
	SimpleConnected Kind = "simple_connected"
	NConnected      Kind = "n_connected"
	NonConnected    Kind = "non_connected"
)

type Shape struct {
	Kind  Kind
	Parts []geom.Polygon // simple: one; n-connected: outer then holes; non-connected: components
}

func Intersect(a, b geom.Polygon) (Shape, error) {
	poly, ok := sutherlandHodgman(a, b)
	if !ok {
		return Shape{Kind: Empty}, nil
	}
	return Shape{Kind: SimpleConnected, Parts: []geom.Polygon{poly}}, nil
}

func Union(a, b geom.Polygon) (Shape, error) {
	rel, _ := classify.Of(a, b, geom.Vec{})
	if rel == classify.Disjoint {
		return Shape{Kind: NonConnected, Parts: []geom.Polygon{a, b}}, nil
	}
	pts := append(append([]geom.Point{}, a.Verts...), b.Verts...)
	h, err := geom.ConvexHull(pts)
	if err != nil {
		return Shape{}, err
	}
	return Shape{Kind: SimpleConnected, Parts: []geom.Polygon{h}}, nil
}

func Subtract(a, b geom.Polygon) (Shape, error) {
	rel, _ := classify.Of(a, b, geom.Vec{})
	if rel == classify.Disjoint {
		return Shape{Kind: SimpleConnected, Parts: []geom.Polygon{a}}, nil
	}
	if containsPoly(b, a) {
		return Shape{Kind: Empty}, nil
	}
	if containsPoly(a, b) && rel == classify.Overlap {
		return Shape{Kind: NConnected, Parts: []geom.Polygon{a, b}}, nil
	}
	diff, ok := convexDifference(a, b)
	if !ok {
		return Shape{Kind: Empty}, nil
	}
	return Shape{Kind: SimpleConnected, Parts: []geom.Polygon{diff}}, nil
}

func containsPoly(outer, inner geom.Polygon) bool {
	for _, v := range inner.Verts {
		if !outer.Contains(v) {
			return false
		}
	}
	return true
}

func sutherlandHodgman(subj, clip geom.Polygon) (geom.Polygon, bool) {
	output := append([]geom.Point{}, subj.Verts...)
	nclip := len(clip.Verts)
	for i := 0; i < nclip; i++ {
		a := clip.Verts[i]
		b := clip.Verts[(i+1)%nclip]
		input := append([]geom.Point(nil), output...)
		output = output[:0]
		if len(input) == 0 {
			return geom.Polygon{}, false
		}
		prev := input[len(input)-1]
		for _, cur := range input {
			curIn := geom.Orient(a, b, cur) >= -geom.Eps
			prevIn := geom.Orient(a, b, prev) >= -geom.Eps
			if curIn {
				if !prevIn {
					if hit, ok := geom.LineHit(prev, cur, a, b); ok {
						output = append(output, hit)
					}
				}
				output = append(output, cur)
			} else if prevIn {
				if hit, ok := geom.LineHit(prev, cur, a, b); ok {
					output = append(output, hit)
				}
			}
			prev = cur
		}
	}
	if len(output) < 3 {
		return geom.Polygon{}, false
	}
	p, err := geom.NewPolygon(output)
	if err != nil {
		return geom.Polygon{}, false
	}
	if p.Area() <= geom.Eps {
		return geom.Polygon{}, false
	}
	return p, true
}

func convexDifference(a, b geom.Polygon) (geom.Polygon, bool) {
	var pts []geom.Point
	n := len(a.Verts)
	for i := 0; i < n; i++ {
		cur := a.Verts[i]
		next := a.Verts[(i+1)%n]
		if !b.ContainsStrict(cur) {
			pts = append(pts, cur)
		}
		if hit, ok := firstHit(cur, next, b); ok {
			pts = append(pts, hit)
		}
	}
	if len(pts) < 3 {
		return geom.Polygon{}, false
	}
	// Result of a convex bite can be non-convex; build an ordered ring
	// from A-outside vertices + intersection points already in edge order.
	pts = dedupPts(pts)
	if len(pts) < 3 {
		return geom.Polygon{}, false
	}
	p := geom.Polygon{Verts: pts}
	if p.SignedArea() < 0 {
		p = p.Reversed()
	}
	if p.Area() <= geom.Eps {
		return geom.Polygon{}, false
	}
	return p, true
}

func firstHit(p, q geom.Point, poly geom.Polygon) (geom.Point, bool) {
	n := len(poly.Verts)
	var best geom.Point
	bestT := math.Inf(1)
	found := false
	d := q.Sub(p)
	for i := 0; i < n; i++ {
		hit, ok := geom.SegmentIntersect(p, q, poly.Verts[i], poly.Verts[(i+1)%n])
		if !ok {
			continue
		}
		var t float64
		if math.Abs(d.X) >= math.Abs(d.Y) && math.Abs(d.X) > geom.Eps {
			t = (hit.X - p.X) / d.X
		} else if math.Abs(d.Y) > geom.Eps {
			t = (hit.Y - p.Y) / d.Y
		} else {
			continue
		}
		if t > geom.Eps && t < 1-geom.Eps && t < bestT {
			bestT = t
			best = hit
			found = true
		}
	}
	return best, found
}

func dedupPts(pts []geom.Point) []geom.Point {
	out := make([]geom.Point, 0, len(pts))
	for _, p := range pts {
		dup := false
		for _, q := range out {
			if math.Hypot(p.X-q.X, p.Y-q.Y) <= geom.Eps {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, p)
		}
	}
	return out
}

func (s Shape) String() string {
	return fmt.Sprintf("%s (%d parts)", s.Kind, len(s.Parts))
}
