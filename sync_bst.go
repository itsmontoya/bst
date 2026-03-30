package bst

import (
	"encoding/json"
	"sync"
)

// NewSync creates a SyncBST preallocated with capacity length.
//
// length must be >= 0. A negative length panics (same behavior as make with
// a negative capacity).
func NewSync[T Entry](length int) (out *SyncBST[T]) {
	var s SyncBST[T]
	s.b = make(BST[T], 0, length)
	return &s
}

// SyncBST wraps a BST with an RWMutex for concurrent access.
//
// The zero value is ready to use:
//
//	var s SyncBST[T]
//
// A nil *SyncBST is not valid. Calling methods on a nil pointer panics.
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

// ForEach calls fn for each entry in key order while holding a read lock.
//
// fn executes before the read lock is released. Calling write methods on the
// same SyncBST instance from fn (for example Insert or Remove) can self-deadlock.
func (b *SyncBST[T]) ForEach(fn func(T) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.b.ForEach(fn)
}

// Cursor calls fn with a cursor while holding a read lock.
//
// fn executes before the read lock is released. Calling write methods on the
// same SyncBST instance from fn (for example Insert or Remove) can self-deadlock.
func (b *SyncBST[T]) Cursor(fn func(*Cursor[T]) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	c := b.b.Cursor()
	defer c.cleanup()
	return fn(c)
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

func (b *SyncBST[T]) MarshalJSON() (bs []byte, err error) {
	return json.Marshal(b.b)
}

func (b *SyncBST[T]) UnmarshalJSON(bs []byte) (err error) {
	return json.Unmarshal(bs, &b.b)
}
