// Package geom is planar geometry for convex polygons.
package geom

import (
	"fmt"
	"math"
)

const Eps = 1e-9

type Point struct {
	X, Y float64
}

type Vec struct {
	X, Y float64
}

func (p Point) Add(v Vec) Point { return Point{p.X + v.X, p.Y + v.Y} }
func (p Point) Sub(q Point) Vec { return Vec{p.X - q.X, p.Y - q.Y} }
func (v Vec) Add(w Vec) Vec     { return Vec{v.X + w.X, v.Y + w.Y} }
func (v Vec) Scale(s float64) Vec {
	return Vec{v.X * s, v.Y * s}
}
func (v Vec) Dot(w Vec) float64   { return v.X*w.X + v.Y*w.Y }
func (v Vec) Cross(w Vec) float64 { return v.X*w.Y - v.Y*w.X }
func (v Vec) Len() float64        { return math.Hypot(v.X, v.Y) }

func (v Vec) Normalized() Vec {
	l := v.Len()
	if l < Eps {
		return Vec{}
	}
	return v.Scale(1 / l)
}

func Orient(a, b, c Point) float64 {
	return b.Sub(a).Cross(c.Sub(a))
}

type Polygon struct {
	Verts []Point
}

func NewPolygon(verts []Point) (Polygon, error) {
	if len(verts) < 3 {
		return Polygon{}, fmt.Errorf("need at least 3 vertices, got %d", len(verts))
	}
	clean := dedupRing(verts)
	if len(clean) < 3 {
		return Polygon{}, fmt.Errorf("degenerate polygon")
	}
	p := Polygon{Verts: clean}
	if p.SignedArea() < 0 {
		p = p.Reversed()
	}
	if !p.IsConvex() {
		return Polygon{}, fmt.Errorf("v1 accepts convex polygons only")
	}
	return p, nil
}

func dedupRing(verts []Point) []Point {
	out := make([]Point, 0, len(verts))
	for i, v := range verts {
		prev := verts[(i+len(verts)-1)%len(verts)]
		if math.Hypot(v.X-prev.X, v.Y-prev.Y) > Eps {
			out = append(out, v)
		}
	}
	return out
}

func (p Polygon) Reversed() Polygon {
	n := len(p.Verts)
	out := make([]Point, n)
	for i, v := range p.Verts {
		out[n-1-i] = v
	}
	return Polygon{Verts: out}
}

func (p Polygon) SignedArea() float64 {
	var acc float64
	n := len(p.Verts)
	for i := 0; i < n; i++ {
		a := p.Verts[i]
		b := p.Verts[(i+1)%n]
		acc += a.X*b.Y - b.X*a.Y
	}
	return acc / 2
}

func (p Polygon) Area() float64 { return math.Abs(p.SignedArea()) }

func (p Polygon) IsConvex() bool {
	n := len(p.Verts)
	if n < 3 {
		return false
	}
	sign := 0
	for i := 0; i < n; i++ {
		c := Orient(p.Verts[i], p.Verts[(i+1)%n], p.Verts[(i+2)%n])
		if math.Abs(c) <= Eps {
			continue
		}
		s := 1
		if c < 0 {
			s = -1
		}
		if sign == 0 {
			sign = s
		} else if s != sign {
			return false
		}
	}
	return true
}

func (p Polygon) Translated(t Vec) Polygon {
	out := make([]Point, len(p.Verts))
	for i, v := range p.Verts {
		out[i] = v.Add(t)
	}
	q, err := NewPolygon(out)
	if err != nil {
		return Polygon{Verts: out}
	}
	return q
}

func (p Polygon) Contains(q Point) bool {
	n := len(p.Verts)
	for i := 0; i < n; i++ {
		if Orient(p.Verts[i], p.Verts[(i+1)%n], q) < -Eps {
			return false
		}
	}
	return true
}

func (p Polygon) ContainsStrict(q Point) bool {
	n := len(p.Verts)
	for i := 0; i < n; i++ {
		if Orient(p.Verts[i], p.Verts[(i+1)%n], q) <= Eps {
			return false
		}
	}
	return true
}

func (p Polygon) OutwardNormals() []Vec {
	n := len(p.Verts)
	out := make([]Vec, n)
	for i := 0; i < n; i++ {
		e := p.Verts[(i+1)%n].Sub(p.Verts[i])
		// CCW polygon: outward is rotate clockwise → (y, -x)
		out[i] = Vec{e.Y, -e.X}.Normalized()
	}
	return out
}

func Project(p Polygon, n Vec) (min, max float64) {
	min = p.Verts[0].X*n.X + p.Verts[0].Y*n.Y
	max = min
	for _, v := range p.Verts[1:] {
		s := v.X*n.X + v.Y*n.Y
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
	}
	return min, max
}

func ConvexHull(pts []Point) (Polygon, error) {
	if len(pts) < 3 {
		return Polygon{}, fmt.Errorf("not enough points for a hull")
	}
	sorted := append([]Point(nil), pts...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].X < sorted[i].X || (sorted[j].X == sorted[i].X && sorted[j].Y < sorted[i].Y) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	uniq := []Point{sorted[0]}
	for _, p := range sorted[1:] {
		last := uniq[len(uniq)-1]
		if math.Hypot(p.X-last.X, p.Y-last.Y) > Eps {
			uniq = append(uniq, p)
		}
	}
	lower := make([]Point, 0, len(uniq))
	for _, p := range uniq {
		for len(lower) >= 2 && Orient(lower[len(lower)-2], lower[len(lower)-1], p) <= Eps {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}
	upper := make([]Point, 0, len(uniq))
	for i := len(uniq) - 1; i >= 0; i-- {
		p := uniq[i]
		for len(upper) >= 2 && Orient(upper[len(upper)-2], upper[len(upper)-1], p) <= Eps {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}
	if len(lower) > 0 {
		lower = lower[:len(lower)-1]
	}
	if len(upper) > 0 {
		upper = upper[:len(upper)-1]
	}
	return NewPolygon(append(lower, upper...))
}

func SegmentIntersect(a1, a2, b1, b2 Point) (Point, bool) {
	p, t, ok := lineHit(a1, a2, b1, b2)
	if !ok || t < -Eps || t > 1+Eps {
		return Point{}, false
	}
	s := b2.Sub(b1)
	den := a2.Sub(a1).Cross(s)
	u := b1.Sub(a1).Cross(a2.Sub(a1)) / den
	if u < -Eps || u > 1+Eps {
		return Point{}, false
	}
	return p, true
}

// LineHit returns the intersection of segment a1-a2 with the infinite line b1-b2.
func LineHit(a1, a2, b1, b2 Point) (Point, bool) {
	p, t, ok := lineHit(a1, a2, b1, b2)
	if !ok || t < -Eps || t > 1+Eps {
		return Point{}, false
	}
	return p, true
}

func lineHit(a1, a2, b1, b2 Point) (Point, float64, bool) {
	r := a2.Sub(a1)
	s := b2.Sub(b1)
	den := r.Cross(s)
	if math.Abs(den) < Eps {
		return Point{}, 0, false
	}
	t := b1.Sub(a1).Cross(s) / den
	return a1.Add(r.Scale(t)), t, true
}
