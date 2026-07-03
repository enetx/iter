package iter

import "slices"

// Find returns the first element that satisfies the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	value, ok := s.Find(func(x int) bool { return x > 3 }) // 4, true
func (s Seq[T]) Find(p func(T) bool) (T, bool) {
	var result T
	found := false
	s(func(v T) bool {
		if p(v) {
			result = v
			found = true
			return false
		}
		return true
	})
	return result, found
}

// Any returns true if any element satisfies the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	hasEven := s.Any(func(x int) bool { return x%2 == 0 }) // true
func (s Seq[T]) Any(p func(T) bool) bool {
	found := false
	s(func(v T) bool {
		if p(v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns true if all elements satisfy the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{2, 4, 6, 8})
//	allEven := s.All(func(x int) bool { return x%2 == 0 }) // true
func (s Seq[T]) All(p func(T) bool) bool {
	all := true
	s(func(v T) bool {
		if !p(v) {
			all = false
			return false
		}
		return true
	})
	return all
}

// Fold reduces the sequence to a single value using an accumulator.
// The accumulator type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4})
//	sum := s.Fold(0, func(acc, x int) int { return acc + x }) // 10
func (s Seq[T]) Fold[A any](acc A, f func(A, T) A) A {
	s(func(v T) bool {
		acc = f(acc, v)
		return true
	})
	return acc
}

// Reduce reduces the sequence to a single value of the same type.
// Returns false if the sequence is empty.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4})
//	sum, ok := s.Reduce(func(a, b int) int { return a + b }) // 10, true
func (s Seq[T]) Reduce(f func(T, T) T) (T, bool) {
	var result T
	first := true
	s(func(v T) bool {
		if first {
			result = v
			first = false
		} else {
			result = f(result, v)
		}
		return true
	})
	return result, !first
}

// MinBy returns the minimum element according to the comparison function.
//
// Example:
//
//	s := iter.FromSlice([]int{3, 1, 4, 1, 5})
//	min, ok := s.MinBy(func(a, b int) bool { return a < b }) // 1, true
func (s Seq[T]) MinBy(less func(a, b T) bool) (T, bool) {
	var min T
	found := false
	s(func(v T) bool {
		if !found || less(v, min) {
			min = v
			found = true
		}
		return true
	})
	return min, found
}

// MaxBy returns the maximum element according to the comparison function.
//
// Example:
//
//	s := iter.FromSlice([]int{3, 1, 4, 1, 5})
//	max, ok := s.MaxBy(func(a, b int) bool { return a < b }) // 5, true
func (s Seq[T]) MaxBy(less func(a, b T) bool) (T, bool) {
	var max T
	found := false
	s(func(v T) bool {
		if !found || less(max, v) {
			max = v
			found = true
		}
		return true
	})
	return max, found
}

// CountBy counts elements by grouping them using a key function.
//
// Example:
//
//	s := iter.FromSlice([]string{"a", "bb", "ccc", "dd"})
//	counts := s.CountBy(func(s string) int { return len(s) })
//	// map[1:1 2:2 3:1]
func (s Seq[T]) CountBy[K comparable](key func(T) K) map[K]int {
	counts := make(map[K]int)
	s(func(v T) bool {
		k := key(v)
		counts[k]++
		return true
	})
	return counts
}

// Partition splits the sequence into two slices based on a predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	evens, odds := s.Partition(func(x int) bool { return x%2 == 0 })
//	// evens: [2, 4], odds: [1, 3, 5]
func (s Seq[T]) Partition(pred func(T) bool) (left, right []T) {
	s(func(v T) bool {
		if pred(v) {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
		return true
	})
	return left, right
}

// Counter returns a map with counts of each element.
// It remains a free function because it requires T to be comparable,
// which cannot be expressed as a constraint on the receiver's type parameter.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 1, 3, 2, 1})
//	counts := iter.Counter(s) // map[1:3 2:2 3:1]
func Counter[T comparable](s Seq[T]) map[T]int {
	counts := make(map[T]int)
	s(func(v T) bool {
		counts[v]++
		return true
	})
	return counts
}

// SortBy returns a sorted sequence using the provided comparison function.
// Note: This collects all elements into a slice first.
//
// Example:
//
//	s := iter.FromSlice([]int{3, 1, 4, 1, 5})
//	s.SortBy(func(a, b int) bool { return a < b }) // yields: 1, 1, 3, 4, 5
func (s Seq[T]) SortBy(less func(a, b T) bool) Seq[T] {
	slice := s.ToSlice()
	slices.SortFunc(slice, func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
	return FromSlice(slice)
}

// Position returns the index of the first element that satisfies the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	pos, ok := s.Position(func(x int) bool { return x > 3 }) // 3, true
func (s Seq[T]) Position(pred func(T) bool) (int, bool) {
	index := 0
	found := false
	s(func(v T) bool {
		if pred(v) {
			found = true
			return false
		}
		index++
		return true
	})

	if found {
		return index, true
	}
	return -1, false
}

// RPosition returns the index of the last element that satisfies the predicate.
// It performs a single forward pass over the sequence, tracking the last match,
// and uses O(1) memory.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 3, 2})
//	pos, ok := s.RPosition(func(x int) bool { return x == 3 }) // 4, true
func (s Seq[T]) RPosition(pred func(T) bool) (int, bool) {
	last := -1
	idx := 0
	s(func(v T) bool {
		if pred(v) {
			last = idx
		}
		idx++
		return true
	})
	return last, last >= 0
}

// IsPartitioned checks if the sequence is partitioned according to the predicate.
// A sequence is partitioned if all elements satisfying the predicate come before
// all elements that don't.
//
// Example:
//
//	s := iter.FromSlice([]int{2, 4, 6, 1, 3, 5})
//	partitioned := s.IsPartitioned(func(x int) bool { return x%2 == 0 }) // true
func (s Seq[T]) IsPartitioned(pred func(T) bool) bool {
	foundFalse := false
	result := true
	s(func(v T) bool {
		p := pred(v)
		if foundFalse && p {
			result = false
			return false
		}
		if !p {
			foundFalse = true
		}
		return true
	})
	return result
}
