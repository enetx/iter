package iter_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/iter"
)

func TestNext2(t *testing.T) {
	// Test Next2 with multiple pairs
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	s := FromPairs(pairs)

	// Get first pair
	k1, v1, rest1, ok1 := Next2(s)
	if !ok1 || k1 != 1 || v1 != "a" {
		t.Errorf("Next2() first = %v, %v, %v, want 1, a, true", k1, v1, ok1)
	}

	// Get second pair
	k2, v2, rest2, ok2 := Next2(rest1)
	if !ok2 || k2 != 2 || v2 != "b" {
		t.Errorf("Next2() second = %v, %v, %v, want 2, b, true", k2, v2, ok2)
	}

	// Check remaining pairs
	remaining := ToPairs(rest2)
	expected := []Pair[int, string]{{3, "c"}, {4, "d"}}
	if !reflect.DeepEqual(remaining, expected) {
		t.Errorf("Next2() remaining = %v, want %v", remaining, expected)
	}

	// Test Next2 with single pair
	single := FromPairs([]Pair[int, string]{{42, "test"}})
	k, v, rest, ok := Next2(single)
	if !ok || k != 42 || v != "test" {
		t.Errorf("Next2() single = %v, %v, %v, want 42, test, true", k, v, ok)
	}
	if rest != nil {
		restPairs := ToPairs(rest)
		if len(restPairs) != 0 {
			t.Errorf("Next2() single rest = %v, want empty", restPairs)
		}
	}

	// Test Next2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	k, v, rest, ok = Next2(empty)
	if ok || k != 0 || v != "" || rest != nil {
		t.Errorf("Next2() empty = %v, %v, %v, %v, want 0, \"\", nil, false", k, v, rest, ok)
	}

	// Test Next2 iterating through all elements
	s2 := FromPairs([]Pair[int, string]{{1, "one"}, {2, "two"}, {3, "three"}})
	count := 0
	for {
		k, v, remaining, ok := Next2(s2)
		if !ok {
			break
		}
		count++
		s2 = remaining

		// Verify we got expected values
		switch count {
		case 1:
			if k != 1 || v != "one" {
				t.Errorf("Next2() iteration %d = %v, %v, want 1, one", count, k, v)
			}
		case 2:
			if k != 2 || v != "two" {
				t.Errorf("Next2() iteration %d = %v, %v, want 2, two", count, k, v)
			}
		case 3:
			if k != 3 || v != "three" {
				t.Errorf("Next2() iteration %d = %v, %v, want 3, three", count, k, v)
			}
		}
	}
	if count != 3 {
		t.Errorf("Next2() iteration count = %v, want 3", count)
	}
}

func TestFirst2(t *testing.T) {
	// Test First2 with non-empty sequence
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	s := FromPairs(pairs)
	k, v, ok := First2(s)
	if !ok || k != 1 || v != "a" {
		t.Errorf("First2() = %v, %v, %v, want 1, a, true", k, v, ok)
	}

	// Test First2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	_, _, ok2 := First2(empty)
	if ok2 {
		t.Errorf("First2(empty) ok = %v, want false", ok2)
	}

	// Test First2 with single pair
	single := FromPairs([]Pair[string, int]{{"hello", 42}})
	k3, v3, ok3 := First2(single)
	if !ok3 || k3 != "hello" || v3 != 42 {
		t.Errorf("First2(single) = %v, %v, %v, want hello, 42, true", k3, v3, ok3)
	}
}

func TestLast2(t *testing.T) {
	// Test Last2 with non-empty sequence
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	s := FromPairs(pairs)
	k, v, ok := Last2(s)
	if !ok || k != 3 || v != "c" {
		t.Errorf("Last2() = %v, %v, %v, want 3, c, true", k, v, ok)
	}

	// Test Last2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	_, _, ok2 := Last2(empty)
	if ok2 {
		t.Errorf("Last2(empty) ok = %v, want false", ok2)
	}

	// Test Last2 with single pair
	single := FromPairs([]Pair[string, int]{{"world", 99}})
	k3, v3, ok3 := Last2(single)
	if !ok3 || k3 != "world" || v3 != 99 {
		t.Errorf("Last2(single) = %v, %v, %v, want world, 99, true", k3, v3, ok3)
	}
}

func TestKeys(t *testing.T) {
	// Test keys operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToSlice(Keys(FromPairs(pairs)))
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Keys() = %v, want %v", result, expected)
	}
}

func TestValues(t *testing.T) {
	// Test values operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToSlice(Values(FromPairs(pairs)))
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Values() = %v, want %v", result, expected)
	}
}

func TestFromPairs(t *testing.T) {
	// Test fromPairs operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := ToPairs(FromPairs(pairs))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromPairs() = %v, want %v", result, expected)
	}
}

func TestFromPairsEarlyTermination(t *testing.T) {
	// Test fromPairs with early termination
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}

	count := 0
	FromPairs(pairs)(func(k int, v string) bool {
		count++
		return count < 3
	})

	if count != 3 {
		t.Errorf("FromPairsEarlyTermination count = %v, want 3", count)
	}
}

func TestOrderByKey(t *testing.T) {
	// Test orderByKey operation
	pairs := []Pair[int, string]{{3, "c"}, {1, "a"}, {2, "b"}}
	result := ToPairs(OrderByKey(FromPairs(pairs), func(a, b int) bool { return a < b }))
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByKey() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(OrderByKey(FromPairs([]Pair[int, string]{}), func(a, b int) bool { return a < b }))
	if len(result2) != 0 {
		t.Errorf("OrderByKey(empty) = %v, want empty", result2)
	}
}

func TestOrderByValue(t *testing.T) {
	// Test orderByValue operation
	pairs := []Pair[int, string]{{1, "c"}, {2, "a"}, {3, "b"}}
	result := ToPairs(OrderByValue(FromPairs(pairs), func(a, b string) bool { return a < b }))
	expected := []Pair[int, string]{{2, "a"}, {3, "b"}, {1, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderByValue() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(OrderByValue(FromPairs([]Pair[int, string]{}), func(a, b string) bool { return a < b }))
	if len(result2) != 0 {
		t.Errorf("OrderByValue(empty) = %v, want empty", result2)
	}
}

func TestSortBy2(t *testing.T) {
	// Test sortBy2 operation
	pairs := []Pair[string, int]{{"c", 3}, {"a", 1}, {"b", 2}}
	result := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	expected := []Pair[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortBy2() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := ToPairs(SortBy2(FromPairs([]Pair[string, int]{}), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	if len(result2) != 0 {
		t.Errorf("SortBy2(empty) = %v, want empty", result2)
	}

	// Test reverse sort by value
	result3 := ToPairs(SortBy2(FromPairs(pairs), func(a, b Pair[string, int]) bool {
		return a.Value > b.Value
	}))
	expected3 := []Pair[string, int]{{"c", 3}, {"b", 2}, {"a", 1}}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("SortBy2(reverse by value) = %v, want %v", result3, expected3)
	}

	// Test with equal elements (to cover return 0 case)
	equals := []Pair[string, int]{{"a", 1}, {"b", 1}, {"c", 1}}
	result4 := ToPairs(SortBy2(FromPairs(equals), func(a, b Pair[string, int]) bool {
		return a.Value < b.Value // All values are equal
	}))
	// Order should remain stable or be in some order
	if len(result4) != 3 {
		t.Errorf("SortBy2(equal values) length = %v, want 3", len(result4))
	}

	// Test single element
	single := []Pair[string, int]{{"a", 1}}
	result5 := ToPairs(SortBy2(FromPairs(single), func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}))
	if !reflect.DeepEqual(result5, single) {
		t.Errorf("SortBy2(single) = %v, want %v", result5, single)
	}

	// Test equal elements
	equal := []Pair[string, int]{{"a", 1}, {"a", 2}, {"a", 3}}
	result6 := ToPairs(SortBy2(FromPairs(equal), func(a, b Pair[string, int]) bool {
		return a.Value < b.Value
	}))
	expected5 := []Pair[string, int]{{"a", 1}, {"a", 2}, {"a", 3}}
	if !reflect.DeepEqual(result6, expected5) {
		t.Errorf("SortBy2(equal keys, sort by value) = %v, want %v", result6, expected5)
	}
}

func TestOrderByKeyAdvanced(t *testing.T) {
	// Test OrderByKey with single element
	single := []Pair[int, string]{{1, "a"}}
	result := ToPairs(OrderByKey(FromPairs(single), func(a, b int) bool { return a < b }))
	if !reflect.DeepEqual(result, single) {
		t.Errorf("OrderByKey(single) = %v, want %v", result, single)
	}

	// Test OrderByKey with already sorted data
	sorted := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result2 := ToPairs(OrderByKey(FromPairs(sorted), func(a, b int) bool { return a < b }))
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("OrderByKey(already sorted) = %v, want %v", result2, sorted)
	}

	// Test OrderByKey with duplicate keys
	duplicates := []Pair[int, string]{{2, "b"}, {1, "a"}, {2, "c"}, {1, "d"}}
	result3 := ToPairs(OrderByKey(FromPairs(duplicates), func(a, b int) bool { return a < b }))
	// Should have keys in order 1, 1, 2, 2
	if len(result3) != 4 || result3[0].Key != 1 || result3[1].Key != 1 || result3[2].Key != 2 || result3[3].Key != 2 {
		t.Errorf("OrderByKey(duplicates) = %v, want keys ordered as [1,1,2,2]", result3)
	}
}

func TestOrderByValueAdvanced(t *testing.T) {
	// Test OrderByValue with single element
	single := []Pair[int, string]{{1, "a"}}
	result := ToPairs(OrderByValue(FromPairs(single), func(a, b string) bool { return a < b }))
	if !reflect.DeepEqual(result, single) {
		t.Errorf("OrderByValue(single) = %v, want %v", result, single)
	}

	// Test OrderByValue with already sorted data
	sorted := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result2 := ToPairs(OrderByValue(FromPairs(sorted), func(a, b string) bool { return a < b }))
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("OrderByValue(already sorted) = %v, want %v", result2, sorted)
	}

	// Test OrderByValue with duplicate values
	duplicates := []Pair[int, string]{{2, "b"}, {1, "a"}, {3, "b"}, {4, "a"}}
	result3 := ToPairs(OrderByValue(FromPairs(duplicates), func(a, b string) bool { return a < b }))
	// Should have values in order a, a, b, b
	if len(result3) != 4 || result3[0].Value != "a" || result3[1].Value != "a" || result3[2].Value != "b" ||
		result3[3].Value != "b" {
		t.Errorf("OrderByValue(duplicates) = %v, want values ordered as [a,a,b,b]", result3)
	}
}

func TestFilterMap2(t *testing.T) {
	s := FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})

	result := ToMap(FilterMap2(s, func(k int, v string) (Pair[int, string], bool) {
		if len(v) > 1 {
			return Pair[int, string]{k * 10, strings.ToUpper(v)}, true
		}
		return Pair[int, string]{}, false
	}))

	expected := map[int]string{20: "BB", 30: "CCC"}

	if len(result) != len(expected) {
		t.Errorf("FilterMap2() result length = %d, want %d", len(result), len(expected))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("FilterMap2() result[%d] = %s, want %s", k, result[k], v)
		}
	}

	// Test empty map
	empty := ToMap(FilterMap2(FromMap(map[int]string{}), func(k int, v string) (Pair[int, string], bool) {
		return Pair[int, string]{k, v}, true
	}))
	if len(empty) != 0 {
		t.Errorf("FilterMap2(empty) = %v, want empty", empty)
	}

	// Test filter all out
	allFiltered := ToMap(FilterMap2(s, func(k int, v string) (Pair[int, string], bool) {
		return Pair[int, string]{k, v}, false // Filter all out
	}))
	if len(allFiltered) != 0 {
		t.Errorf("FilterMap2(filter all) = %v, want empty", allFiltered)
	}
}
