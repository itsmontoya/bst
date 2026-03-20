package bst

// Entry is a value that can be stored in a BST.
type Entry interface {
	// Key returns the entry's sort key.
	Key() string
}
