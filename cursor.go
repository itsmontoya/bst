package bst

// Cursor navigates entries in a BST.
type Cursor[T Entry] struct {
	b BST[T]

	index int
}

// Seek positions the cursor at key and returns the matching entry.
func (c *Cursor[T]) Seek(key string) (val T, ok bool) {
	var (
		out   T
		match bool
	)

	if out, c.index, match = c.b.get(key); match {
		return out, true
	}

	return val, false
}

// Prev moves the cursor to the previous entry.
func (c *Cursor[T]) Prev() (val T, ok bool) {
	if c.index-1 < 0 {
		return val, false
	}

	c.index--

	return c.b[c.index], true
}

// Next moves the cursor to the next entry.
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
