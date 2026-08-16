# phi-algebra

Φ-functions and boolean algebra for **convex** parts. The math kernel from an NTU «KhPI» diploma (Stoyan school); `nest-playground` will sit on top of this. No UI.

```
Φ > 0  disjoint     (gap)
Φ = 0  touching
Φ < 0  overlapping  (negated penetration)
```

```bash
go test ./...
go run ./cmd/phi classify --a testdata/square.json --b testdata/small.json --tx 3
go run ./cmd/phi eval --a testdata/square.json --b testdata/small.json --tx 0.5 --ty 0.5
go run ./cmd/phi bool --op intersect --a testdata/square.json --b testdata/small.json --tx 1 --ty 1
```

v1 is convex-only. Union of overlapping parts is the convex hull of the vertices (honest superset). Subtract of a part strictly inside another is tagged `n_connected` (outer + hole).

MIT. Rewritten from Diplom 7.0 `FiObject`, not a C++ Builder dump.
