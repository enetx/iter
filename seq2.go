package iter

import "slices"

// Next extracts the first key-value pair from the sequence and returns the remaining sequence.
// Returns (key, value, remainingSeq, true) if a pair exists, or (zeroK, zeroV, nil, false) if empty.
// This is similar to Rust's Iterator::next() method for key-value pairs.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b", 3: "c"})
//	k, v, rest, ok := s.Next()
//	// k = 1, v = "a", ok = true (order not guaranteed for maps)
//	// rest yields remaining pairs
//
//	k2, v2, rest2, ok2 := rest.Next()
//	// k2 = 2, v2 = "b", ok2 = true
//	// rest2 yields remaining pairs
func (s Seq2[K, V]) Next() (K, V, Seq2[K, V], bool) {
	next, stop := s.Pull()

	firstK, firstV, ok := next()
	if !ok {
		stop()
		var zeroK K
		var zeroV V
		return zeroK, zeroV, nil, false
	}

	// The remaining sequence continues from the same pull iterator, so the source
	// is walked exactly once. This makes Next O(1) per element and correct for
	// non-deterministic sources (e.g. maps), at the cost of the remaining
	// sequence being single-use.
	remaining := func(yield func(K, V) bool) {
		defer stop()
		for {
			k, v, ok := next()
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}

	return firstK, firstV, remaining, true
}

// First returns the first key-value pair from the sequence.
// Returns (key, value, true) if a pair exists, or (zeroK, zeroV, false) if empty.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b", 3: "c"})
//	k, v, ok := s.First() // might return: 1, "a", true
func (s Seq2[K, V]) First() (K, V, bool) {
	var resultK K
	var resultV V
	found := false
	s(func(k K, v V) bool {
		resultK = k
		resultV = v
		found = true
		return false
	})
	return resultK, resultV, found
}

// Last returns the last key-value pair from the sequence.
// Returns (key, value, true) if a pair exists, or (zeroK, zeroV, false) if empty.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
//	s := iter.FromPairs(pairs)
//	k, v, ok := s.Last() // 3, "c", true
func (s Seq2[K, V]) Last() (K, V, bool) {
	var resultK K
	var resultV V
	found := false
	s(func(k K, v V) bool {
		resultK = k
		resultV = v
		found = true
		return true
	})
	return resultK, resultV, found
}

// ForEach applies a function to each key-value pair in the sequence.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	s.ForEach(func(k int, v string) { fmt.Printf("%d: %s\n", k, v) })
func (s Seq2[K, V]) ForEach(fn func(K, V)) {
	s(func(k K, v V) bool { fn(k, v); return true })
}

// Count returns the number of key-value pairs in the sequence.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	count := s.Count() // 2
func (s Seq2[K, V]) Count() int {
	count := 0
	s(func(K, V) bool { count++; return true })
	return count
}

// Range applies a function to each key-value pair until it returns false.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b", 3: "c"})
//	s.Range(func(k int, v string) bool {
//	  fmt.Printf("%d: %s\n", k, v)
//	  return k != 2 // Stop at key 2
//	})
func (s Seq2[K, V]) Range(fn func(K, V) bool) { s(fn) }

// Map applies a function to each key-value pair, producing a new Seq2.
// The key and value types of the result may differ from the input types.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	s.Map(func(k int, v string) (int, string) {
//	  return k*10, strings.ToUpper(v)
//	}) // yields: (10, "A"), (20, "B")
func (s Seq2[K, V]) Map[K2, V2 any](f func(K, V) (K2, V2)) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		s(func(k K, v V) bool {
			k2, v2 := f(k, v)
			return yield(k2, v2)
		})
	}
}

// Filter returns a new Seq2 containing only pairs that satisfy the predicate.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})
//	s.Filter(func(k int, v string) bool { return len(v) > 1 })
//	// yields pairs where value length > 1
func (s Seq2[K, V]) Filter(p func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s(func(k K, v V) bool {
			if p(k, v) {
				return yield(k, v)
			}
			return true
		})
	}
}

// Exclude returns a new Seq2 containing only pairs that do not satisfy the predicate.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})
//	s.Exclude(func(k int, v string) bool { return len(v) > 1 })
//	// yields pairs where value length <= 1
func (s Seq2[K, V]) Exclude(p func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s(func(k K, v V) bool {
			if !p(k, v) {
				return yield(k, v)
			}
			return true
		})
	}
}

// FilterMap applies a function to each key-value pair and filters out pairs for which
// it returns false. The key and value types of the result may differ from the input types.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})
//	s.FilterMap(func(k int, v string) (Pair[int, string], bool) {
//	  if len(v) > 1 {
//	    return Pair[int, string]{k*10, strings.ToUpper(v)}, true
//	  }
//	  return Pair[int, string]{}, false
//	}) // yields: (20, "BB"), (30, "CCC")
func (s Seq2[K, V]) FilterMap[K2, V2 any](f func(K, V) (Pair[K2, V2], bool)) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		s(func(k K, v V) bool {
			if pair, ok := f(k, v); ok {
				return yield(pair.Key, pair.Value)
			}
			return true
		})
	}
}

// MapWhile applies a function to key-value pairs while it returns true, stopping at the first false.
// The key and value types of the result may differ from the input types.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {-1, "c"}, {4, "d"}}
//	s := iter.FromPairs(pairs)
//	s.MapWhile(func(k int, v string) (Pair[int, string], bool) {
//	  if k > 0 {
//	    return Pair[int, string]{k * 10, strings.ToUpper(v)}, true
//	  }
//	  return Pair[int, string]{}, false
//	}) // yields: (10, "A"), (20, "B")
func (s Seq2[K, V]) MapWhile[K2, V2 any](f func(K, V) (Pair[K2, V2], bool)) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		s(func(k K, v V) bool {
			if pair, ok := f(k, v); ok {
				return yield(pair.Key, pair.Value)
			}
			return false
		})
	}
}

// Find returns the first key-value pair that satisfies the predicate.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})
//	k, v, ok := s.Find(func(k int, v string) bool { return len(v) > 1 })
//	// might return: 2, "bb", true
func (s Seq2[K, V]) Find(p func(K, V) bool) (K, V, bool) {
	var resultK K
	var resultV V
	found := false
	s(func(k K, v V) bool {
		if p(k, v) {
			resultK = k
			resultV = v
			found = true
			return false
		}
		return true
	})
	return resultK, resultV, found
}

// Keys extracts all keys from the Seq2 into a Seq.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	s.Keys() // yields: 1, 2 (order not guaranteed)
func (s Seq2[K, V]) Keys() Seq[K] {
	return func(y func(K) bool) { s(func(k K, _ V) bool { return y(k) }) }
}

// Values extracts all values from the Seq2 into a Seq.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	s.Values() // yields: "a", "b" (order not guaranteed)
func (s Seq2[K, V]) Values() Seq[V] {
	return func(y func(V) bool) { s(func(_ K, v V) bool { return y(v) }) }
}

// SortByKey sorts the sequence by keys using the comparison function.
// Note: This collects all elements into a slice first.
func (s Seq2[K, V]) SortByKey(less func(a, b K) bool) Seq2[K, V] {
	buf := s.ToPairs()
	slices.SortFunc(buf, func(a, b Pair[K, V]) int {
		switch {
		case less(a.Key, b.Key):
			return -1
		case less(b.Key, a.Key):
			return 1
		default:
			return 0
		}
	})
	return FromPairs(buf)
}

// SortByValue sorts the sequence by values using the comparison function.
// Note: This collects all elements into a slice first.
func (s Seq2[K, V]) SortByValue(less func(a, b V) bool) Seq2[K, V] {
	buf := s.ToPairs()
	slices.SortFunc(buf, func(a, b Pair[K, V]) int {
		switch {
		case less(a.Value, b.Value):
			return -1
		case less(b.Value, a.Value):
			return 1
		default:
			return 0
		}
	})
	return FromPairs(buf)
}

// Take returns the first n key-value pairs of the sequence.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b", 3: "c"})
//	s.Take(2) // yields first 2 pairs
func (s Seq2[K, V]) Take(n int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		count := 0
		s(func(k K, v V) bool {
			if count >= n {
				return false
			}
			count++
			return yield(k, v)
		})
	}
}

// Skip skips the first n key-value pairs and returns the rest.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b", 3: "c"})
//	s.Skip(1) // yields all pairs except the first
func (s Seq2[K, V]) Skip(n int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			s(yield)
			return
		}
		count := 0
		s(func(k K, v V) bool {
			if count < n {
				count++
				return true
			}
			return yield(k, v)
		})
	}
}

// TakeWhile yields key-value pairs while the predicate returns true.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
//	s := iter.FromPairs(pairs)
//	s.TakeWhile(func(k int, v string) bool { return k < 3 }) // yields: (1, "a"), (2, "b")
func (s Seq2[K, V]) TakeWhile(pred func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s(func(k K, v V) bool {
			if !pred(k, v) {
				return false
			}
			return yield(k, v)
		})
	}
}

// SkipWhile skips key-value pairs while the predicate returns true, then yields the rest.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
//	s := iter.FromPairs(pairs)
//	s.SkipWhile(func(k int, v string) bool { return k < 2 }) // yields: (2, "b"), (3, "c")
func (s Seq2[K, V]) SkipWhile(pred func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipping := true
		s(func(k K, v V) bool {
			if skipping && pred(k, v) {
				return true
			}
			skipping = false
			return yield(k, v)
		})
	}
}

// StepBy returns every nth key-value pair from the sequence.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
//	s := iter.FromPairs(pairs)
//	s.StepBy(2) // yields: (1, "a"), (3, "c")
func (s Seq2[K, V]) StepBy(step int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if step <= 0 {
			return
		}
		index := 0
		s(func(k K, v V) bool {
			if index%step == 0 {
				if !yield(k, v) {
					return false
				}
			}
			index++
			return true
		})
	}
}

// SortBy sorts the sequence using the provided comparison function on Pair values.
// Note: This collects all elements into a slice first.
func (s Seq2[K, V]) SortBy(less func(a, b Pair[K, V]) bool) Seq2[K, V] {
	buf := s.ToPairs()
	slices.SortFunc(buf, func(a, b Pair[K, V]) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
	return FromPairs(buf)
}

// Inspect applies a function to each key-value pair for side effects while passing through the original pairs.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	s.Inspect(func(k int, v string) { fmt.Printf("Processing: %d=%s\n", k, v) })
func (s Seq2[K, V]) Inspect(fn func(K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) { s(func(k K, v V) bool { fn(k, v); return yield(k, v) }) }
}

// Any returns true if any key-value pair satisfies the predicate.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "bb"})
//	hasLongValue := s.Any(func(k int, v string) bool { return len(v) > 1 }) // true
func (s Seq2[K, V]) Any(p func(K, V) bool) bool {
	found := false
	s(func(k K, v V) bool {
		if p(k, v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns true if all key-value pairs satisfy the predicate.
//
// Example:
//
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"})
//	allShortValues := s.All(func(k int, v string) bool { return len(v) == 1 }) // true
func (s Seq2[K, V]) All(p func(K, V) bool) bool {
	all := true
	s(func(k K, v V) bool {
		if !p(k, v) {
			all = false
			return false
		}
		return true
	})
	return all
}

// Fold reduces the Seq2 to a single value using an accumulator.
// The accumulator type may differ from the key and value types.
//
// Example:
//
//	s := iter.FromMap(map[int]int{1: 10, 2: 20, 3: 30})
//	sum := s.Fold(0, func(acc, k, v int) int { return acc + k + v }) // 66
func (s Seq2[K, V]) Fold[A any](acc A, f func(A, K, V) A) A {
	s(func(k K, v V) bool { acc = f(acc, k, v); return true })
	return acc
}

// Reduce reduces the Seq2 to a single Pair.
// Returns false if the sequence is empty.
//
// Example:
//
//	s := iter.FromMap(map[int]int{1: 10, 2: 20})
//	result, ok := s.Reduce(func(a, b Pair[int, int]) Pair[int, int] {
//	  return Pair[int, int]{a.Key + b.Key, a.Value + b.Value}
//	}) // Pair{3, 30}, true
func (s Seq2[K, V]) Reduce(f func(Pair[K, V], Pair[K, V]) Pair[K, V]) (Pair[K, V], bool) {
	var result Pair[K, V]
	first := true
	s(func(k K, v V) bool {
		kv := Pair[K, V]{k, v}
		if first {
			result = kv
			first = false
		} else {
			result = f(result, kv)
		}
		return true
	})
	return result, !first
}

// Nth returns the nth key-value pair (0-indexed) from the sequence.
//
// Example:
//
//	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
//	s := iter.FromPairs(pairs)
//	k, v, ok := s.Nth(1) // 2, "b", true
func (s Seq2[K, V]) Nth(n int) (K, V, bool) {
	var resultK K
	var resultV V
	found := false
	index := 0
	s(func(k K, v V) bool {
		if index == n {
			resultK = k
			resultV = v
			found = true
			return false
		}
		index++
		return true
	})
	return resultK, resultV, found
}
