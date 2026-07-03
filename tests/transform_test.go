package iter_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/enetx/iter"
)

func TestScan(t *testing.T) {
	// Test basic scan operation
	result := FromSlice([]int{1, 2, 3}).Scan(0, func(acc, x int) int {
		return acc + x
	}).ToSlice()
	expected := []int{1, 3, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Scan() = %v, want %v", result, expected)
	}

	// Test scan with empty sequence
	result2 := FromSlice([]int{}).Scan(10, func(acc, x int) int {
		return acc + x
	}).ToSlice()
	if len(result2) != 0 {
		t.Errorf("empty.Scan() = %v, want empty slice", result2)
	}
}

func TestMap(t *testing.T) {
	// Test basic map operation
	result := FromSlice([]int{1, 2, 3}).Map(func(x int) int { return x * 2 }).ToSlice()
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Map() = %v, want %v", result, expected)
	}
}

func TestMapTo(t *testing.T) {
	// Test map to different type
	result := FromSlice([]int{1, 2, 3}).Map(func(x int) string { return fmt.Sprintf("%d", x) }).ToSlice()
	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Map() = %v, want %v", result, expected)
	}
}

func TestInspect(t *testing.T) {
	// Test inspect operation
	var sum int
	FromSlice([]int{1, 2, 3}).Inspect(func(x int) { sum += x }).ToSlice()
	if sum != 6 {
		t.Errorf(".Inspect() sum = %v, want 6", sum)
	}
}

func TestFilter(t *testing.T) {
	// Test basic filter operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 }).ToSlice()
	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Filter() = %v, want %v", result, expected)
	}
}

func TestExclude(t *testing.T) {
	// Test basic exclude operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).Exclude(func(x int) bool { return x%2 == 0 }).ToSlice()
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Exclude() = %v, want %v", result, expected)
	}
}

func TestFilterMap(t *testing.T) {
	// Test filter map operation
	result := FromSlice([]string{"1", "2", "abc", "3"}).FilterMap(func(s string) (int, bool) {
		if len(s) == 1 {
			return len(s), true
		}
		return 0, false
	}).ToSlice()
	expected := []int{1, 1, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".FilterMap() = %v, want %v", result, expected)
	}
}

func TestEnumerate(t *testing.T) {
	// Test enumerate operation
	result := FromSlice([]string{"a", "b", "c"}).Enumerate(0).ToPairs()
	expected := []Pair[int, string]{{Key: 0, Value: "a"}, {Key: 1, Value: "b"}, {Key: 2, Value: "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Enumerate() = %v, want %v", result, expected)
	}
}

func TestUnique(t *testing.T) {
	// Test unique operation
	result := FromSlice([]int{1, 2, 2, 3, 3, 3, 4}).Unique().ToSlice()
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Unique() = %v, want %v", result, expected)
	}
}

func TestUniqueBy(t *testing.T) {
	// Test unique by operation - consecutive elements with same key are removed
	result := FromSlice([]string{"a", "b", "bb", "ccc", "dd", "e"}).UniqueBy(func(s string) int { return len(s) }).ToSlice()
	expected := []string{"a", "bb", "ccc", "dd", "e"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".UniqueBy() = %v, want %v", result, expected)
	}
}

func TestMap2(t *testing.T) {
	// Test map2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	result := FromPairs(pairs).Map(func(k int, v string) (int, string) {
		return k * 2, v + v
	}).ToPairs()
	expected := []Pair[int, string]{{Key: 2, Value: "aa"}, {Key: 4, Value: "bb"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Map() = %v, want %v", result, expected)
	}
}

func TestFilter2(t *testing.T) {
	// Test filter2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result := FromPairs(pairs).Filter(func(k int, v string) bool { return k%2 == 0 }).ToPairs()
	expected := []Pair[int, string]{{Key: 2, Value: "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Filter() = %v, want %v", result, expected)
	}
}

func TestExclude2(t *testing.T) {
	// Test exclude2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}, {Key: 3, Value: "c"}}
	result := FromPairs(pairs).Exclude(func(k int, v string) bool { return k%2 == 0 }).ToPairs()
	expected := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 3, Value: "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Exclude() = %v, want %v", result, expected)
	}
}

func TestInspect2(t *testing.T) {
	// Test inspect2 operation
	var sum int
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	FromPairs(pairs).Inspect(func(k int, v string) { sum += k }).ToPairs()
	if sum != 3 {
		t.Errorf(".Inspect() sum = %v, want 3", sum)
	}
}

func TestEnumerateNonZeroStart(t *testing.T) {
	// Test enumerate with non-zero start
	result := FromSlice([]string{"a", "b"}).Enumerate(5).ToPairs()
	expected := []Pair[int, string]{{Key: 5, Value: "a"}, {Key: 6, Value: "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("start=5.Enumerate() = %v, want %v", result, expected)
	}
}

func TestMapWhile(t *testing.T) {
	// Test mapWhile operation
	result := FromSlice([]string{"1", "2", "abc", "3"}).MapWhile(func(s string) (int, bool) {
		if len(s) == 1 {
			return len(s), true
		}
		return 0, false // Stop at "abc"
	}).ToSlice()
	expected := []int{1, 1} // Should stop at "abc"
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".MapWhile() = %v, want %v", result, expected)
	}

	// Test mapWhile with all matching elements
	result2 := FromSlice([]string{"a", "b", "c"}).MapWhile(func(s string) (int, bool) {
		return len(s), true
	}).ToSlice()
	expected2 := []int{1, 1, 1}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("all match.MapWhile() = %v, want %v", result2, expected2)
	}

	// Test mapWhile with empty sequence
	result3 := FromSlice([]string{}).MapWhile(func(s string) (int, bool) {
		return len(s), true
	}).ToSlice()
	if len(result3) != 0 {
		t.Errorf("empty.MapWhile() = %v, want empty", result3)
	}

	// Test mapWhile with no matching elements
	result4 := FromSlice([]string{"abc", "def"}).MapWhile(func(s string) (int, bool) {
		return len(s), len(s) == 1 // None match
	}).ToSlice()
	if len(result4) != 0 {
		t.Errorf("no match.MapWhile() = %v, want empty", result4)
	}
}
