<p align="center">
    <img width="2509" height="2509" alt="GOPHER_ROCKS" src="https://github.com/user-attachments/assets/d6cbc3a0-1e0c-460a-87d3-07b3f6d4cb08" />
</p>

# iter — Lazy Iterators with Generics for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/enetx/iter.svg)](https://pkg.go.dev/github.com/enetx/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/enetx/iter)](https://goreportcard.com/report/github.com/enetx/iter)
[![Coverage Status](https://coveralls.io/repos/github/enetx/iter/badge.svg?branch=main&service=github)](https://coveralls.io/github/enetx/iter?branch=main)
[![Go](https://github.com/enetx/iter/actions/workflows/go.yml/badge.svg)](https://github.com/enetx/iter/actions/workflows/go.yml)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/enetx/iter)

`iter` is a library of lazy iterators for Go. Built on pure functions, no unnecessary allocations, with full generics support and a clean functional style.

Apply transformations like `Map`, `Filter`, `Take`, `Fold`, `Zip` directly to sequences — lazily and efficiently.

---

## 📦 Installation

```bash
go get github.com/enetx/iter
```

---

## 🚀 Example

```go
package main

import (
	"fmt"
	"github.com/enetx/iter"
)

func main() {
	// Extract even numbers and square them — chained, Rust-style:
	out := iter.FromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(x int) bool { return x%2 == 0 }). // 2, 4
		Map(func(x int) int { return x * x }).        // 4, 16
		ToSlice()
	fmt.Println(out) // [4 16]
}
```

Requires **Go 1.27+**: transformations are generic methods, so `Map`, `FilterMap`, `Scan`,
`Fold`, `Zip` can change element types right in the chain.

---

## ✨ Features

- ✅ **Lazy sequences** — transformations are not applied until iteration begins.
- ✅ **Generics-powered** — works with any types `T`, `U`, structs, etc.
- ✅ **Zero-allocation** — avoids unnecessary memory allocations.
- ✅ **Chainable methods** — every transformation is a method: `Map`, `Filter`, `Take`, `Skip`, `StepBy`, `Zip`, `Enumerate`, `Fold`, `Reduce`, and many more. Type-changing transformations are generic methods (Go 1.27).
- ✅ **Supports pairs, maps, channels** — `Seq2` with the same method set (no `*2` suffixes), `FromMap`, `FromChan`, and beyond.
- ✅ **Free functions where methods can't go** — element-type constraints (`Dedup`, `Contains`, `Counter`, `ToMap`) and slice-of-element results (`Windows`, `Chunks`, `Flatten`, `Combinations`, `Permutations`).

---

## 🔬 More Examples

```go
// Simple range with step
r := iter.Iota(1, 10) // 1..9
r = iter.StepBy(r, 3)  // 1, 4, 7
fmt.Println(iter.ToSlice(r)) // → [1 4 7]

// Iterate over a map
m := map[string]int{"a": 1, "b": 2}
s2 := iter.FromMap(m)
iter.ForEach2(s2, func(k string, v int) {
	fmt.Printf("%s → %d\n", k, v)
})
```

---

## ⚖️ License

MIT
