package bst

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"testing"
)

func TestSyncBSTGet(t *testing.T) {
	t.Parallel()

	tree := syncBSTFromEntries(
		testEntry{K: "a", V: "alpha"},
		testEntry{K: "c", V: "charlie"},
		testEntry{K: "e", V: "echo"},
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
			want:   testEntry{K: "c", V: "charlie"},
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
		testEntry{K: "a", V: "alpha"},
		testEntry{K: "b", V: "bravo"},
		testEntry{K: "c", V: "charlie"},
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
				if entry.K == "b" {
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
				gotKeys = append(gotKeys, entry.K)
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
		testEntry{K: "a", V: "alpha"},
		testEntry{K: "c", V: "charlie"},
		testEntry{K: "e", V: "echo"},
	)

	err := tree.Cursor(func(cursor *Cursor[testEntry]) error {
		if cursor == nil {
			t.Fatal("Cursor() passed nil cursor")
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Cursor() error = %v, want nil", err)
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
			val:  testEntry{K: "b", V: "bravo"},
			want: []testEntry{
				{K: "b", V: "bravo"},
			},
		},
		{
			name: "insert at beginning",
			tree: []testEntry{
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
			val: testEntry{K: "a", V: "alpha"},
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
		},
		{
			name: "insert in middle",
			tree: []testEntry{
				{K: "a", V: "alpha"},
				{K: "e", V: "echo"},
			},
			val: testEntry{K: "c", V: "charlie"},
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
		},
		{
			name: "replace existing key",
			tree: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
			},
			val: testEntry{K: "c", V: "updated"},
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "updated"},
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
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
			},
			key: "b",
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
			},
		},
		{
			name: "remove first entry",
			tree: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
			key: "a",
			want: []testEntry{
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
		},
		{
			name: "remove middle entry",
			tree: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
			key: "c",
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "e", V: "echo"},
			},
		},
		{
			name: "remove last entry",
			tree: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
				{K: "e", V: "echo"},
			},
			key: "e",
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "c", V: "charlie"},
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

func TestSyncBSTMarshalJSON(t *testing.T) {
	t.Parallel()

	tree := &SyncBST[testEntry]{}
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "c", V: "charlie"})

	bs, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("Marshal() error = %v, want nil", err)
	}

	var got []testEntry
	if err = json.Unmarshal(bs, &got); err != nil {
		t.Fatalf("Unmarshal(marshaled tree) error = %v, want nil", err)
	}

	want := []testEntry{
		{K: "a", V: "alpha"},
		{K: "b", V: "bravo"},
		{K: "c", V: "charlie"},
	}

	if !slices.Equal(got, want) {
		t.Fatalf("Marshal() = %#v, want %#v", got, want)
	}
}

func TestSyncBSTUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []testEntry
		wantErr bool
	}{
		{
			name:  "sorts inbound entries by key",
			input: `[{"key":"c","value":"charlie"},{"key":"a","value":"alpha"},{"key":"b","value":"bravo"}]`,
			want: []testEntry{
				{K: "a", V: "alpha"},
				{K: "b", V: "bravo"},
				{K: "c", V: "charlie"},
			},
		},
		{
			name:    "returns error for invalid json",
			input:   `{"key":"a"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var tree SyncBST[testEntry]
			err := json.Unmarshal([]byte(tt.input), &tree)
			if tt.wantErr {
				if err == nil {
					t.Fatal("Unmarshal() error = nil, want non-nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unmarshal() error = %v, want nil", err)
			}

			var got []testEntry
			err = tree.ForEach(func(entry testEntry) error {
				got = append(got, entry)
				return nil
			})
			if err != nil {
				t.Fatalf("ForEach() error = %v, want nil", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Fatalf("Unmarshal() entries = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func ExampleSyncBST() {
	tree := NewSync[testEntry](1024)

	_ = tree

	// Output:
}

func ExampleSyncBST_Insert() {
	tree := NewSync[testEntry](1024)

	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	// Output:
}

func ExampleSyncBST_Get() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	val, ok := tree.Get("a")
	fmt.Printf("exampleSyncBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleSyncBST.Get("a"): {a alpha} / true
}

func ExampleSyncBST_ForEach() {
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	if err := tree.ForEach(func(te testEntry) error {
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
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	if err := tree.Cursor(func(cursor *Cursor[testEntry]) error {
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
	tree := NewSync[testEntry](1024)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	tree.Remove("b")
	fmt.Printf("exampleSyncBST.Length(): %d\n", tree.Length())

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
	K string `json:"key"`
	V string `json:"value"`
}

func (e testEntry) Key() string {
	return e.K
}
