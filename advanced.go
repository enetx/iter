package iter

// Cycle creates an infinite sequence by repeating the given sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.Cycle().Take(7) // yields: 1, 2, 3, 1, 2, 3, 1
func (s Seq[T]) Cycle() Seq[T] {
	return func(yield func(T) bool) {
		has := false
		s(func(T) bool {
			has = true
			return false
		})
		if !has {
			return
		}
		for {
			keep := true
			s(func(v T) bool {
				if !yield(v) {
					keep = false
					return false
				}
				return true
			})
			if !keep {
				return
			}
		}
	}
}

// Dedup removes consecutive duplicate elements.
// It remains a free function because it requires T to be comparable,
// which cannot be expressed as a constraint on the receiver's type parameter.
// For a method-chained variant use DedupBy.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 1, 2, 2, 2, 3, 3})
//	iter.Dedup(s) // yields: 1, 2, 3
func Dedup[T comparable](s Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		var prev T
		first := true
		s(func(v T) bool {
			if first || v != prev {
				prev = v
				first = false
				return yield(v)
			}
			return true
		})
	}
}

// DedupBy removes consecutive elements where the provided function returns the same value.
//
// Example:
//
//	s := iter.FromSlice([]string{"a", "aa", "b", "bb", "cc"})
//	s.DedupBy(func(a, b string) bool { return len(a) == len(b) })
//	// yields: "a", "b", "cc"
func (s Seq[T]) DedupBy(eq func(a, b T) bool) Seq[T] {
	return func(yield func(T) bool) {
		var prev T
		first := true
		s(func(v T) bool {
			if first || !eq(prev, v) {
				prev = v
				first = false
				return yield(v)
			}
			return true
		})
	}
}

// Intersperse inserts a separator between each element.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.Intersperse(0) // yields: 1, 0, 2, 0, 3
func (s Seq[T]) Intersperse(sep T) Seq[T] {
	return func(yield func(T) bool) {
		first := true
		s(func(v T) bool {
			if !first {
				if !yield(sep) {
					return false
				}
			}
			first = false
			return yield(v)
		})
	}
}

// Flatten flattens a sequence of slices into a single sequence.
// It remains a free function because the element type of the receiver
// cannot be constrained to a slice type.
//
// Example:
//
//	s := iter.FromSlice([][]int{{1, 2}, {3, 4}, {5}})
//	iter.Flatten(s) // yields: 1, 2, 3, 4, 5
func Flatten[T any](ss Seq[[]T]) Seq[T] {
	return func(yield func(T) bool) {
		ss(func(slice []T) bool {
			for _, v := range slice {
				if !yield(v) {
					return false
				}
			}
			return true
		})
	}
}

// FlattenSeq flattens a sequence of sequences into a single sequence.
// It remains a free function because the element type of the receiver
// cannot be constrained to a sequence type.
//
// Example:
//
//	seqs := []iter.Seq[int]{
//	  iter.FromSlice([]int{1, 2}),
//	  iter.FromSlice([]int{3, 4}),
//	}
//	s := iter.FromSlice(seqs)
//	iter.FlattenSeq(s) // yields: 1, 2, 3, 4
func FlattenSeq[T any](ss Seq[Seq[T]]) Seq[T] {
	return func(yield func(T) bool) {
		ss(func(s Seq[T]) bool {
			keep := true
			s(func(v T) bool {
				if !yield(v) {
					keep = false
					return false
				}
				return true
			})
			return keep
		})
	}
}

// FlatMap applies a function producing a sequence to each element and flattens the results.
// The result type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.FlatMap(func(x int) iter.Seq[int] { return iter.FromSlice([]int{x, x * 10}) })
//	// yields: 1, 10, 2, 20, 3, 30
func (s Seq[T]) FlatMap[U any](f func(T) Seq[U]) Seq[U] {
	return FlattenSeq(s.Map(f))
}

// Combinations generates all combinations of k elements from the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4})
//	iter.Combinations(s, 2) // yields: [1,2], [1,3], [1,4], [2,3], [2,4], [3,4]
func Combinations[T any](s Seq[T], k int) Seq[[]T] {
	return func(yield func([]T) bool) {
		slice := s.ToSlice()
		n := len(slice)

		if k > n || k <= 0 {
			return
		}

		indices := make([]int, k)
		for i := range indices {
			indices[i] = i
		}

		for {
			combination := make([]T, k)
			for i, idx := range indices {
				combination[i] = slice[idx]
			}

			if !yield(combination) {
				return
			}

			i := k - 1
			for i >= 0 && indices[i] == n-k+i {
				i--
			}

			if i < 0 {
				break
			}

			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}
		}
	}
}

// Permutations generates all permutations of the sequence elements.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	iter.Permutations(s) // yields: [1,2,3], [1,3,2], [2,1,3], [2,3,1], [3,1,2], [3,2,1]
func Permutations[T any](s Seq[T]) Seq[[]T] {
	return func(yield func([]T) bool) {
		slice := s.ToSlice()
		n := len(slice)

		if n == 0 {
			return
		}

		var generate func(int) bool
		generate = func(k int) bool {
			if k == 1 {
				perm := make([]T, n)
				copy(perm, slice)
				return yield(perm)
			}

			for i := range k {
				if !generate(k - 1) {
					return false
				}

				if k%2 == 0 {
					slice[i], slice[k-1] = slice[k-1], slice[i]
				} else {
					slice[0], slice[k-1] = slice[k-1], slice[0]
				}
			}

			return true
		}

		generate(n)
	}
}
