package boolop

import (
	"testing"

	"github.com/mykolapodpriatov/phi-algebra/geom"
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

func TestIntersectOverlap(t *testing.T) {
	a := sq(0, 0, 2)
	b := sq(1, 1, 2)
	s, err := Intersect(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if s.Kind != SimpleConnected {
		t.Fatalf("kind=%s", s.Kind)
	}
	if s.Parts[0].Area() < 0.99 || s.Parts[0].Area() > 1.01 {
		t.Fatalf("area=%v", s.Parts[0].Area())
	}
}

func TestUnionDisjoint(t *testing.T) {
	s, err := Union(sq(0, 0, 1), sq(3, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	if s.Kind != NonConnected || len(s.Parts) != 2 {
		t.Fatalf("%s %d", s.Kind, len(s.Parts))
	}
}

func TestSubtractDisjointAndContained(t *testing.T) {
	a := sq(0, 0, 4)
	s, err := Subtract(a, sq(5, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	if s.Kind != SimpleConnected || s.Parts[0].Area() < 15.9 {
		t.Fatalf("disjoint subtract: %s %v", s.Kind, s)
	}
	inner := sq(1, 1, 1)
	hole, err := Subtract(a, inner)
	if err != nil {
		t.Fatal(err)
	}
	if hole.Kind != NConnected || len(hole.Parts) != 2 {
		t.Fatalf("expected hole, got %s %d", hole.Kind, len(hole.Parts))
	}
}

func TestIntersectContained(t *testing.T) {
	a := sq(0, 0, 2)
	b := sq(0, 0, 1)
	s, err := Intersect(a, b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("kind=%s parts=%d", s.Kind, len(s.Parts))
	if s.Kind == SimpleConnected {
		t.Logf("area=%v verts=%v", s.Parts[0].Area(), s.Parts[0].Verts)
	}
	if s.Kind != SimpleConnected {
		t.Fatalf("expected intersection, got %s", s.Kind)
	}
}
