package geom

import "testing"

func square(x0, y0, s float64) Polygon {
	p, err := NewPolygon([]Point{
		{x0, y0}, {x0 + s, y0}, {x0 + s, y0 + s}, {x0, y0 + s},
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TestAreaAndConvex(t *testing.T) {
	p := square(0, 0, 2)
	if p.Area() < 3.99 || p.Area() > 4.01 {
		t.Fatalf("area=%v", p.Area())
	}
	if !p.IsConvex() {
		t.Fatal("expected convex")
	}
	if !p.Contains(Point{1, 1}) || p.Contains(Point{-1, 0}) {
		t.Fatal("contains")
	}
}

func TestRejectsConcave(t *testing.T) {
	_, err := NewPolygon([]Point{
		{0, 0}, {3, 0}, {1, 1}, {3, 2}, {0, 2},
	})
	if err == nil {
		t.Fatal("expected concave rejection")
	}
}
