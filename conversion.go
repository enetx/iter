package iter

import (
	"context"
	"iter"
)

// Pull converts a push-style iterator (Seq) to a pull-style iterator.
// Returns a next function that yields the next value and a boolean indicating if valid,
// and a stop function that should be called to release resources.
//
// Example:
//
//	next, stop := iter.FromSlice([]int{1, 2, 3}).Pull()
//	defer stop()
//	for {
//	  v, ok := next()
//	  if !ok { break }
//	  fmt.Println(v)
//	}
func (s Seq[T]) Pull() (next func() (T, bool), stop func()) {
	return iter.Pull(iter.Seq[T](s))
}

// Pull converts a push-style iterator (Seq2) to a pull-style iterator.
func (s Seq2[K, V]) Pull() (next func() (K, V, bool), stop func()) {
	return iter.Pull2(iter.Seq2[K, V](s))
}

// Context wraps a sequence with context cancellation.
// If the context is cancelled, iteration stops early.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	s := iter.FromSlice([]int{1, 2, 3}).Context(ctx)
func (s Seq[T]) Context(ctx context.Context) Seq[T] {
	return func(yield func(T) bool) {
		if err := ctx.Err(); err != nil {
			return
		}
		s(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			default:
				return yield(v)
			}
		})
	}
}

// Context wraps a key-value sequence with context cancellation.
// If the context is cancelled, iteration stops early.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	s := iter.FromMap(map[int]string{1: "a", 2: "b"}).Context(ctx)
func (s Seq2[K, V]) Context(ctx context.Context) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if err := ctx.Err(); err != nil {
			return
		}
		s(func(k K, v V) bool {
			select {
			case <-ctx.Done():
				return false
			default:
				return yield(k, v)
			}
		})
	}
}

// ToChan converts a sequence to a channel.
// The channel is closed when the sequence is exhausted or context is cancelled.
//
// Example:
//
//	ctx := context.Background()
//	ch := iter.FromSlice([]int{1, 2, 3}).ToChan(ctx)
//	for v := range ch {
//	  fmt.Println(v)
//	}
func (s Seq[T]) ToChan(ctx context.Context) chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if err := ctx.Err(); err != nil {
			return
		}
		s(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
				return true
			}
		})
	}()
	return ch
}

// ToChan converts a key-value sequence to a channel of Pair pairs.
// The channel is closed when the sequence is exhausted or context is cancelled.
//
// Example:
//
//	ctx := context.Background()
//	ch := iter.FromMap(map[int]string{1: "a", 2: "b"}).ToChan(ctx)
//	for kv := range ch {
//	  fmt.Printf("%d: %s\n", kv.Key, kv.Value)
//	}
func (s Seq2[K, V]) ToChan(ctx context.Context) chan Pair[K, V] {
	ch := make(chan Pair[K, V])
	go func() {
		defer close(ch)
		if err := ctx.Err(); err != nil {
			return
		}
		s(func(k K, v V) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- Pair[K, V]{k, v}:
				return true
			}
		})
	}()
	return ch
}

// ToSlice collects all elements from the sequence into a slice.
//
// Example:
//
//	sl := iter.FromSlice([]int{1, 2, 3}).ToSlice() // [1, 2, 3]
func (s Seq[T]) ToSlice() []T {
	out := make([]T, 0)
	s(func(v T) bool {
		out = append(out, v)
		return true
	})
	return out
}

// ToMap collects all key-value pairs into a map.
// Later pairs with the same key will overwrite earlier ones.
// It remains a free function because it requires K to be comparable,
// which cannot be expressed as a constraint on the receiver's type parameter.
func ToMap[K comparable, V any](s Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	s(func(k K, v V) bool {
		m[k] = v
		return true
	})
	return m
}

// ToPairs collects all key-value pairs into a slice of Pair structs.
func (s Seq2[K, V]) ToPairs() []Pair[K, V] {
	out := make([]Pair[K, V], 0)
	s(func(k K, v V) bool {
		out = append(out, Pair[K, V]{k, v})
		return true
	})
	return out
}
