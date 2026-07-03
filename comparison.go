package iter

import "reflect"

// Cmp compares two sequences lexicographically using the provided comparison function.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 3})
//	s2 := iter.FromSlice([]int{1, 2, 4})
//	result := s1.Cmp(s2, func(a, b int) int {
//	  if a < b { return -1 }
//	  if a > b { return 1 }
//	  return 0
//	}) // -1
func (a Seq[T]) Cmp(b Seq[T], cmp func(T, T) int) int {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if !aok && !bok {
			return 0
		}
		if !aok {
			return -1
		}
		if !bok {
			return 1
		}

		if c := cmp(av, bv); c != 0 {
			return c
		}
	}
}

// Equal checks if two sequences are equal (same elements in same order).
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 3})
//	s2 := iter.FromSlice([]int{1, 2, 3})
//	equal := s1.Equal(s2) // true
func (a Seq[T]) Equal(b Seq[T]) bool {
	var zero T
	if reflect.ValueOf(zero).Comparable() {
		return a.EqualBy(b, func(x, y T) bool { return any(x) == any(y) })
	}
	return a.EqualBy(b, func(x, y T) bool { return reflect.DeepEqual(x, y) })
}

// EqualBy checks if two sequences are equal using a custom equality function
// (same length and pairwise-equal elements in the same order).
//
// Example:
//
//	s1 := iter.FromSlice([]string{"a", "B"})
//	s2 := iter.FromSlice([]string{"A", "b"})
//	equal := s1.EqualBy(s2, strings.EqualFold) // true
func (a Seq[T]) EqualBy(b Seq[T], eq func(T, T) bool) bool {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if aok != bok {
			return false
		}
		if !aok {
			return true
		}
		if !eq(av, bv) {
			return false
		}
	}
}

// Lt checks if sequence a is lexicographically less than sequence b.
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 3})
//	s2 := iter.FromSlice([]int{1, 2, 4})
//	isLess := s1.Lt(s2, func(a, b int) bool { return a < b }) // true
func (a Seq[T]) Lt(b Seq[T], less func(T, T) bool) bool {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if !aok && bok {
			return true
		}
		if !aok || !bok {
			return false
		}

		if less(av, bv) {
			return true
		}
		if less(bv, av) {
			return false
		}
	}
}

// Le checks if sequence a is lexicographically less than or equal to sequence b.
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 3})
//	s2 := iter.FromSlice([]int{1, 2, 3})
//	isLessOrEqual := s1.Le(s2, func(a, b int) bool { return a < b }) // true
func (a Seq[T]) Le(b Seq[T], less func(T, T) bool) bool {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if !aok {
			return true
		}
		if !bok {
			return false
		}

		if less(av, bv) {
			return true
		}
		if less(bv, av) {
			return false
		}
	}
}

// Gt checks if sequence a is lexicographically greater than sequence b.
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 4})
//	s2 := iter.FromSlice([]int{1, 2, 3})
//	isGreater := s1.Gt(s2, func(a, b int) bool { return a < b }) // true
func (a Seq[T]) Gt(b Seq[T], less func(T, T) bool) bool {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if !bok && aok {
			return true
		}
		if !aok || !bok {
			return false
		}

		if less(bv, av) {
			return true
		}
		if less(av, bv) {
			return false
		}
	}
}

// Ge checks if sequence a is lexicographically greater than or equal to sequence b.
//
// Example:
//
//	s1 := iter.FromSlice([]int{1, 2, 3})
//	s2 := iter.FromSlice([]int{1, 2, 3})
//	isGreaterOrEqual := s1.Ge(s2, func(a, b int) bool { return a < b }) // true
func (a Seq[T]) Ge(b Seq[T], less func(T, T) bool) bool {
	an, as := a.Pull()
	defer as()
	bn, bs := b.Pull()
	defer bs()

	for {
		av, aok := an()
		bv, bok := bn()

		if !bok {
			return true
		}
		if !aok {
			return false
		}

		if less(bv, av) {
			return true
		}
		if less(av, bv) {
			return false
		}
	}
}
