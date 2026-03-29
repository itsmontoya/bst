package bst

import (
	"fmt"
	"testing"
)

func TestCursorSeek(t *testing.T) {
	t.Parallel()

	tree := cursorBSTFromEntries(
		testEntry{key: "a", value: "alpha"},
		testEntry{key: "c", value: "charlie"},
		testEntry{key: "e", value: "echo"},
	)

	tests := []struct {
		name   string
		key    string
		want   testEntry
		wantOK bool
	}{
		{
			name:   "hit",
			key:    "c",
			want:   testEntry{key: "c", value: "charlie"},
			wantOK: true,
		},
		{
			name:   "miss between entries",
			key:    "d",
			wantOK: false,
		},
		{
			name:   "miss before first entry",
			key:    "0",
			wantOK: false,
		},
		{
			name:   "miss after last entry",
			key:    "z",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cursor := tree.Cursor()
			got, gotOK := cursor.Seek(tt.key)
			if gotOK != tt.wantOK {
				t.Fatalf("Seek(%q) ok = %v, want %v", tt.key, gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Seek(%q) = %#v, want %#v", tt.key, got, tt.want)
			}
		})
	}
}

func TestCursorPrev(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		key    string
		want   testEntry
		wantOK bool
		seekOK bool
	}{
		{
			name:   "returns previous entry",
			key:    "c",
			want:   testEntry{key: "a", value: "alpha"},
			wantOK: true,
			seekOK: true,
		},
		{
			name:   "at beginning returns false",
			key:    "a",
			wantOK: false,
			seekOK: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := cursorBSTFromEntries(
				testEntry{key: "a", value: "alpha"},
				testEntry{key: "c", value: "charlie"},
				testEntry{key: "e", value: "echo"},
			)

			cursor := tree.Cursor()
			_, seekOK := cursor.Seek(tt.key)
			if seekOK != tt.seekOK {
				t.Fatalf("Seek(%q) ok = %v, want %v", tt.key, seekOK, tt.seekOK)
			}

			got, gotOK := cursor.Prev()
			if gotOK != tt.wantOK {
				t.Fatalf("Prev() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Prev() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCursorNext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		key    string
		want   testEntry
		wantOK bool
		seekOK bool
	}{
		{
			name:   "returns next entry",
			key:    "c",
			want:   testEntry{key: "e", value: "echo"},
			wantOK: true,
			seekOK: true,
		},
		{
			name:   "at end returns false",
			key:    "e",
			wantOK: false,
			seekOK: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := cursorBSTFromEntries(
				testEntry{key: "a", value: "alpha"},
				testEntry{key: "c", value: "charlie"},
				testEntry{key: "e", value: "echo"},
			)

			cursor := tree.Cursor()
			_, seekOK := cursor.Seek(tt.key)
			if seekOK != tt.seekOK {
				t.Fatalf("Seek(%q) ok = %v, want %v", tt.key, seekOK, tt.seekOK)
			}

			got, gotOK := cursor.Next()
			if gotOK != tt.wantOK {
				t.Fatalf("Next() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Next() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCursorPrevAfterSeekMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		seekKey  string
		want     testEntry
		wantOK   bool
		wantSeek bool
	}{
		{
			name:     "miss before first returns false",
			seekKey:  "0",
			wantOK:   false,
			wantSeek: false,
		},
		{
			name:     "miss between entries returns lower neighbor",
			seekKey:  "d",
			want:     testEntry{key: "c", value: "charlie"},
			wantOK:   true,
			wantSeek: false,
		},
		{
			name:     "miss after last returns last entry",
			seekKey:  "z",
			want:     testEntry{key: "e", value: "echo"},
			wantOK:   true,
			wantSeek: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := cursorBSTFromEntries(
				testEntry{key: "a", value: "alpha"},
				testEntry{key: "c", value: "charlie"},
				testEntry{key: "e", value: "echo"},
			)

			cursor := tree.Cursor()
			_, seekOK := cursor.Seek(tt.seekKey)
			if seekOK != tt.wantSeek {
				t.Fatalf("Seek(%q) ok = %v, want %v", tt.seekKey, seekOK, tt.wantSeek)
			}

			got, gotOK := cursor.Prev()
			if gotOK != tt.wantOK {
				t.Fatalf("Prev() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Prev() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCursorNextAfterSeekMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		seekKey  string
		want     testEntry
		wantOK   bool
		wantSeek bool
	}{
		{
			name:     "miss before first returns second entry",
			seekKey:  "0",
			want:     testEntry{key: "c", value: "charlie"},
			wantOK:   true,
			wantSeek: false,
		},
		{
			name:     "miss between entries returns false",
			seekKey:  "d",
			wantOK:   false,
			wantSeek: false,
		},
		{
			name:     "miss after last returns false",
			seekKey:  "z",
			wantOK:   false,
			wantSeek: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := cursorBSTFromEntries(
				testEntry{key: "a", value: "alpha"},
				testEntry{key: "c", value: "charlie"},
				testEntry{key: "e", value: "echo"},
			)

			cursor := tree.Cursor()
			_, seekOK := cursor.Seek(tt.seekKey)
			if seekOK != tt.wantSeek {
				t.Fatalf("Seek(%q) ok = %v, want %v", tt.seekKey, seekOK, tt.wantSeek)
			}

			got, gotOK := cursor.Next()
			if gotOK != tt.wantOK {
				t.Fatalf("Next() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Next() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCursorFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tree    BST[testEntry]
		want    testEntry
		wantOK  bool
		wantKey string
	}{
		{
			name:   "empty tree",
			tree:   cursorBSTFromEntries(),
			wantOK: false,
		},
		{
			name:    "returns first entry and positions cursor at beginning",
			tree:    cursorBSTFromEntries(testEntry{key: "a", value: "alpha"}, testEntry{key: "c", value: "charlie"}, testEntry{key: "e", value: "echo"}),
			want:    testEntry{key: "a", value: "alpha"},
			wantOK:  true,
			wantKey: "c",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cursor := tt.tree.Cursor()
			got, gotOK := cursor.First()
			if gotOK != tt.wantOK {
				t.Fatalf("First() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("First() = %#v, want %#v", got, tt.want)
			}

			if !tt.wantOK {
				return
			}

			next, nextOK := cursor.Next()
			if !nextOK || next.key != tt.wantKey {
				t.Fatalf("Next() = (%#v, %v), want key %q and ok=true", next, nextOK, tt.wantKey)
			}
		})
	}
}

func TestCursorLast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tree    BST[testEntry]
		want    testEntry
		wantOK  bool
		wantKey string
	}{
		{
			name:   "empty tree",
			tree:   cursorBSTFromEntries(),
			wantOK: false,
		},
		{
			name:    "returns last entry and positions cursor at end",
			tree:    cursorBSTFromEntries(testEntry{key: "a", value: "alpha"}, testEntry{key: "c", value: "charlie"}, testEntry{key: "e", value: "echo"}),
			want:    testEntry{key: "e", value: "echo"},
			wantOK:  true,
			wantKey: "c",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cursor := tt.tree.Cursor()
			got, gotOK := cursor.Last()
			if gotOK != tt.wantOK {
				t.Fatalf("Last() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Last() = %#v, want %#v", got, tt.want)
			}

			if !tt.wantOK {
				return
			}

			prev, prevOK := cursor.Prev()
			if !prevOK || prev.key != tt.wantKey {
				t.Fatalf("Prev() = (%#v, %v), want key %q and ok=true", prev, prevOK, tt.wantKey)
			}
		})
	}
}

func ExampleCursor() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})

	_ = tree.Cursor()

	// Output:
}

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

func cursorBSTFromEntries(entries ...testEntry) (tree BST[testEntry]) {
	tree = make(BST[testEntry], 0, len(entries))
	for _, entry := range entries {
		tree.Insert(entry)
	}

	return tree
}
