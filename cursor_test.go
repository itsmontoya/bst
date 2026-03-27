package bst

import (
	"fmt"
	"testing"
)

var exampleCursor *Cursor[testEntry]

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
	exampleCursor = exampleBST.Cursor()

	// Output:
}

func ExampleCursor_Seek() {
	val, ok := exampleCursor.Seek("d")
	fmt.Printf("cursor.Seek(%q): %v / %v\n", "d", val, ok)

	// Output:
	// cursor.Seek("d"): {d delta} / true
}

func ExampleCursor_Prev() {
	val, ok := exampleCursor.Prev()
	fmt.Printf("cursor.Prev(): %v / %v\n", val, ok)

	// Output:
	// cursor.Prev(): {c charlie} / true
}

func ExampleCursor_Next() {
	val, ok := exampleCursor.Next()
	fmt.Printf("cursor.Next(): %v / %v\n", val, ok)

	// Output:
	// cursor.Next(): {d delta} / true
}

func cursorBSTFromEntries(entries ...testEntry) (tree BST[testEntry]) {
	tree = make(BST[testEntry], 0, len(entries))
	for _, entry := range entries {
		tree.Insert(entry)
	}

	return tree
}
