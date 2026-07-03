# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking:** requires **Go 1.27** (currently `1.27.0-rc.1`): the API is now
  method-based. Package-level functions (`Map`, `Filter`, `Fold`, `Take`, …)
  became methods on `Seq[T]`, using Go 1.27 generic methods (`Map[U]`,
  `FilterMap[U]`, `MapWhile[U]`, `Scan[S]`, `Fold[A]`, …) where the element
  type changes mid-chain. `MapTo` is folded into the generic `Seq.Map[U]`.
- **Breaking:** `Seq2` operations are methods without the `*2` suffix:
  `Map2` → `Seq2.Map`, `Filter2` → `Seq2.Filter`, `Fold2` → `Seq2.Fold[A]`,
  `Next2` → `Seq2.Next`, and so on.
- **Breaking:** `OrderByKey`/`OrderByValue` are renamed to
  `Seq2.SortByKey`/`Seq2.SortByValue`.
- **Breaking:** `UniqueBy` is renamed to `DedupByKey[K]` — the name now matches
  its semantics (it deduplicates *consecutive* runs only).
- **Breaking:** the free function `Unique` requires `T comparable` and uses a
  typed `map[T]struct{}` (no more `any`-boxing). For non-comparable element
  types use the new `Seq.Unique` method, which keeps the reflection-free
  `any`-keyed fallback. Operations that need `T comparable` (`Unique`, `Dedup`,
  `Counter`, …) remain free functions, since that constraint cannot be
  expressed on a method receiver.

### Added

- `Seq2.TakeWhile`, `Seq2.SkipWhile`, `Seq2.MapWhile[K2, V2]`, plus generic
  `Seq2.Map[K2, V2]`, `Seq2.FilterMap[K2, V2]`, and `Seq2.Fold[A]`.
- `Seq.FlatMap[U]`.
- `Seq.Unique` method — `any`-keyed fallback for non-comparable element types
  (complements the typed free function `Unique[T comparable]`).

### Fixed

- `Seq.Cycle`: no longer probes the source with an extra up-front pass (which
  consumed an element of single-use sources), and terminates once the source
  stops yielding elements instead of spinning forever.
- `Seq.Next`/`Seq2.Next`: now advance via a pull iterator, so the source is
  walked exactly once — O(1) per element and correct for non-deterministic
  sources (map iteration, channels); previously each call re-ran the source
  from the start and skipped a prefix.
- `Seq.RPosition`: single forward pass with O(1) memory instead of
  materializing the entire sequence into a slice.
