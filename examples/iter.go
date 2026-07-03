package main

import (
	"fmt"

	"github.com/enetx/iter"
)

func main() {
	// Extract even numbers and square them — chained, Rust-style:
	evens := iter.FromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		ToSlice()
	fmt.Println(evens) // → [4 16]

	// Map can change the element type right in the chain:
	strs := iter.FromSlice([]int{1, 2, 3}).
		Map(func(x int) string { return fmt.Sprintf("#%d", x) }).
		ToSlice()
	fmt.Println(strs) // → [#1 #2 #3]

	// Simple range with step
	r := iter.Iota(1, 10).StepBy(3) // 1, 4, 7
	fmt.Println(r.ToSlice())        // → [1 4 7]

	// Zip two sequences of different types
	iter.FromSlice([]int{1, 2, 3}).
		Zip(iter.FromSlice([]string{"one", "two", "three"})).
		ForEach(func(n int, s string) {
			fmt.Printf("%d → %s\n", n, s)
		})

	// Iterate over a map
	m := map[string]int{"a": 1, "b": 2}
	iter.FromMap(m).ForEach(func(k string, v int) {
		fmt.Printf("%s → %d\n", k, v)
	})
}
