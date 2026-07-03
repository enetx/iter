package iter_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/iter"
)

func TestFind(t *testing.T) {
	// Test find operation
	result, found := FromSlice([]int{1, 2, 3, 4, 5}).Find(func(x int) bool { return x > 3 })
	if !found || result != 4 {
		t.Errorf(".Find() = %v, %v, want 4, true", result, found)
	}

	// Test find not found
	_, found2 := FromSlice([]int{1, 2, 3}).Find(func(x int) bool { return x > 10 })
	if found2 {
		t.Errorf("not found.Find() found = %v, want false", found2)
	}
}

func TestAny(t *testing.T) {
	// Test any operation
	result := FromSlice([]int{1, 2, 3, 4}).Any(func(x int) bool { return x > 3 })
	if !result {
		t.Errorf(".Any() = %v, want true", result)
	}

	// Test any false
	result2 := FromSlice([]int{1, 2, 3}).Any(func(x int) bool { return x > 10 })
	if result2 {
		t.Errorf("false.Any() = %v, want false", result2)
	}
}

func TestAll(t *testing.T) {
	// Test all operation
	result := FromSlice([]int{1, 2, 3, 4}).All(func(x int) bool { return x > 0 })
	if !result {
		t.Errorf(".All() = %v, want true", result)
	}

	// Test all false
	result2 := FromSlice([]int{1, 2, 3}).All(func(x int) bool { return x > 2 })
	if result2 {
		t.Errorf("false.All() = %v, want false", result2)
	}
}

func TestFold(t *testing.T) {
	// Test fold operation
	result := FromSlice([]int{1, 2, 3, 4}).Fold(0, func(acc, x int) int { return acc + x })
	if result != 10 {
		t.Errorf(".Fold() = %v, want 10", result)
	}

	// Test fold with empty sequence
	result2 := FromSlice([]int{}).Fold(5, func(acc, x int) int { return acc + x })
	if result2 != 5 {
		t.Errorf("empty.Fold() = %v, want 5", result2)
	}
}

func TestReduce(t *testing.T) {
	// Test reduce operation
	result, found := FromSlice([]int{1, 2, 3, 4}).Reduce(func(a, b int) int { return a + b })
	if !found || result != 10 {
		t.Errorf(".Reduce() = %v, %v, want 10, true", result, found)
	}

	// Test reduce with empty sequence
	_, found2 := FromSlice([]int{}).Reduce(func(a, b int) int { return a + b })
	if found2 {
		t.Errorf("empty.Reduce() found = %v, want false", found2)
	}
}

func TestContains(t *testing.T) {
	// Test contains operation
	result := Contains(FromSlice([]int{1, 2, 3, 4}), 3)
	if !result {
		t.Errorf("Contains() = %v, want true", result)
	}

	// Test contains false
	result2 := Contains(FromSlice([]int{1, 2, 3}), 5)
	if result2 {
		t.Errorf("Contains(false) = %v, want false", result2)
	}
}

func TestFind2(t *testing.T) {
	// Test find2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "bb"}, {Key: 3, Value: "c"}}
	k, v, found := FromPairs(pairs).Find(func(key int, val string) bool { return len(val) > 1 })
	if !found || k != 2 || v != "bb" {
		t.Errorf(".Find() = %v, %v, %v, want 2, bb, true", k, v, found)
	}
}

func TestAny2(t *testing.T) {
	// Test any2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "bb"}}
	result := FromPairs(pairs).Any(func(k int, v string) bool { return len(v) > 1 })
	if !result {
		t.Errorf(".Any() = %v, want true", result)
	}
}

func TestAll2(t *testing.T) {
	// Test all2 operation
	pairs := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "b"}}
	result := FromPairs(pairs).All(func(k int, v string) bool { return len(v) == 1 })
	if !result {
		t.Errorf(".All() = %v, want true", result)
	}

	// Test all2 with early false
	pairs2 := []Pair[int, string]{{Key: 1, Value: "a"}, {Key: 2, Value: "bb"}, {Key: 3, Value: "c"}}
	result2 := FromPairs(pairs2).All(func(k int, v string) bool { return len(v) == 1 })
	if result2 {
		t.Errorf("early false.All() = %v, want false", result2)
	}

	// Test all2 with empty sequence
	result3 := FromPairs([]Pair[int, string]{}).All(func(k int, v string) bool { return len(v) == 1 })
	if !result3 {
		t.Errorf("empty.All() = %v, want true", result3)
	}

	// Test all2 early termination
	count := 0
	pairs3 := []Pair[int, string]{{Key: 1, Value: "aa"}, {Key: 2, Value: "bb"}, {Key: 3, Value: "cc"}}
	FromPairs(pairs3).All(func(k int, v string) bool {
		count++
		return false // First check returns false, should terminate
	})
	if count != 1 {
		t.Errorf("All2 early termination count = %v, want 1", count)
	}
}

func TestFold2(t *testing.T) {
	// Test fold2 operation
	pairs := []Pair[int, int]{{Key: 1, Value: 10}, {Key: 2, Value: 20}}
	result := FromPairs(pairs).Fold(0, func(acc, k, v int) int { return acc + k + v })
	if result != 33 { // 0 + 1 + 10 + 2 + 20
		t.Errorf(".Fold() = %v, want 33", result)
	}
}

func TestReduce2(t *testing.T) {
	// Test reduce2 operation
	pairs := []Pair[int, int]{{Key: 1, Value: 10}, {Key: 2, Value: 20}}
	result, found := FromPairs(pairs).Reduce(func(a, b Pair[int, int]) Pair[int, int] {
		return Pair[int, int]{Key: a.Key + b.Key, Value: a.Value + b.Value}
	})
	if !found || result.Key != 3 || result.Value != 30 {
		t.Errorf(".Reduce() = %v, %v, want {3, 30}, true", result, found)
	}
}

func TestMinBy(t *testing.T) {
	// Test basic minBy operation
	result, found := FromSlice([]int{3, 1, 4, 1, 5}).MinBy(func(a, b int) bool { return a < b })
	if !found || result != 1 {
		t.Errorf(".MinBy() = %v, %v, want 1, true", result, found)
	}

	// Test empty sequence
	_, found2 := FromSlice([]int{}).MinBy(func(a, b int) bool { return a < b })
	if found2 {
		t.Errorf("empty.MinBy() found = %v, want false", found2)
	}
}

func TestMaxBy(t *testing.T) {
	// Test basic maxBy operation
	result, found := FromSlice([]int{3, 1, 4, 1, 5}).MaxBy(func(a, b int) bool { return a < b })
	if !found || result != 5 {
		t.Errorf(".MaxBy() = %v, %v, want 5, true", result, found)
	}

	// Test empty sequence
	_, found2 := FromSlice([]int{}).MaxBy(func(a, b int) bool { return a < b })
	if found2 {
		t.Errorf("empty.MaxBy() found = %v, want false", found2)
	}
}

func TestCountBy(t *testing.T) {
	// Test countBy operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).CountBy(func(x int) bool { return x%2 == 0 })
	expected := map[bool]int{true: 2, false: 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".CountBy() = %v, want %v", result, expected)
	}

	// Test empty sequence
	result2 := FromSlice([]int{}).CountBy(func(x int) bool { return x%2 == 0 })
	if len(result2) != 0 {
		t.Errorf("empty.CountBy() = %v, want empty map", result2)
	}
}

func TestPartition(t *testing.T) {
	// Test partition operation
	left, right := FromSlice([]int{1, 2, 3, 4, 5}).Partition(func(x int) bool { return x%2 == 0 })

	expectedLeft := []int{2, 4}
	expectedRight := []int{1, 3, 5}

	if !reflect.DeepEqual(left, expectedLeft) {
		t.Errorf(".Partition() left = %v, want %v", left, expectedLeft)
	}
	if !reflect.DeepEqual(right, expectedRight) {
		t.Errorf(".Partition() right = %v, want %v", right, expectedRight)
	}
}

func TestSortBy(t *testing.T) {
	// Test sortBy operation
	result := FromSlice([]int{3, 1, 4, 1, 5}).SortBy(func(a, b int) bool { return a < b }).ToSlice()
	expected := []int{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".SortBy() = %v, want %v", result, expected)
	}

	// Test reverse sort
	result2 := FromSlice([]int{3, 1, 4, 1, 5}).SortBy(func(a, b int) bool { return a > b }).ToSlice()
	expected2 := []int{5, 4, 3, 1, 1}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("reverse.SortBy() = %v, want %v", result2, expected2)
	}
}

func TestPosition(t *testing.T) {
	// Test position operation
	result, found := FromSlice([]int{1, 2, 3, 4, 5}).Position(func(x int) bool { return x == 3 })
	if !found || result != 2 {
		t.Errorf(".Position() = %v, %v, want 2, true", result, found)
	}

	// Test not found
	_, found2 := FromSlice([]int{1, 2, 3, 4, 5}).Position(func(x int) bool { return x == 10 })
	if found2 {
		t.Errorf("not found.Position() found = %v, want false", found2)
	}
}

func TestRPosition(t *testing.T) {
	// Test rposition operation
	result, found := FromSlice([]int{1, 2, 3, 2, 5}).RPosition(func(x int) bool { return x == 2 })
	if !found || result != 3 {
		t.Errorf(".RPosition() = %v, %v, want 3, true", result, found)
	}

	// Test not found
	_, found2 := FromSlice([]int{1, 2, 3, 4, 5}).RPosition(func(x int) bool { return x == 10 })
	if found2 {
		t.Errorf("not found.RPosition() found = %v, want false", found2)
	}
}

func TestIsPartitioned(t *testing.T) {
	// Test is_partitioned operation
	result := FromSlice([]int{2, 4, 1, 3, 5}).IsPartitioned(func(x int) bool { return x%2 == 0 })
	if !result {
		t.Errorf(".IsPartitioned() = %v, want true", result)
	}

	// Test not partitioned
	result2 := FromSlice([]int{1, 2, 3, 4, 5}).IsPartitioned(func(x int) bool { return x%2 == 0 })
	if result2 {
		t.Errorf("not partitioned.IsPartitioned() = %v, want false", result2)
	}

	// Test empty sequence
	result3 := FromSlice([]int{}).IsPartitioned(func(x int) bool { return x%2 == 0 })
	if !result3 {
		t.Errorf("empty.IsPartitioned() = %v, want true", result3)
	}
}
