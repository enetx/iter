package main

import (
	"fmt"

	"github.com/enetx/iter"
)

func main() {
	// Extract even numbers and square them:
	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
	s = iter.Filter(s, func(x int) bool { return x%2 == 0 })
	s = iter.Map(s, func(x int) int { return x * x })
	fmt.Println(iter.ToSlice(s)) // → [4 16]

	// Simple range with step
	r := iter.Iota(1, 10)        // 1..9
	r = iter.StepBy(r, 3)        // 1, 4, 7
	fmt.Println(iter.ToSlice(r)) // → [1 4 7]

	// Iterate over a map
	m := map[string]int{"a": 1, "b": 2}
	s2 := iter.FromMap(m)
	iter.ForEach2(s2, func(k string, v int) {
		fmt.Printf("%s → %d\n", k, v)
	})
}
