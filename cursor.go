package bst

// Cursor navigates entries in a BST.
type Cursor[T Entry] struct {
	b BST[T]

	index int
}

// Seek positions the cursor using BST search semantics and returns a value at or after key.
//
// If key is present, Seek returns the matching entry and ok=true.
//
// If key is missing but there is a later entry, Seek returns that next larger entry and ok=true.
//
// If key is after the last entry, Seek returns (zero, false).
//
// Note: ok reports whether a value was found at or after key, not whether key matched exactly.
func (c *Cursor[T]) Seek(key string) (val T, ok bool) {
	var (
		out T
		i   int
	)

	if out, i, _ = c.b.get(key); i >= len(c.b) {
		return val, false
	}

	return out, true
}

// Prev moves the cursor to the previous entry relative to the current index.
//
// After Seek miss, Prev uses the insertion index semantics from Seek.
// Examples:
//   - miss before first: Prev() returns (zero, false)
//   - miss between entries: Prev() returns the lower neighbor
//   - miss after last: Prev() returns the last entry
func (c *Cursor[T]) Prev() (val T, ok bool) {
	if c.index-1 < 0 {
		return val, false
	}

	c.index--

	return c.b[c.index], true
}

// Next moves the cursor to the next entry relative to the current index.
//
// After Seek miss, Next uses the insertion index semantics from Seek.
// Examples:
//   - miss before first: Next() returns the second entry (if present)
//   - miss between entries: Next() returns (zero, false)
//   - miss after last: Next() returns (zero, false)
func (c *Cursor[T]) Next() (val T, ok bool) {
	if c.index+1 >= len(c.b) {
		return val, false
	}

	c.index++

	return c.b[c.index], true
}

// First moves the cursor to the first entry.
func (c *Cursor[T]) First() (val T, ok bool) {
	if len(c.b) == 0 {
		return val, false
	}

	c.index = 0

	return c.b[c.index], true
}

// Last moves the cursor to the last entry.
func (c *Cursor[T]) Last() (val T, ok bool) {
	if len(c.b) == 0 {
		return val, false
	}

	c.index = len(c.b) - 1

	return c.b[c.index], true
}

func (c *Cursor[T]) cleanup() {
	c.b = nil
}
