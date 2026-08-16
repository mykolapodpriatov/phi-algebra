package phi_test

import (
	"math"
	"testing"

	"github.com/mykolapodpriatov/phi-algebra/boolop"
	"github.com/mykolapodpriatov/phi-algebra/classify"
	"github.com/mykolapodpriatov/phi-algebra/geom"
	"github.com/mykolapodpriatov/phi-algebra/phi"
)

func sq(x0, y0, s float64) geom.Polygon {
	p, err := geom.NewPolygon([]geom.Point{
		{X: x0, Y: y0}, {X: x0 + s, Y: y0}, {X: x0 + s, Y: y0 + s}, {X: x0, Y: y0 + s},
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TestPhiSigns(t *testing.T) {
	a := sq(0, 0, 1)
	cases := []struct {
		name string
		b    geom.Polygon
		want classify.Relation
	}{
		{"disjoint", sq(3, 0, 1), classify.Disjoint},
		{"touch", sq(1, 0, 1), classify.Touch},
		{"overlap", sq(0.4, 0, 1), classify.Overlap},
	}
	for _, tc := range cases {
		rel, r := classify.Of(a, tc.b, geom.Vec{})
		if rel != tc.want {
			t.Fatalf("%s: got %s Φ=%v", tc.name, rel, r.Value)
		}
	}
}

func TestPhiAgreesWithIntersection(t *testing.T) {
	a := sq(0, 0, 2)
	shifts := []geom.Vec{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 3, Y: 0},
		{X: 0.5, Y: 0.5},
		{X: 2, Y: 2},
	}
	for _, txy := range shifts {
		b := sq(0, 0, 1)
		r := phi.Of(a, b, txy)
		inter, err := boolop.Intersect(a, b.Translated(txy))
		if err != nil {
			t.Fatal(err)
		}
		has := inter.Kind != boolop.Empty
		if r.Value > geom.Eps && has {
			t.Fatalf("t=%v Φ=%v but intersection nonempty", txy, r.Value)
		}
		if r.Value < -geom.Eps && !has {
			t.Fatalf("t=%v Φ=%v but intersection empty", txy, r.Value)
		}
		if has && inter.Parts[0].Area() > math.Min(a.Area(), b.Area())+geom.Eps {
			t.Fatalf("intersection area too large: %v", inter.Parts[0].Area())
		}
	}
}

func TestGradientSeparates(t *testing.T) {
	a := sq(0, 0, 1)
	b := sq(0.3, 0, 1)
	r := phi.Of(a, b, geom.Vec{})
	if r.Value >= 0 {
		t.Fatalf("expected overlap, Φ=%v", r.Value)
	}
	// A step along +grad should increase Φ (less overlap / more separation).
	step := r.Gradient.Scale(math.Abs(r.Value) + 0.05)
	r2 := phi.Of(a, b, step)
	if r2.Value <= r.Value {
		t.Fatalf("Φ did not increase: %v -> %v grad=%v", r.Value, r2.Value, r.Gradient)
	}
}
