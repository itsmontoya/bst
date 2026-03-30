# bst &emsp; [![GoDoc][GoDoc Badge]][GoDoc URL] ![Coverage] [![Go Report Card][Report Card Badge]][Report Card URL] [![MIT licensed][License Badge]][License URL]

[GoDoc Badge]: https://pkg.go.dev/badge/github.com/itsmontoya/bst
[GoDoc URL]: https://pkg.go.dev/github.com/itsmontoya/bst
[Coverage]: https://img.shields.io/badge/coverage-100%25-brightgreen
[License Badge]: https://img.shields.io/badge/license-MIT-blue.svg
[License URL]: https://github.com/itsmontoya/bst/blob/main/LICENSE
[Report Card Badge]: https://goreportcard.com/badge/github.com/itsmontoya/bst
[Report Card URL]: https://goreportcard.com/report/github.com/itsmontoya/bst

![banner](https://res.cloudinary.com/dryepxxoy/image/upload/v1774414470/bst_1920_lpr9y2.webp "BST banner")

`bst` is a Go library that provides a **sorted string-keyed collection** with binary-search lookups.

It stores entries in key order, supports in-place insert/replace semantics, and includes cursor helpers for neighbor traversal.

For concurrent access, `SyncBST` wraps the tree with an `RWMutex` and provides the same operations under read/write locks.

## Motivation

Go gives you maps for fast keyed lookup, but maps do not preserve ordering.

If your workload needs:

- Stable key ordering
- Predictable iteration order
- Keyed lookup
- Neighbor navigation (`Prev`/`Next`)

then a sorted structure like `bst` can be a better fit.

## Examples

### Initialize BST
```go
func ExampleBST() {
	tree := make(BST[testEntry], 0, 1024)

	_ = tree
}
```

### BST.Insert
```go
func ExampleBST_Insert() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	fmt.Printf("exampleBST: %v\n", tree)

	// Output:
	// exampleBST: [{a alpha} {b bravo} {c charlie} {d delta}]
}
```

### BST.Get
```go
func ExampleBST_Get() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	val, ok := tree.Get("a")
	fmt.Printf("exampleBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleBST.Get("a"): {a alpha} / true
}
```

### BST.ForEach
```go
func ExampleBST_ForEach() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	if err := tree.ForEach(func(te testEntry) error {
		fmt.Printf("exampleBST.ForEach(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// exampleBST.ForEach(): {a alpha}
	// exampleBST.ForEach(): {b bravo}
	// exampleBST.ForEach(): {c charlie}
	// exampleBST.ForEach(): {d delta}
}
```

### BST.Cursor
```go
func ExampleBST_Cursor() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	_ = tree.Cursor()
}
```

### BST.Remove
```go
func ExampleBST_Remove() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	tree.Remove("b")
	fmt.Printf("exampleBS.Remove(): %v\n", tree)

	// Output:
	// exampleBS.Remove(): [{a alpha} {c charlie} {d delta}]
}
```

### BST.UnmarshalJSON

Inbound JSON entries are sorted by `Key()` during unmarshal before assignment.

```go
func ExampleBST_UnmarshalJSON() {
	var tree BST[testEntry]

	if err := json.Unmarshal([]byte(`[{"key":"c","value":"charlie"},{"key":"a","value":"alpha"},{"key":"b","value":"bravo"}]`), &tree); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("exampleBST.UnmarshalJSON(): %v\n", tree)

	// Output:
	// exampleBST.UnmarshalJSON(): [{a alpha} {b bravo} {c charlie}]
}
```

### Cursor.Seek
```go
func ExampleCursor_Seek() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	cursor := tree.Cursor()

	val, ok := cursor.Seek("d")
	fmt.Printf("cursor.Seek(%q): %v / %v\n", "d", val, ok)

	// Output:
	// cursor.Seek("d"): {d delta} / true
}
```

### Cursor.Prev
```go
func ExampleCursor_Prev() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	cursor := tree.Cursor()
	_, _ = cursor.Seek("d")

	val, ok := cursor.Prev()
	fmt.Printf("cursor.Prev(): %v / %v\n", val, ok)

	// Output:
	// cursor.Prev(): {c charlie} / true
}
```

### Cursor.Next
```go
func ExampleCursor_Next() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	cursor := tree.Cursor()
	_, _ = cursor.Seek("c")

	val, ok := cursor.Next()
	fmt.Printf("cursor.Next(): %v / %v\n", val, ok)

	// Output:
	// cursor.Next(): {d delta} / true
}
```

### Cursor.First
```go
func ExampleCursor_First() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	cursor := tree.Cursor()

	val, ok := cursor.First()
	fmt.Printf("cursor.First(): %v / %v\n", val, ok)

	// Output:
	// cursor.First(): {a alpha} / true
}
```

### Cursor.Last
```go
func ExampleCursor_Last() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	cursor := tree.Cursor()

	val, ok := cursor.Last()
	fmt.Printf("cursor.Last(): %v / %v\n", val, ok)

	// Output:
	// cursor.Last(): {d delta} / true
}
```

### Initialize SyncBST
```go
func ExampleSyncBST() {
	tree := NewSync[testEntry](1024)

	_ = tree
}
```

`SyncBST` has a usable zero value, so this is also valid:

```go
func ExampleSyncBST_ZeroValue() {
	var tree SyncBST[testEntry]

	tree.Insert(testEntry{key: "a", value: "alpha"})
}
```

Use a non-nil value for method calls. `var tree *SyncBST[testEntry]` is nil, and
calling methods on it panics.

`NewSync(length int)` follows Go's `make` behavior: negative `length` panics.

### SyncBST.Insert
```go
func ExampleSyncBST_Insert() {
	tree := NewSync[testEntry](1024)

	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
}
```

### SyncBST.Get
```go
func ExampleSyncBST_Get() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	val, ok := tree.Get("a")
	fmt.Printf("exampleSyncBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleSyncBST.Get("a"): {a alpha} / true
}
```

### SyncBST.ForEach
```go
func ExampleSyncBST_ForEach() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	if err := tree.ForEach(func(te testEntry) error {
		fmt.Printf("exampleSyncBST.ForEach(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// exampleSyncBST.ForEach(): {a alpha}
	// exampleSyncBST.ForEach(): {b bravo}
	// exampleSyncBST.ForEach(): {c charlie}
	// exampleSyncBST.ForEach(): {d delta}
}
```

`SyncBST.ForEach` runs the callback while `RLock` is held. Do not call `Insert` or `Remove`
on that same `SyncBST` from inside the callback, or you can self-deadlock.

### SyncBST.Cursor
```go
func ExampleSyncBST_Cursor() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	if err := tree.Cursor(func(cursor *Cursor[testEntry]) error {
		val, ok := cursor.Seek("b")
		fmt.Printf("cursor.Seek(%q): %v / %v\n", "b", val, ok)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// cursor.Seek("b"): {b bravo} / true
}
```

`SyncBST.Cursor` runs the callback while `RLock` is held. Do not call `Insert` or `Remove`
on that same `SyncBST` from inside the callback, or you can self-deadlock.

### SyncBST.Remove
```go
func ExampleSyncBST_Remove() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	tree.Remove("b")
	fmt.Printf("exampleSyncBST.Length(): %d\n", tree.Length())

	// Output:
	// exampleSyncBST.Length(): 3
}
```

### SyncBST.MarshalJSON
```go
func ExampleSyncBST_MarshalJSON() {
	tree := NewSync[testEntry](0)
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "a", V: "alpha"})

	bs, err := json.Marshal(tree)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("exampleSyncBST.MarshalJSON(): %s\n", bs)

	// Output:
	// exampleSyncBST.MarshalJSON(): [{"key":"a","value":"alpha"},{"key":"b","value":"bravo"}]
}
```

### SyncBST.UnmarshalJSON

Inbound JSON entries are sorted by `Key()` during unmarshal before assignment.

```go
func ExampleSyncBST_UnmarshalJSON() {
	var tree SyncBST[testEntry]

	if err := json.Unmarshal([]byte(`[{"key":"c","value":"charlie"},{"key":"a","value":"alpha"}]`), &tree); err != nil {
		log.Fatal(err)
	}

	if err := tree.ForEach(func(te testEntry) error {
		fmt.Printf("exampleSyncBST.UnmarshalJSON(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// exampleSyncBST.UnmarshalJSON(): {a alpha}
	// exampleSyncBST.UnmarshalJSON(): {c charlie}
}
```

## Core Concepts

### Sorted storage

Entries are always maintained in ascending `Key()` order.

### Lookup and update behavior

`Get` uses binary search.

`Insert` replaces an existing entry when keys match, otherwise inserts at the sorted position.

### Deterministic iteration

`ForEach` always walks entries in key order.

### Cursor navigation

A cursor can seek to a key, then move to adjacent entries with `Prev` and `Next`.

### Concurrency model

`BST` is not synchronized.

`SyncBST` provides thread-safe access by wrapping operations with an `RWMutex`.

`SyncBST.ForEach` and `SyncBST.Cursor` hold `RLock` for the entire callback execution.
Callbacks should be read-only with respect to the same `SyncBST` instance.
Calling `Insert` or `Remove` from those callbacks can self-deadlock.

## AI Usage and Authorship

This project is intentionally **human-authored** for all logic.

To be explicit:

- AI does **not** write or modify non-test code in this repository.
- AI does **not** make architectural or behavioral decisions.
- AI may assist with documentation, comments, and test scaffolding only.
- All implementation logic is written and reviewed by human maintainers.

These boundaries are enforced in `AGENTS.md` and are part of this repository's contribution discipline.

## Contributors

- Human maintainers: library design, implementation, and behavior decisions.
- ChatGPT Codex: documentation, test coverage support, and comments.
- Google Gemini: README artwork generation.

![footer](https://res.cloudinary.com/dryepxxoy/image/upload/v1774414470/bst_footer_1920_ia720l.webp "BST footer")
