package iter_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/enetx/iter"
)

func TestNext2(t *testing.T) {
	// Test Next2 with multiple pairs
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}, {Key: 4, Value: "d"}}
	s := FromPairs(pairs)

	// Get first pair
	k1, v1, rest1, ok1 := s.Next()
	if !ok1 || k1 != 1 || v1 != "a" {
		t.Errorf(".Next() first = %v, %v, %v, want 1, a, true", k1, v1, ok1)
	}

	// Get second pair
	k2, v2, rest2, ok2 := rest1.Next()
	if !ok2 || k2 != 2 || v2 != "b" {
		t.Errorf(".Next() second = %v, %v, %v, want 2, b, true", k2, v2, ok2)
	}

	// Check remaining pairs
	remaining := rest2.ToPairs()
	expected := []Pair[int, string]{{Key: 3, Value: "c"}, {Key: 4, Value: "d"}}
	if !reflect.DeepEqual(remaining, expected) {
		t.Errorf(".Next() remaining = %v, want %v", remaining, expected)
	}

	// Test Next2 with single pair
	single := FromPairs([]Pair[int, string]{{Key: 42, Value: "test"}})
	k, v, rest, ok := single.Next()
	if !ok || k != 42 || v != "test" {
		t.Errorf(".Next() single = %v, %v, %v, want 42, test, true", k, v, ok)
	}
	if rest != nil {
		restPairs := rest.ToPairs()
		if len(restPairs) != 0 {
			t.Errorf(".Next() single rest = %v, want empty", restPairs)
		}
	}

	// Test Next2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	k, v, rest, ok = empty.Next()
	if ok || k != 0 || v != "" || rest != nil {
		t.Errorf(".Next() empty = %v, %v, %v, %v, want 0, \"\", nil, false", k, v, rest, ok)
	}

	// Test Next2 iterating through all elements
	s2 := FromPairs([]Pair[int, string]{{Key: 1, Value: "one"}, {Key: 2, Value: "two"}, {Key: 3, Value: "three"}})
	count := 0
	for {
		k, v, remaining, ok := s2.Next()
		if !ok {
			break
		}
		count++
		s2 = remaining

		// Verify we got expected values
		switch count {
		case 1:
			if k != 1 || v != "one" {
				t.Errorf(".Next() iteration %d = %v, %v, want 1, one", count, k, v)
			}
		case 2:
			if k != 2 || v != "two" {
				t.Errorf(".Next() iteration %d = %v, %v, want 2, two", count, k, v)
			}
		case 3:
			if k != 3 || v != "three" {
				t.Errorf(".Next() iteration %d = %v, %v, want 3, three", count, k, v)
			}
		}
	}
	if count != 3 {
		t.Errorf(".Next() iteration count = %v, want 3", count)
	}
}

func TestFirst2(t *testing.T) {
	// Test First2 with non-empty sequence
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	s := FromPairs(pairs)
	k, v, ok := s.First()
	if !ok || k != 1 || v != "a" {
		t.Errorf(".First() = %v, %v, %v, want 1, a, true", k, v, ok)
	}

	// Test First2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	_, _, ok2 := empty.First()
	if ok2 {
		t.Errorf("empty.First() ok = %v, want false", ok2)
	}

	// Test First2 with single pair
	single := FromPairs([]Pair[string, int]{{Key: "hello", Value: 42}})
	k3, v3, ok3 := single.First()
	if !ok3 || k3 != "hello" || v3 != 42 {
		t.Errorf("single.First() = %v, %v, %v, want hello, 42, true", k3, v3, ok3)
	}
}

func TestLast2(t *testing.T) {
	// Test Last2 with non-empty sequence
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	s := FromPairs(pairs)
	k, v, ok := s.Last()
	if !ok || k != 3 || v != "c" {
		t.Errorf(".Last() = %v, %v, %v, want 3, c, true", k, v, ok)
	}

	// Test Last2 with empty sequence
	empty := FromPairs([]Pair[int, string]{})
	_, _, ok2 := empty.Last()
	if ok2 {
		t.Errorf("empty.Last() ok = %v, want false", ok2)
	}

	// Test Last2 with single pair
	single := FromPairs([]Pair[string, int]{{Key: "world", Value: 99}})
	k3, v3, ok3 := single.Last()
	if !ok3 || k3 != "world" || v3 != 99 {
		t.Errorf("single.Last() = %v, %v, %v, want world, 99, true", k3, v3, ok3)
	}
}

func TestKeys(t *testing.T) {
	// Test keys operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result := FromPairs(pairs).Keys().ToSlice()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Keys() = %v, want %v", result, expected)
	}
}

func TestValues(t *testing.T) {
	// Test values operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result := FromPairs(pairs).Values().ToSlice()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Values() = %v, want %v", result, expected)
	}
}

func TestFromPairs(t *testing.T) {
	// Test fromPairs operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	result := FromPairs(pairs).ToPairs()
	expected := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FromPairs() = %v, want %v", result, expected)
	}
}

func TestFromPairsEarlyTermination(t *testing.T) {
	// Test fromPairs with early termination
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}, {Key: 4, Value: "d"}}

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
	pairs := []Pair[int, string]{{Key: 3, Value: "c"}, {Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	result := FromPairs(pairs).OrderByKey(func(a, b int) bool { return a < b }).ToPairs()
	expected := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".OrderByKey() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := FromPairs([]Pair[int, string]{}).OrderByKey(func(a, b int) bool { return a < b }).ToPairs()
	if len(result2) != 0 {
		t.Errorf("empty.OrderByKey() = %v, want empty", result2)
	}
}

func TestOrderByValue(t *testing.T) {
	// Test orderByValue operation
	pairs := []Pair[int, string]{{Key: 1, Value: "c"}, {Key: 2, Value: "a"}, {Key: 3, Value: "b"}}
	result := FromPairs(pairs).OrderByValue(func(a, b string) bool { return a < b }).ToPairs()
	expected := []Pair[int, string]{{Key: 2, Value: "a"}, {Key: 3, Value: "b"}, {Key: 1, Value: "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".OrderByValue() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := FromPairs([]Pair[int, string]{}).OrderByValue(func(a, b string) bool { return a < b }).ToPairs()
	if len(result2) != 0 {
		t.Errorf("empty.OrderByValue() = %v, want empty", result2)
	}
}

func TestSortBy2(t *testing.T) {
	// Test sortBy2 operation
	pairs := []Pair[string, int]{{Key: "c", Value: 3}, {Key: "a", Value: 1}, {Key: "b", Value: 2}}
	result := FromPairs(pairs).SortBy(func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}).ToPairs()
	expected := []Pair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "c", Value: 3}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".SortBy() = %v, want %v", result, expected)
	}

	// Test with empty sequence
	result2 := FromPairs([]Pair[string, int]{}).SortBy(func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}).ToPairs()
	if len(result2) != 0 {
		t.Errorf("empty.SortBy() = %v, want empty", result2)
	}

	// Test reverse sort by value
	result3 := FromPairs(pairs).SortBy(func(a, b Pair[string, int]) bool {
		return a.Value > b.Value
	}).ToPairs()
	expected3 := []Pair[string, int]{{Key: "c", Value: 3}, {Key: "b", Value: 2}, {Key: "a", Value: 1}}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("reverse by value.SortBy() = %v, want %v", result3, expected3)
	}

	// Test with equal elements (to cover return 0 case)
	equals := []Pair[string, int]{{Key: "a", Value: 1}, {Key: "b", Value: 1}, {Key: "c", Value: 1}}
	result4 := FromPairs(equals).SortBy(func(a, b Pair[string, int]) bool {
		return a.Value < b.Value // All values are equal
	}).ToPairs()
	// Order should remain stable or be in some order
	if len(result4) != 3 {
		t.Errorf("equal values.SortBy() length = %v, want 3", len(result4))
	}

	// Test single element
	single := []Pair[string, int]{{Key: "a", Value: 1}}
	result5 := FromPairs(single).SortBy(func(a, b Pair[string, int]) bool {
		return a.Key < b.Key
	}).ToPairs()
	if !reflect.DeepEqual(result5, single) {
		t.Errorf("single.SortBy() = %v, want %v", result5, single)
	}

	// Test equal elements
	equal := []Pair[string, int]{{Key: "a", Value: 1}, {Key: "a", Value: 2}, {Key: "a", Value: 3}}
	result6 := FromPairs(equal).SortBy(func(a, b Pair[string, int]) bool {
		return a.Value < b.Value
	}).ToPairs()
	expected5 := []Pair[string, int]{{Key: "a", Value: 1}, {Key: "a", Value: 2}, {Key: "a", Value: 3}}
	if !reflect.DeepEqual(result6, expected5) {
		t.Errorf("equal keys.SortBy(sort by value) = %v, want %v", result6, expected5)
	}
}

func TestOrderByKeyAdvanced(t *testing.T) {
	// Test OrderByKey with single element
	single := []Pair[int, string]{{Key: 1, Value: "a"}}
	result := FromPairs(single).OrderByKey(func(a, b int) bool { return a < b }).ToPairs()
	if !reflect.DeepEqual(result, single) {
		t.Errorf("single.OrderByKey() = %v, want %v", result, single)
	}

	// Test OrderByKey with already sorted data
	sorted := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result2 := FromPairs(sorted).OrderByKey(func(a, b int) bool { return a < b }).ToPairs()
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("already sorted.OrderByKey() = %v, want %v", result2, sorted)
	}

	// Test OrderByKey with duplicate keys
	duplicates := []Pair[int, string]{{Key: 2, Value: "b"}, {Key: 1, Value: "a"}, {Key: 2, Value: "c"}, {Key: 1, Value: "d"}}
	result3 := FromPairs(duplicates).OrderByKey(func(a, b int) bool { return a < b }).ToPairs()
	// Should have keys in order 1, 1, 2, 2
	if len(result3) != 4 || result3[0].Key != 1 || result3[1].Key != 1 || result3[2].Key != 2 || result3[3].Key != 2 {
		t.Errorf("duplicates.OrderByKey() = %v, want keys ordered as [1,1,2,2]", result3)
	}
}

func TestOrderByValueAdvanced(t *testing.T) {
	// Test OrderByValue with single element
	single := []Pair[int, string]{{Key: 1, Value: "a"}}
	result := FromPairs(single).OrderByValue(func(a, b string) bool { return a < b }).ToPairs()
	if !reflect.DeepEqual(result, single) {
		t.Errorf("single.OrderByValue() = %v, want %v", result, single)
	}

	// Test OrderByValue with already sorted data
	sorted := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result2 := FromPairs(sorted).OrderByValue(func(a, b string) bool { return a < b }).ToPairs()
	if !reflect.DeepEqual(result2, sorted) {
		t.Errorf("already sorted.OrderByValue() = %v, want %v", result2, sorted)
	}

	// Test OrderByValue with duplicate values
	duplicates := []Pair[int, string]{{Key: 2, Value: "b"}, {Key: 1, Value: "a"}, {Key: 3, Value: "b"}, {Key: 4, Value: "a"}}
	result3 := FromPairs(duplicates).OrderByValue(func(a, b string) bool { return a < b }).ToPairs()
	// Should have values in order a, a, b, b
	if len(result3) != 4 || result3[0].Value != "a" || result3[1].Value != "a" || result3[2].Value != "b" ||
		result3[3].Value != "b" {
		t.Errorf("duplicates.OrderByValue() = %v, want values ordered as [a,a,b,b]", result3)
	}
}

func TestFilterMap2(t *testing.T) {
	s := FromMap(map[int]string{1: "a", 2: "bb", 3: "ccc"})

	result := ToMap(s.FilterMap(func(k int, v string) (Pair[int, string], bool) {
		if len(v) > 1 {
			return Pair[int, string]{Key: k * 10, Value: strings.ToUpper(v)}, true
		}
		return Pair[int, string]{}, false
	}))

	expected := map[int]string{20: "BB", 30: "CCC"}

	if len(result) != len(expected) {
		t.Errorf(".FilterMap() result length = %d, want %d", len(result), len(expected))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf(".FilterMap() result[%d] = %s, want %s", k, result[k], v)
		}
	}

	// Test empty map
	empty := ToMap(FromMap(map[int]string{}).FilterMap(func(k int, v string) (Pair[int, string], bool) {
		return Pair[int, string]{Key: k, Value: v}, true
	}))
	if len(empty) != 0 {
		t.Errorf("empty.FilterMap() = %v, want empty", empty)
	}

	// Test filter all out
	allFiltered := ToMap(s.FilterMap(func(k int, v string) (Pair[int, string], bool) {
		return Pair[int, string]{Key: k, Value: v}, false // Filter all out
	}))
	if len(allFiltered) != 0 {
		t.Errorf("filter all.FilterMap() = %v, want empty", allFiltered)
	}
}
