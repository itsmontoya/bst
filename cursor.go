package bst

type Cursor[T Entry] struct {
	b BST[T]

	index int
}

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

func (c *Cursor[T]) Prev() (val T, ok bool) {
	if c.index-1 < 0 {
		return val, false
	}

	c.index--

	return c.b[c.index], true
}

func (c *Cursor[T]) Next() (val T, ok bool) {
	if c.index+1 >= len(c.b) {
		return val, false
	}

	c.index++

	return c.b[c.index], true
}
