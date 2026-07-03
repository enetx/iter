package iter_test

import (
	"strings"
	"testing"

	. "github.com/enetx/iter"
)

func TestCmp(t *testing.T) {
	// Test equal sequences
	cmp1 := FromSlice([]int{1, 2, 3}).Cmp(FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != 0 {
		t.Errorf("equal.Cmp() = %d, want 0", cmp1)
	}

	// Test first less than second
	cmp2 := FromSlice([]int{1, 2}).Cmp(FromSlice([]int{1, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp2 != -1 {
		t.Errorf("first < second.Cmp() = %d, want -1", cmp2)
	}

	// Test first greater than second
	cmp3 := FromSlice([]int{1, 3}).Cmp(FromSlice([]int{1, 2}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp3 != 1 {
		t.Errorf("first > second.Cmp() = %d, want 1", cmp3)
	}
}

func TestEqual(t *testing.T) {
	// Test equal sequences
	result := FromSlice([]int{1, 2, 3}).Equal(FromSlice([]int{1, 2, 3}))
	if !result {
		t.Errorf(".Equal() = %v, want true", result)
	}

	// Test unequal sequences
	result2 := FromSlice([]int{1, 2, 3}).Equal(FromSlice([]int{1, 2, 4}))
	if result2 {
		t.Errorf("unequal.Equal() = %v, want false", result2)
	}

	// Test different lengths
	result3 := FromSlice([]int{1, 2}).Equal(FromSlice([]int{1, 2, 3}))
	if result3 {
		t.Errorf("different lengths.Equal() = %v, want false", result3)
	}

	// Test empty sequences
	result4 := FromSlice([]int{}).Equal(FromSlice([]int{}))
	if !result4 {
		t.Errorf("empty.Equal(empty) = %v, want true", result4)
	}

	// Test one empty, one not
	result5 := FromSlice([]int{}).Equal(FromSlice([]int{1}))
	if result5 {
		t.Errorf("empty.Equal(non-empty) = %v, want false", result5)
	}

	// Test non-comparable types (slices)
	slices1 := [][]int{{1, 2}, {3, 4}}
	slices2 := [][]int{{1, 2}, {3, 4}}
	result6 := FromSlice(slices1).Equal(FromSlice(slices2))
	if !result6 {
		t.Errorf("non-comparable equal.Equal() = %v, want true", result6)
	}

	// Test non-comparable unequal types
	slices3 := [][]int{{1, 2}, {3, 5}}
	result7 := FromSlice(slices1).Equal(FromSlice(slices3))
	if result7 {
		t.Errorf("non-comparable unequal.Equal() = %v, want false", result7)
	}
}

func TestLt(t *testing.T) {
	// Test first less than second
	result := FromSlice([]int{1, 2}).Lt(FromSlice([]int{1, 3}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf(".Lt() = %v, want true", result)
	}

	// Test first not less than second
	result2 := FromSlice([]int{1, 3}).Lt(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("not less.Lt() = %v, want false", result2)
	}

	// Test first shorter than second
	result3 := FromSlice([]int{1, 2}).Lt(FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("first shorter.Lt() = %v, want true", result3)
	}
}

func TestLe(t *testing.T) {
	// Test first less than or equal to second
	result := FromSlice([]int{1, 2}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf(".Le() = %v, want true", result)
	}

	// Test first longer than second (should be false)
	result2 := FromSlice([]int{1, 2, 3}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("first longer.Le() = %v, want false", result2)
	}

	// Test first element greater than second (less(bv, av) case)
	result3 := FromSlice([]int{3, 2}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result3 {
		t.Errorf("first > second.Le() = %v, want false", result3)
	}
}

func TestGt(t *testing.T) {
	// Test first greater than second
	result := FromSlice([]int{1, 3}).Gt(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf(".Gt() = %v, want true", result)
	}

	// Test first not greater than second
	result2 := FromSlice([]int{1, 2}).Gt(FromSlice([]int{1, 3}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("not greater.Gt() = %v, want false", result2)
	}

	// Test first longer than second
	result3 := FromSlice([]int{1, 2, 3}).Gt(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("first longer.Gt() = %v, want true", result3)
	}
}

func TestGe(t *testing.T) {
	// Test first greater than or equal to second
	result := FromSlice([]int{1, 2}).Ge(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf(".Ge() = %v, want true", result)
	}

	// Test first shorter than second (should be false)
	result2 := FromSlice([]int{1, 2}).Ge(FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("first shorter.Ge() = %v, want false", result2)
	}

	// Test first element less than second (less(av, bv) case)
	result3 := FromSlice([]int{1, 2}).Ge(FromSlice([]int{3, 2}), func(a, b int) bool {
		return a < b
	})
	if result3 {
		t.Errorf("first < second.Ge() = %v, want false", result3)
	}
}

func TestCmpDifferentLengths(t *testing.T) {
	// Test comparison with different lengths
	cmp1 := FromSlice([]int{1, 2}).Cmp(FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != -1 {
		t.Errorf("different lengths.Cmp(first shorter) = %d, want -1", cmp1)
	}
}

func TestCmpSecondLonger(t *testing.T) {
	// Test comparison when second is longer
	cmp1 := FromSlice([]int{1, 2, 3}).Cmp(FromSlice([]int{1, 2}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != 1 {
		t.Errorf("second shorter.Cmp() = %d, want 1", cmp1)
	}
}

func TestLtWithEqualElements(t *testing.T) {
	// Test Lt with equal elements
	result := FromSlice([]int{1, 2}).Lt(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result {
		t.Errorf("equal.Lt() = %v, want false", result)
	}
}

func TestLeWithEqualElements(t *testing.T) {
	// Test Le with equal elements
	result := FromSlice([]int{1, 2}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("equal.Le() = %v, want true", result)
	}
}

func TestGtWithEqualElements(t *testing.T) {
	// Test Gt with equal elements
	result := FromSlice([]int{1, 2}).Gt(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result {
		t.Errorf("equal.Gt() = %v, want false", result)
	}
}

func TestGeWithEqualElements(t *testing.T) {
	// Test Ge with equal elements
	result := FromSlice([]int{1, 2}).Ge(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("equal.Ge() = %v, want true", result)
	}
}

func TestEqualBy(t *testing.T) {
	// Test EqualBy with custom comparison by length
	result := FromSlice([]string{"a", "bb"}).EqualBy(FromSlice([]string{"x", "yy"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if !result {
		t.Errorf("same lengths.EqualBy() = %v, want true", result)
	}

	// Test EqualBy with different lengths
	result2 := FromSlice([]string{"a", "bb"}).EqualBy(FromSlice([]string{"x", "yyy"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if result2 {
		t.Errorf("different lengths.EqualBy() = %v, want false", result2)
	}

	// Test EqualBy early termination
	result3 := FromSlice([]string{"a", "bb", "ccc"}).EqualBy(
		FromSlice([]string{"x", "yyy", "dddd"}),
		func(a, b string) bool {
			return len(a) == len(b)
		},
	)
	if result3 {
		t.Errorf("early mismatch.EqualBy() = %v, want false", result3)
	}

	// Test EqualBy different sequence lengths
	result4 := FromSlice([]string{"a", "bb"}).EqualBy(FromSlice([]string{"x", "yy", "zzz"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if result4 {
		t.Errorf("different seq lengths.EqualBy() = %v, want false", result4)
	}

	// Test EqualBy empty sequences
	result5 := FromSlice([]string{}).EqualBy(FromSlice([]string{}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if !result5 {
		t.Errorf("empty.EqualBy() = %v, want true", result5)
	}

	// Test EqualBy case-insensitive comparison
	result6 := FromSlice([]string{"Hello", "WORLD"}).EqualBy(
		FromSlice([]string{"hello", "world"}),
		func(a, b string) bool {
			return len(a) == len(b) && strings.EqualFold(a, b)
		},
	)
	if !result6 {
		t.Errorf("case insensitive.EqualBy() = %v, want true", result6)
	}

	// Test EqualBy with early false return
	count := 0
	result7 := FromSlice([]string{"a", "bb", "ccc"}).EqualBy(
		FromSlice([]string{"x", "yy", "ddd"}),
		func(a, b string) bool {
			count++
			return len(a) == len(b)
		},
	)
	if !result7 {
		t.Errorf("all match.EqualBy() = %v, want true", result7)
	}
	if count != 3 {
		t.Errorf("EqualBy comparison count = %v, want 3", count)
	}
}

func TestLeAdvanced(t *testing.T) {
	// Test Le with different length sequences (first shorter)
	result := FromSlice([]int{1, 2}).Le(FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("first shorter.Le() = %v, want true", result)
	}

	// Test Le with different length sequences (second shorter)
	result2 := FromSlice([]int{1, 2, 3}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("second shorter.Le() = %v, want false", result2)
	}

	// Test Le with first less than second
	result3 := FromSlice([]int{1, 1}).Le(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("first < second.Le() = %v, want true", result3)
	}

	// Test Le with empty sequences
	result4 := FromSlice([]int{}).Le(FromSlice([]int{}), func(a, b int) bool {
		return a < b
	})
	if !result4 {
		t.Errorf("empty sequences.Le() = %v, want true", result4)
	}
}

func TestGeAdvanced(t *testing.T) {
	// Test Ge with different length sequences (first longer)
	result := FromSlice([]int{1, 2, 3}).Ge(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("first longer.Ge() = %v, want true", result)
	}

	// Test Ge with different length sequences (second longer)
	result2 := FromSlice([]int{1, 2}).Ge(FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("second longer.Ge() = %v, want false", result2)
	}

	// Test Ge with first greater than second
	result3 := FromSlice([]int{1, 3}).Ge(FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("first > second.Ge() = %v, want true", result3)
	}

	// Test Ge with empty sequences
	result4 := FromSlice([]int{}).Ge(FromSlice([]int{}), func(a, b int) bool {
		return a < b
	})
	if !result4 {
		t.Errorf("empty sequences.Ge() = %v, want true", result4)
	}
}
