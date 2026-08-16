package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mykolapodpriatov/phi-algebra/boolop"
	"github.com/mykolapodpriatov/phi-algebra/classify"
	"github.com/mykolapodpriatov/phi-algebra/geom"
	"github.com/mykolapodpriatov/phi-algebra/phi"
)

type polyJSON struct {
	Vertices [][2]float64 `json:"vertices"`
}

func loadPoly(path string) (geom.Polygon, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return geom.Polygon{}, err
	}
	var doc polyJSON
	if err := json.Unmarshal(raw, &doc); err != nil {
		return geom.Polygon{}, err
	}
	verts := make([]geom.Point, len(doc.Vertices))
	for i, v := range doc.Vertices {
		verts[i] = geom.Point{X: v[0], Y: v[1]}
	}
	return geom.NewPolygon(verts)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: phi eval|classify|bool [--a file] [--b file] [--tx n] [--ty n]")
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	aPath := fs.String("a", "", "JSON polygon A")
	bPath := fs.String("b", "", "JSON polygon B")
	tx := fs.Float64("tx", 0, "translate B by tx")
	ty := fs.Float64("ty", 0, "translate B by ty")
	op := fs.String("op", "intersect", "union|intersect|subtract")
	must(fs.Parse(os.Args[2:]))
	if *aPath == "" || *bPath == "" {
		fmt.Fprintln(os.Stderr, "--a and --b are required")
		os.Exit(2)
	}
	a, err := loadPoly(*aPath)
	must(err)
	b, err := loadPoly(*bPath)
	must(err)
	t := geom.Vec{X: *tx, Y: *ty}
	bt := b.Translated(t)

	switch cmd {
	case "eval":
		r := phi.Of(a, b, t)
		enc := json.NewEncoder(os.Stdout)
		must(enc.Encode(map[string]any{
			"phi":      r.Value,
			"gradient": []float64{r.Gradient.X, r.Gradient.Y},
		}))
	case "classify":
		rel, r := classify.Of(a, b, t)
		fmt.Printf("%s  phi=%.6g\n", rel, r.Value)
	case "bool":
		var s boolop.Shape
		switch *op {
		case "union":
			s, err = boolop.Union(a, bt)
		case "intersect":
			s, err = boolop.Intersect(a, bt)
		case "subtract":
			s, err = boolop.Subtract(a, bt)
		default:
			must(fmt.Errorf("unknown op %q", *op))
		}
		must(err)
		fmt.Printf("%s\n", s.Kind)
		for i, p := range s.Parts {
			fmt.Printf("  part %d area=%.6g verts=%d\n", i, p.Area(), len(p.Verts))
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown command")
		os.Exit(2)
	}
}
