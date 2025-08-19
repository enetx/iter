package iter

// ForEach applies a function to each element in the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	iter.ForEach(s, func(x int) { fmt.Println(x) })
func ForEach[T any](s Seq[T], fn func(T)) {
	s(func(v T) bool { fn(v); return true })
}

// Count returns the number of elements in the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	count := iter.Count(s) // 3
func Count[T any](s Seq[T]) int {
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
//	iter.Range(s, func(x int) bool {
//	  fmt.Println(x)
//	  return x != 3 // Stop at 3
//	})
func Range[T any](s Seq[T], fn func(T) bool) { s(fn) }

// Take returns the first n elements of the sequence.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	iter.Take(s, 3) // yields: 1, 2, 3
func Take[T any](s Seq[T], n int) Seq[T] {
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
//	iter.Skip(s, 2) // yields: 3, 4, 5
func Skip[T any](s Seq[T], n int) Seq[T] {
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
//	iter.StepBy(s, 2) // yields: 1, 3, 5
func StepBy[T any](s Seq[T], step int) Seq[T] {
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
//	iter.TakeWhile(s, func(x int) bool { return x < 4 }) // yields: 1, 2, 3
func TakeWhile[T any](s Seq[T], pred func(T) bool) Seq[T] {
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
//	iter.SkipWhile(s, func(x int) bool { return x < 3 }) // yields: 3, 4, 5
func SkipWhile[T any](s Seq[T], pred func(T) bool) Seq[T] {
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
//	value, ok := iter.Nth(s, 2) // 3, true
func Nth[T any](s Seq[T], n int) (T, bool) {
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
