# bst &emsp; [![GoDoc][GoDoc Badge]][GoDoc URL] [![Go Report Card][Report Card Badge]][Report Card URL]

[GoDoc Badge]: https://pkg.go.dev/badge/github.com/itsmontoya/bst
[GoDoc URL]: https://pkg.go.dev/github.com/itsmontoya/bst
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
	users := make(bst.BST[userEntry], 0, 8)
	_ = users
}
```

### BST.Insert
```go
func ExampleBST_Insert() {
	users := make(bst.BST[userEntry], 0, 8)

	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "carol"})
	users.Insert(userEntry{id: 3, name: "bob"})
	users.Insert(userEntry{id: 4, name: "carol"}) // replaces existing "carol"
}
```

### BST.Get
```go
func ExampleBST_Get() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})

	entry, ok := users.Get("alice")
	fmt.Println(entry.id, ok)
}
```

### BST.ForEach
```go
func ExampleBST_ForEach() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	_ = users.ForEach(func(entry userEntry) error {
		fmt.Println(entry.name)
		return nil
	})
}
```

### BST.Cursor
```go
func ExampleBST_Cursor() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})
	users.Insert(userEntry{id: 3, name: "carol"})

	cursor := users.Cursor()
	_, _ = cursor.Seek("bob")
	prev, _ := cursor.Prev()
	next, _ := cursor.Next()

	fmt.Println(prev.name, next.name)
}
```

### BST.Remove
```go
func ExampleBST_Remove() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	users.Remove("alice")
}
```

### Cursor.Seek
```go
func ExampleCursor_Seek() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	cursor := users.Cursor()
	entry, ok := cursor.Seek("bob")
	fmt.Println(entry.id, ok)
}
```

### Cursor.Prev
```go
func ExampleCursor_Prev() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	cursor := users.Cursor()
	_, _ = cursor.Seek("bob")
	entry, ok := cursor.Prev()
	fmt.Println(entry.name, ok)
}
```

### Cursor.Next
```go
func ExampleCursor_Next() {
	users := make(bst.BST[userEntry], 0, 8)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	cursor := users.Cursor()
	_, _ = cursor.Seek("alice")
	entry, ok := cursor.Next()
	fmt.Println(entry.name, ok)
}
```

### NewSync
```go
func ExampleNewSync() {
	users := bst.NewSync[userEntry](1024)
	_ = users
}
```

### SyncBST.Insert
```go
func ExampleSyncBST_Insert() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})
}
```

### SyncBST.Get
```go
func ExampleSyncBST_Get() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})

	entry, ok := users.Get("alice")
	fmt.Println(entry.id, ok)
}
```

### SyncBST.ForEach
```go
func ExampleSyncBST_ForEach() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	_ = users.ForEach(func(entry userEntry) error {
		fmt.Println(entry.name)
		return nil
	})
}
```

### SyncBST.Cursor
```go
func ExampleSyncBST_Cursor() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})
	users.Insert(userEntry{id: 3, name: "carol"})

	_ = users.Cursor(func(cursor *bst.Cursor[userEntry]) error {
		_, _ = cursor.Seek("bob")
		prev, _ := cursor.Prev()
		next, _ := cursor.Next()
		fmt.Println(prev.name, next.name)
		return nil
	})
}
```

### SyncBST.Remove
```go
func ExampleSyncBST_Remove() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Remove("alice")
}
```

### SyncBST.Length
```go
func ExampleSyncBST_Length() {
	users := bst.NewSync[userEntry](1024)
	users.Insert(userEntry{id: 1, name: "alice"})
	users.Insert(userEntry{id: 2, name: "bob"})

	fmt.Println(users.Length())
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