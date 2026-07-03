package iter

// Next extracts the first element from the sequence and returns the remaining sequence.
// Returns (value, remainingSeq, true) if an element exists, or (zero, nil, false) if empty.
// This is similar to Rust's Iterator::next() method.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	val, rest, ok := s.Next()
//	// val = 1, ok = true
//	// rest yields: 2, 3, 4, 5
//
//	val2, rest2, ok2 := rest.Next()
//	// val2 = 2, ok2 = true
//	// rest2 yields: 3, 4, 5
func (s Seq[T]) Next() (T, Seq[T], bool) {
	next, stop := s.Pull()

	first, ok := next()
	if !ok {
		stop()
		var zero T
		return zero, nil, false
	}

	// The remaining sequence continues from the same pull iterator, so the source
	// is walked exactly once. This makes Next O(1) per element and correct for
	// non-deterministic sources (e.g. map-backed sets, channels), at the cost of
	// the remaining sequence being single-use.
	remaining := func(yield func(T) bool) {
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			if !yield(v) {
				return
			}
		}
	}

	return first, remaining, true
}

// First returns the first element from the sequence.
// Returns (value, true) if an element exists, or (zero, false) if empty.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	value, ok := s.First() // 1, true
func (s Seq[T]) First() (T, bool) {
	var result T
	found := false
	s(func(v T) bool {
		result = v
		found = true
		return false
	})
	return result, found
}

// Last returns the last element from the sequence.
// Returns (value, true) if an element exists, or (zero, false) if empty.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	value, ok := s.Last() // 5, true
func (s Seq[T]) Last() (T, bool) {
	var result T
	found := false
	s(func(v T) bool {
		result = v
		found = true
		return true
	})
	return result, found
}

// ForEach applies a function to each element in the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.ForEach(func(x int) { fmt.Println(x) })
func (s Seq[T]) ForEach(fn func(T)) {
	s(func(v T) bool { fn(v); return true })
}

// Count returns the number of elements in the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	count := s.Count() // 3
func (s Seq[T]) Count() int {
	count := 0
	s(func(T) bool {
		count++
		return true
	})
	return count
}

// Range applies a function to each element until it returns false.
// This is the same as the sequence's yield function.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.Range(func(x int) bool {
//	  fmt.Println(x)
//	  return x != 3 // Stop at 3
//	})
func (s Seq[T]) Range(fn func(T) bool) { s(fn) }

// Take returns the first n elements of the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.Take(3) // yields: 1, 2, 3
func (s Seq[T]) Take(n int) Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		count := 0
		s(func(v T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(v)
		})
	}
}

// Skip skips the first n elements and returns the rest.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.Skip(2) // yields: 3, 4, 5
func (s Seq[T]) Skip(n int) Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			s(yield)
			return
		}
		count := 0
		s(func(v T) bool {
			if count < n {
				count++
				return true
			}
			return yield(v)
		})
	}
}

// StepBy returns every nth element from the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5, 6})
//	s.StepBy(2) // yields: 1, 3, 5
func (s Seq[T]) StepBy(step int) Seq[T] {
	return func(yield func(T) bool) {
		if step <= 0 {
			return
		}
		index := 0
		s(func(v T) bool {
			if index%step == 0 {
				if !yield(v) {
					return false
				}
			}
			index++
			return true
		})
	}
}

// TakeWhile yields elements while the predicate returns true.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.TakeWhile(func(x int) bool { return x < 4 }) // yields: 1, 2, 3
func (s Seq[T]) TakeWhile(pred func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		s(func(v T) bool {
			if !pred(v) {
				return false
			}
			return yield(v)
		})
	}
}

// SkipWhile skips elements while the predicate returns true, then yields the rest.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.SkipWhile(func(x int) bool { return x < 3 }) // yields: 3, 4, 5
func (s Seq[T]) SkipWhile(pred func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		skipping := true
		s(func(v T) bool {
			if skipping && pred(v) {
				return true
			}
			skipping = false
			return yield(v)
		})
	}
}

// Nth returns the nth element (0-indexed) from the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	value, ok := s.Nth(2) // 3, true
func (s Seq[T]) Nth(n int) (T, bool) {
	var result T
	found := false
	index := 0
	s(func(v T) bool {
		if index == n {
			result = v
			found = true
			return false
		}
		index++
		return true
	})
	return result, found
}

// Contains checks if the sequence contains the given value.
// It remains a free function because it requires T to be comparable,
// which cannot be expressed as a constraint on the receiver's type parameter.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	has := iter.Contains(s, 2) // true
func Contains[T comparable](s Seq[T], x T) bool {
	found := false
	s(func(v T) bool {
		if v == x {
			found = true
			return false
		}
		return true
	})
	return found
}
