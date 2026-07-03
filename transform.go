package iter

// Map applies a function to each element, producing a new sequence.
// The result type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.Map(func(x int) int { return x * 2 }) // yields: 2, 4, 6
//	s.Map(func(x int) string { return fmt.Sprintf("%d", x) }) // yields: "1", "2", "3"
func (s Seq[T]) Map[U any](f func(T) U) Seq[U] {
	return func(yield func(U) bool) { s(func(v T) bool { return yield(f(v)) }) }
}

// Inspect applies a function to each element for side effects while passing through the original values.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3})
//	s.Inspect(func(x int) { fmt.Printf("Processing: %d\n", x) })
func (s Seq[T]) Inspect(fn func(T)) Seq[T] {
	return func(yield func(T) bool) { s(func(v T) bool { fn(v); return yield(v) }) }
}

// Filter returns a new sequence containing only elements that satisfy the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.Filter(func(x int) bool { return x%2 == 0 }) // yields: 2, 4
func (s Seq[T]) Filter(p func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		s(func(v T) bool {
			if p(v) {
				return yield(v)
			}
			return true
		})
	}
}

// Exclude returns a new sequence containing only elements that do not satisfy the predicate.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4, 5})
//	s.Exclude(func(x int) bool { return x%2 == 0 }) // yields: 1, 3, 5
func (s Seq[T]) Exclude(p func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		s(func(v T) bool {
			if !p(v) {
				return yield(v)
			}
			return true
		})
	}
}

// FilterMap applies a function to each element and filters out elements for which it returns false.
// The result type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]string{"1", "2", "abc", "3"})
//	s.FilterMap(func(s string) (int, bool) {
//	  if i, err := strconv.Atoi(s); err == nil {
//	    return i, true
//	  }
//	  return 0, false
//	}) // yields: 1, 2, 3
func (s Seq[T]) FilterMap[U any](f func(T) (U, bool)) Seq[U] {
	return func(yield func(U) bool) {
		s(func(v T) bool {
			if u, ok := f(v); ok {
				return yield(u)
			}
			return true
		})
	}
}

// MapWhile applies a function to elements while it returns true, stopping at the first false.
// The result type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, -1, 4})
//	s.MapWhile(func(x int) (int, bool) {
//	  if x > 0 { return x * 2, true }
//	  return 0, false
//	}) // yields: 2, 4
func (s Seq[T]) MapWhile[U any](f func(T) (U, bool)) Seq[U] {
	return func(yield func(U) bool) {
		s(func(v T) bool {
			if u, ok := f(v); ok {
				return yield(u)
			}
			return false
		})
	}
}

// Enumerate returns a sequence of (index, value) pairs starting from the given start index.
//
// Example:
//
//	s := iter.FromSlice([]string{"a", "b", "c"})
//	s.Enumerate(0) // yields: (0, "a"), (1, "b"), (2, "c")
func (s Seq[T]) Enumerate(start int) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := start
		s(func(v T) bool {
			result := yield(index, v)
			index++
			return result
		})
	}
}

// Scan is similar to Fold, but emits intermediate accumulator values.
// The accumulator type may differ from the element type.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 2, 3, 4})
//	s.Scan(0, func(acc, x int) int { return acc + x }) // yields: 1, 3, 6, 10
func (s Seq[T]) Scan[S any](init S, f func(S, T) S) Seq[S] {
	return func(yield func(S) bool) {
		acc := init
		s(func(v T) bool {
			acc = f(acc, v)
			return yield(acc)
		})
	}
}

// Unique returns a sequence with all duplicate elements removed.
// Works with any type by storing seen elements as map keys of type any.
//
// Because keys are stored as any, T is not required to be comparable at
// compile time: if T is not comparable at runtime (e.g. a slice, map,
// function, or a struct containing one), inserting an element into the
// seen map panics with "hash of unhashable type". For such types use
// DedupByKey with a comparable key function instead.
//
// This method is the fallback for element types that cannot satisfy the
// comparable constraint. When T is comparable, prefer the free function
// Unique, which uses a typed map[T]struct{} and avoids boxing every
// element into any (the same free-function/method split as Dedup/DedupBy
// and Counter/CounterBy).
//
// Example:
//
//	s := iter.FromSlice([]int{1, 1, 2, 2, 3, 1})
//	s.Unique() // yields: 1, 2, 3
func (s Seq[T]) Unique() Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[any]struct{})
		s(func(v T) bool {
			k := any(v)
			if _, exists := seen[k]; !exists {
				seen[k] = struct{}{}
				return yield(v)
			}
			return true
		})
	}
}

// Unique returns a sequence with all duplicate elements removed, keeping the
// first occurrence of each distinct value.
// It remains a free function because it requires T to be comparable,
// which cannot be expressed as a constraint on the receiver's type parameter.
// Unlike the Seq.Unique method, it stores seen elements in a typed
// map[T]struct{}, avoiding any-boxing. For non-comparable element types use
// the Seq.Unique method instead.
//
// Example:
//
//	s := iter.FromSlice([]int{1, 1, 2, 2, 3, 1})
//	iter.Unique(s) // yields: 1, 2, 3
func Unique[T comparable](s Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		s(func(v T) bool {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				return yield(v)
			}
			return true
		})
	}
}

// DedupByKey removes consecutive elements whose key function returns the same value,
// keeping the first element of each run. Unlike Unique, it does not track all seen keys:
// non-adjacent duplicates are preserved.
//
// Example:
//
//	s := iter.FromSlice([]string{"aa", "bb", "a", "ccc"})
//	s.DedupByKey(func(s string) int { return len(s) }) // yields: "aa", "a", "ccc"
func (s Seq[T]) DedupByKey[K comparable](key func(T) K) Seq[T] {
	return func(yield func(T) bool) {
		var prevKey K
		first := true
		s(func(v T) bool {
			k := key(v)
			if first || k != prevKey {
				prevKey = k
				first = false
				return yield(v)
			}
			return true
		})
	}
}
