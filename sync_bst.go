package bst

import "sync"

func NewSync[T Entry](length int) (out *SyncBST[T]) {
	var s SyncBST[T]
	s.b = make(BST[T], 0, length)
	return &s
}

// SyncBST wraps a BST with an RWMutex for concurrent access.
type SyncBST[T Entry] struct {
	mux sync.RWMutex

	b BST[T]
}

// Get returns the entry with the given key under a read lock.
func (b *SyncBST[T]) Get(key string) (out T, ok bool) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.b.Get(key)
}

// ForEach calls fn for each entry in key order under a read lock.
func (b *SyncBST[T]) ForEach(fn func(T) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.b.ForEach(fn)
}

// Cursor calls fn with a cursor while holding a read lock.
func (b *SyncBST[T]) Cursor(fn func(*Cursor[T]) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	c := b.b.Cursor()
	err = fn(c)
	c.b = nil
	return
}

// Insert adds val or replaces the existing entry with the same key under a write lock.
func (b *SyncBST[T]) Insert(val T) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.b.Insert(val)
}

// Remove deletes the entry with the given key under a write lock.
func (b *SyncBST[T]) Remove(key string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.b.Remove(key)
}

// Length returns the number of entries under a read lock.
func (b *SyncBST[T]) Length() (n int) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return len(b.b)
}
