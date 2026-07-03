// Package iter provides lazy, composable iterator sequences built on the
// standard library's iter.Seq and iter.Seq2.
//
// # Core types
//
// [Seq] is a single-value sequence and [Seq2] is a two-value (typically
// key-value) sequence. Both are thin named types over the standard library's
// iterator function types, so they interoperate directly with range-over-func
// and the std iter package.
//
// # Laziness
//
// All transformations (Map, Filter, Take, Zip, ...) are lazy: they wrap the
// upstream sequence and do no work until a consumer such as ForEach, Fold,
// ToSlice or a range loop actually pulls elements. Building a chain allocates
// closures only; elements flow through one at a time.
//
// # Generic methods (Go 1.27)
//
// Transformations that change the element type are methods with their own
// type parameters, so chains can switch types mid-stream:
//
//	iter.FromSlice([]int{1, 2, 3}).Map[string](strconv.Itoa).ToSlice()
//
// # Free functions and why they exist
//
// A few operations remain package-level functions instead of methods, for one
// of three reasons:
//
//   - Constructors and generators have no receiver to hang off:
//     [FromSlice], [FromMap], [FromChan], [Iota], [Once], [Repeat], [Empty], ...
//   - Some operations need extra constraints on the receiver's existing type
//     parameters (e.g. comparable), which Go cannot express on a method:
//     [Dedup], [Contains], [Counter], [ToMap].
//   - Some operations would have to instantiate Seq with a type that contains
//     the receiver's own type parameter (Seq[[]T] or Seq[Seq[T]]), which
//     creates a generic instantiation cycle: [Windows], [Chunks],
//     [GroupByAdjacent], [Combinations], [Permutations], [Flatten], [FlattenSeq].
package iter
