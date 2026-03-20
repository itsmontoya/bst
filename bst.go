package bst

import "sort"

type BST[T Entry] []T

func (b BST[T]) Get(key string) (out T, ok bool) {
	val, _, match := b.get(key)
	if !match {
		return out, false
	}

	return val, true
}

func (b BST[T]) ForEach(fn func(T) error) (err error) {
	for _, t := range b {
		if err = fn(t); err != nil {
			return err
		}
	}

	return nil
}

func (b BST[T]) Cursor() (out *Cursor[T]) {
	var c Cursor[T]
	c.b = b
	return &c
}

func (b *BST[T]) Insert(val T) {
	_, index, match := b.get(val.Key())
	if match {
		(*b)[index] = val
		return
	}

	b.insert(index, val)
}

func (b *BST[T]) Remove(key string) {
	_, index, match := b.get(key)
	if !match {
		return
	}

	b.remove(index)
}

func (b BST[T]) get(key string) (out T, index int, match bool) {
	index = sort.Search(len(b), func(i int) bool {
		return b[i].Key() >= key
	})

	if index >= len(b) {
		return out, index, false
	}

	out = b[index]
	return out, index, out.Key() == key
}

func (b *BST[T]) insert(index int, val T) {
	first := (*b)[:index]
	second := append([]T{val}, (*b)[index:]...)
	*b = append(first, second...)
}

func (b *BST[T]) remove(index int) {
	copy((*b)[index:], (*b)[index+1:])
	(*b)[len(*b)-1] = b.zero()
	*b = (*b)[:len(*b)-1]
}

func (b *BST[T]) zero() (out T) {
	return out
}
