package iter_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/iter"
)

func TestNext(t *testing.T) {
	// Test Next with multiple elements
	s := FromSlice([]int{1, 2, 3, 4, 5})

	// Get first element
	val1, rest1, ok1 := s.Next()
	if !ok1 || val1 != 1 {
		t.Errorf(".Next() first = %v, %v, want 1, true", val1, ok1)
	}

	// Get second element
	val2, rest2, ok2 := rest1.Next()
	if !ok2 || val2 != 2 {
		t.Errorf(".Next() second = %v, %v, want 2, true", val2, ok2)
	}

	// Check remaining elements
	remaining := rest2.ToSlice()
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(remaining, expected) {
		t.Errorf(".Next() remaining = %v, want %v", remaining, expected)
	}

	// Test Next with single element
	single := FromSlice([]int{42})
	val, rest, ok := single.Next()
	if !ok || val != 42 {
		t.Errorf(".Next() single = %v, %v, want 42, true", val, ok)
	}
	if rest != nil {
		restSlice := rest.ToSlice()
		if len(restSlice) != 0 {
			t.Errorf(".Next() single rest = %v, want empty", restSlice)
		}
	}

	// Test Next with empty sequence
	empty := FromSlice([]int{})
	val, rest, ok = empty.Next()
	if ok || val != 0 || rest != nil {
		t.Errorf(".Next() empty = %v, %v, %v, want 0, nil, false", val, rest, ok)
	}
}

func TestFirst(t *testing.T) {
	// Test First with non-empty sequence
	result, ok := FromSlice([]int{1, 2, 3, 4, 5}).First()
	if !ok || result != 1 {
		t.Errorf(".First() = %v, %v, want 1, true", result, ok)
	}

	// Test First with empty sequence
	_, ok2 := FromSlice([]int{}).First()
	if ok2 {
		t.Errorf("empty.First() ok = %v, want false", ok2)
	}

	// Test First with single element
	result3, ok3 := FromSlice([]string{"hello"}).First()
	if !ok3 || result3 != "hello" {
		t.Errorf("single.First() = %v, %v, want hello, true", result3, ok3)
	}
}

func TestLast(t *testing.T) {
	// Test Last with non-empty sequence
	result, ok := FromSlice([]int{1, 2, 3, 4, 5}).Last()
	if !ok || result != 5 {
		t.Errorf(".Last() = %v, %v, want 5, true", result, ok)
	}

	// Test Last with empty sequence
	_, ok2 := FromSlice([]int{}).Last()
	if ok2 {
		t.Errorf("empty.Last() ok = %v, want false", ok2)
	}

	// Test Last with single element
	result3, ok3 := FromSlice([]string{"world"}).Last()
	if !ok3 || result3 != "world" {
		t.Errorf("single.Last() = %v, %v, want world, true", result3, ok3)
	}
}

func TestForEach(t *testing.T) {
	// Test forEach operation
	var sum int
	FromSlice([]int{1, 2, 3}).ForEach(func(x int) { sum += x })
	if sum != 6 {
		t.Errorf(".ForEach() sum = %v, want 6", sum)
	}
}

func TestCount(t *testing.T) {
	// Test count operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).Count()
	if result != 5 {
		t.Errorf(".Count() = %v, want 5", result)
	}

	// Test count with empty sequence
	result2 := FromSlice([]int{}).Count()
	if result2 != 0 {
		t.Errorf("empty.Count() = %v, want 0", result2)
	}
}

func TestRange(t *testing.T) {
	// Test range operation
	var sum int
	FromSlice([]int{1, 2, 3, 4, 5}).Range(func(x int) bool {
		sum += x
		return x < 3
	})
	if sum != 6 { // 1 + 2 + 3
		t.Errorf(".Range() sum = %v, want 6", sum)
	}
}

func TestTake(t *testing.T) {
	// Test take operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).Take(3).ToSlice()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Take() = %v, want %v", result, expected)
	}

	// Test take with more than available
	result2 := FromSlice([]int{1, 2}).Take(5).ToSlice()
	expected2 := []int{1, 2}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("more than available.Take() = %v, want %v", result2, expected2)
	}
}

func TestSkip(t *testing.T) {
	// Test skip operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice()
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Skip() = %v, want %v", result, expected)
	}

	// Test skip with more than available
	result2 := FromSlice([]int{1, 2}).Skip(5).ToSlice()
	if len(result2) != 0 {
		t.Errorf("more than available.Skip() = %v, want empty slice", result2)
	}
}

func TestStepBy(t *testing.T) {
	// Test stepBy operation
	result := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8}).StepBy(3).ToSlice()
	expected := []int{1, 4, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".StepBy() = %v, want %v", result, expected)
	}
}

func TestTakeWhile(t *testing.T) {
	// Test takeWhile operation
	result := FromSlice([]int{1, 2, 3, 4, 2, 1}).TakeWhile(func(x int) bool { return x < 4 }).ToSlice()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".TakeWhile() = %v, want %v", result, expected)
	}
}

func TestSkipWhile(t *testing.T) {
	// Test skipWhile operation
	result := FromSlice([]int{1, 2, 3, 4, 5}).SkipWhile(func(x int) bool { return x < 3 }).ToSlice()
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".SkipWhile() = %v, want %v", result, expected)
	}
}

func TestNth(t *testing.T) {
	// Test nth operation
	result, found := FromSlice([]int{10, 20, 30, 40}).Nth(2)
	if !found || result != 30 {
		t.Errorf(".Nth() = %v, %v, want 30, true", result, found)
	}

	// Test nth out of bounds
	_, found2 := FromSlice([]int{1, 2}).Nth(5)
	if found2 {
		t.Errorf("out of bounds.Nth() found = %v, want false", found2)
	}
}

func TestForEach2(t *testing.T) {
	// Test forEach2 operation
	var sum int
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	FromPairs(pairs).ForEach(func(k int, v string) { sum += k })
	if sum != 3 {
		t.Errorf(".ForEach() sum = %v, want 3", sum)
	}
}

func TestCount2(t *testing.T) {
	// Test count2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := FromPairs(pairs).Count()
	if result != 2 {
		t.Errorf(".Count() = %v, want 2", result)
	}
}

func TestRange2(t *testing.T) {
	// Test range2 operation
	var sum int
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	FromPairs(pairs).Range(func(k int, v string) bool {
		sum += k
		return k < 2
	})
	if sum != 3 { // 1 + 2
		t.Errorf(".Range() sum = %v, want 3", sum)
	}
}

func TestTake2(t *testing.T) {
	// Test take2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := FromPairs(pairs).Take(2).ToPairs()
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Take() = %v, want %v", result, expected)
	}
}

func TestSkip2(t *testing.T) {
	// Test skip2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	result := FromPairs(pairs).Skip(1).ToPairs()
	expected := []Pair[int, string]{{2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".Skip() = %v, want %v", result, expected)
	}
}

func TestStepBy2(t *testing.T) {
	// Test stepBy2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	result := FromPairs(pairs).StepBy(2).ToPairs()
	expected := []Pair[int, string]{{1, "a"}, {3, "c"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf(".StepBy() = %v, want %v", result, expected)
	}
}

func TestNth2(t *testing.T) {
	// Test nth2 operation
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}
	k, v, found := FromPairs(pairs).Nth(1)
	if !found || k != 2 || v != "b" {
		t.Errorf(".Nth() = %v, %v, %v, want 2, b, true", k, v, found)
	}
}

// Edge cases tests

func TestSkipNegative(t *testing.T) {
	// Test skip with negative value
	result := FromSlice([]int{1, 2, 3}).Skip(-1).ToSlice()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("negative.Skip() = %v, want %v", result, expected)
	}
}

func TestStepByNegative(t *testing.T) {
	// Test stepBy with invalid step
	result := FromSlice([]int{1, 2, 3}).StepBy(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("0.StepBy() = %v, want empty slice", result)
	}
}

func TestTakeZero(t *testing.T) {
	// Test take with zero
	result := FromSlice([]int{1, 2, 3}).Take(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("0.Take() = %v, want empty slice", result)
	}
}

func TestTake2Zero(t *testing.T) {
	// Test take2 with zero
	pairs := []Pair[int, string]{{1, "a"}}
	result := FromPairs(pairs).Take(0).ToPairs()
	if len(result) != 0 {
		t.Errorf("0.Take() = %v, want empty slice", result)
	}
}

func TestSkipZero(t *testing.T) {
	// Test skip with zero
	result := FromSlice([]int{1, 2, 3}).Skip(0).ToSlice()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("0.Skip() = %v, want %v", result, expected)
	}
}

func TestSkip2Zero(t *testing.T) {
	// Test skip2 with zero
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := FromPairs(pairs).Skip(0).ToPairs()
	expected := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("0.Skip() = %v, want %v", result, expected)
	}
}

func TestStepByZero(t *testing.T) {
	// Test stepBy with zero step
	result := FromSlice([]int{1, 2, 3}).StepBy(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("0.StepBy() = %v, want empty slice", result)
	}
}

func TestStepBy2Zero(t *testing.T) {
	// Test stepBy2 with zero step
	pairs := []Pair[int, string]{{1, "a"}}
	result := FromPairs(pairs).StepBy(0).ToPairs()
	if len(result) != 0 {
		t.Errorf("0.StepBy() = %v, want empty slice", result)
	}
}

func TestStepByZeroStepEdgeCase(t *testing.T) {
	// Test stepBy with zero step edge case
	result := FromSlice([]int{1, 2, 3, 4, 5}).StepBy(0).ToSlice()
	if len(result) != 0 {
		t.Errorf("StepBy zero step = %v, want empty slice", result)
	}
}

func TestStepBy2ZeroStepEdgeCase(t *testing.T) {
	// Test stepBy2 with zero step edge case
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}}
	result := FromPairs(pairs).StepBy(0).ToPairs()
	if len(result) != 0 {
		t.Errorf("StepBy2 zero step = %v, want empty slice", result)
	}
}

func TestStepByEarlyTermination(t *testing.T) {
	// Test stepBy with early termination
	count := 0
	FromSlice([]int{1, 2, 3, 4, 5, 6}).StepBy(2)(func(x int) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("StepByEarlyTermination count = %v, want 2", count)
	}
}

func TestStepBy2EarlyTermination(t *testing.T) {
	// Test stepBy2 with early termination
	count := 0
	pairs := []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}}
	FromPairs(pairs).StepBy(2)(func(k int, v string) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("StepBy2EarlyTermination count = %v, want 2", count)
	}
}
