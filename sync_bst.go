package bst

import "sync"

type SyncBST[T Entry] struct {
	mux sync.RWMutex

	b BST[T]
}

func (b *SyncBST[T]) Get(key string) (out T, ok bool) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.b.Get(key)
}

func (b *SyncBST[T]) ForEach(fn func(T) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.b.ForEach(fn)
}

func (b *SyncBST[T]) Cursor(fn func(*Cursor[T]) error) (err error) {
	b.mux.RLock()
	defer b.mux.RUnlock()
	c := b.b.Cursor()
	err = fn(c)
	c.b = nil
	return
}

func (b *SyncBST[T]) Insert(val T) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.b.Insert(val)
}

func (b *SyncBST[T]) Remove(key string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.b.Remove(key)
}
