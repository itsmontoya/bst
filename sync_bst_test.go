package bst

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"testing"
)

var exampleSyncBST *SyncBST[testEntry]

func TestSyncBSTGet(t *testing.T) {
	t.Parallel()

	tree := syncBSTFromEntries(
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
			name:   "miss after last entry",
			key:    "z",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotOK := tree.Get(tt.key)
			if gotOK != tt.wantOK {
				t.Fatalf("Get(%q) ok = %v, want %v", tt.key, gotOK, tt.wantOK)
			}

			if got != tt.want {
				t.Fatalf("Get(%q) = %#v, want %#v", tt.key, got, tt.want)
			}
		})
	}
}

func TestSyncBSTForEach(t *testing.T) {
	t.Parallel()

	tree := syncBSTFromEntries(
		testEntry{key: "a", value: "alpha"},
		testEntry{key: "b", value: "bravo"},
		testEntry{key: "c", value: "charlie"},
	)

	sentinel := errors.New("stop")

	tests := []struct {
		name     string
		fn       func(testEntry) error
		wantKeys []string
		wantErr  error
	}{
		{
			name: "iterates in order",
			fn: func(testEntry) error {
				return nil
			},
			wantKeys: []string{"a", "b", "c"},
		},
		{
			name: "stops on error",
			fn: func(entry testEntry) error {
				if entry.key == "b" {
					return sentinel
				}

				return nil
			},
			wantKeys: []string{"a", "b"},
			wantErr:  sentinel,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotKeys []string
			err := tree.ForEach(func(entry testEntry) error {
				gotKeys = append(gotKeys, entry.key)
				return tt.fn(entry)
			})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ForEach() error = %v, want %v", err, tt.wantErr)
			}

			if !slices.Equal(gotKeys, tt.wantKeys) {
				t.Fatalf("ForEach() visited %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func TestSyncBSTCursor(t *testing.T) {
	t.Parallel()

	tree := syncBSTFromEntries(
		testEntry{key: "a", value: "alpha"},
		testEntry{key: "c", value: "charlie"},
		testEntry{key: "e", value: "echo"},
	)

	tests := []struct {
		name string
		run  func(*testing.T, *SyncBST[testEntry])
	}{
		{
			name: "seek existing and navigate neighbors",
			run: func(t *testing.T, tree *SyncBST[testEntry]) {
				err := tree.Cursor(func(c *Cursor[testEntry]) error {
					got, ok := c.Seek("c")
					if !ok || got.key != "c" {
						t.Fatalf("Seek(%q) = (%#v, %v), want key %q and ok=true", "c", got, ok, "c")
					}

					prev, ok := c.Prev()
					if !ok || prev.key != "a" {
						t.Fatalf("Prev() = (%#v, %v), want key %q and ok=true", prev, ok, "a")
					}

					next, ok := c.Next()
					if !ok || next.key != "c" {
						t.Fatalf("Next() = (%#v, %v), want key %q and ok=true", next, ok, "c")
					}

					return nil
				})
				if err != nil {
					t.Fatalf("Cursor() error = %v, want nil", err)
				}
			},
		},
		{
			name: "seek missing returns false",
			run: func(t *testing.T, tree *SyncBST[testEntry]) {
				err := tree.Cursor(func(c *Cursor[testEntry]) error {
					got, ok := c.Seek("d")
					if ok {
						t.Fatalf("Seek(%q) = (%#v, %v), want ok=false", "d", got, ok)
					}

					return nil
				})
				if err != nil {
					t.Fatalf("Cursor() error = %v, want nil", err)
				}
			},
		},
		{
			name: "prev at beginning returns false",
			run: func(t *testing.T, tree *SyncBST[testEntry]) {
				err := tree.Cursor(func(c *Cursor[testEntry]) error {
					_, _ = c.Seek("a")

					got, ok := c.Prev()
					if ok {
						t.Fatalf("Prev() = (%#v, %v), want ok=false", got, ok)
					}

					return nil
				})
				if err != nil {
					t.Fatalf("Cursor() error = %v, want nil", err)
				}
			},
		},
		{
			name: "next at end returns false",
			run: func(t *testing.T, tree *SyncBST[testEntry]) {
				err := tree.Cursor(func(c *Cursor[testEntry]) error {
					_, _ = c.Seek("e")

					got, ok := c.Next()
					if ok {
						t.Fatalf("Next() = (%#v, %v), want ok=false", got, ok)
					}

					return nil
				})
				if err != nil {
					t.Fatalf("Cursor() error = %v, want nil", err)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.run(t, tree)
		})
	}
}

func TestSyncBSTInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tree []testEntry
		val  testEntry
		want []testEntry
	}{
		{
			name: "insert into empty tree",
			val:  testEntry{key: "b", value: "bravo"},
			want: []testEntry{
				{key: "b", value: "bravo"},
			},
		},
		{
			name: "insert at beginning",
			tree: []testEntry{
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
			val: testEntry{key: "a", value: "alpha"},
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
		},
		{
			name: "insert in middle",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "e", value: "echo"},
			},
			val: testEntry{key: "c", value: "charlie"},
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
		},
		{
			name: "replace existing key",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
			},
			val: testEntry{key: "c", value: "updated"},
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "updated"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := syncBSTFromEntries(tt.tree...)
			tree.Insert(tt.val)

			got, err := collectSyncBST(tree)
			if err != nil {
				t.Fatalf("collectSyncBST() error = %v, want nil", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Fatalf("Insert(%#v) = %#v, want %#v", tt.val, got, tt.want)
			}
		})
	}
}

func TestSyncBSTRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tree []testEntry
		key  string
		want []testEntry
	}{
		{
			name: "remove missing key",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
			},
			key: "b",
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
			},
		},
		{
			name: "remove first entry",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
			key: "a",
			want: []testEntry{
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
		},
		{
			name: "remove middle entry",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
			key: "c",
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "e", value: "echo"},
			},
		},
		{
			name: "remove last entry",
			tree: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
				{key: "e", value: "echo"},
			},
			key: "e",
			want: []testEntry{
				{key: "a", value: "alpha"},
				{key: "c", value: "charlie"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tree := syncBSTFromEntries(tt.tree...)
			tree.Remove(tt.key)

			got, err := collectSyncBST(tree)
			if err != nil {
				t.Fatalf("collectSyncBST() error = %v, want nil", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Fatalf("Remove(%q) = %#v, want %#v", tt.key, got, tt.want)
			}
		})
	}
}

func ExampleSyncBST() {
	exampleSyncBST = NewSync[testEntry](1024)

	// Output:
}

func ExampleSyncBST_Insert() {
	exampleSyncBST.Insert(testEntry{key: "a", value: "alpha"})
	exampleSyncBST.Insert(testEntry{key: "b", value: "bravo"})
	exampleSyncBST.Insert(testEntry{key: "c", value: "charlie"})
	exampleSyncBST.Insert(testEntry{key: "d", value: "delta"})

	// Output:
}

func ExampleSyncBST_Get() {
	val, ok := exampleSyncBST.Get("a")
	fmt.Printf("exampleSyncBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleSyncBST.Get("a"): {a alpha} / true
}

func ExampleSyncBST_ForEach() {
	if err := exampleSyncBST.ForEach(func(te testEntry) error {
		fmt.Printf("exampleSyncBST.ForEach(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// exampleSyncBST.ForEach(): {a alpha}
	// exampleSyncBST.ForEach(): {b bravo}
	// exampleSyncBST.ForEach(): {c charlie}
	// exampleSyncBST.ForEach(): {d delta}
}

func ExampleSyncBST_Cursor() {
	if err := exampleSyncBST.Cursor(func(cursor *Cursor[testEntry]) error {
		val, ok := cursor.Seek("b")
		fmt.Printf("cursor.Seek(%q): %v / %v\n", "b", val, ok)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// cursor.Seek("b"): {b bravo} / true
}

func ExampleSyncBST_Remove() {
	exampleSyncBST.Remove("b")
	fmt.Printf("exampleSyncBST.Length(): %d\n", exampleSyncBST.Length())

	// Output:
	// exampleSyncBST.Length(): 3
}

func syncBSTFromEntries(entries ...testEntry) *SyncBST[testEntry] {
	tree := &SyncBST[testEntry]{}
	for _, entry := range entries {
		tree.Insert(entry)
	}

	return tree
}

func collectSyncBST(tree *SyncBST[testEntry]) ([]testEntry, error) {
	var out []testEntry
	err := tree.ForEach(func(entry testEntry) error {
		out = append(out, entry)
		return nil
	})

	return out, err
}

type testEntry struct {
	key   string
	value string
}

func (e testEntry) Key() string {
	return e.key
}
