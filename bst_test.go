package bst

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"testing"
)

func TestBSTUnmarshalJSON(t *testing.T) {
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
		{
			name:    "returns error for invalid entry type",
			input:   `[{"key":123,"value":"alpha"}]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var tree BST[testEntry]
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

			if !slices.Equal([]testEntry(tree), tt.want) {
				t.Fatalf("Unmarshal() = %#v, want %#v", []testEntry(tree), tt.want)
			}
		})
	}
}

func ExampleBST() {
	tree := make(BST[testEntry], 0, 1024)

	_ = tree

	// Output:
}

func ExampleBST_Insert() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	fmt.Printf("exampleBST: %v\n", tree)

	// Output:
	// exampleBST: [{a alpha} {b bravo} {c charlie} {d delta}]
}

func ExampleBST_Get() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	val, ok := tree.Get("a")
	fmt.Printf("exampleBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleBST.Get("a"): {a alpha} / true
}

func ExampleBST_ForEach() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	if err := tree.ForEach(func(te testEntry) error {
		fmt.Printf("exampleBST.ForEach(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// exampleBST.ForEach(): {a alpha}
	// exampleBST.ForEach(): {b bravo}
	// exampleBST.ForEach(): {c charlie}
	// exampleBST.ForEach(): {d delta}
}

func ExampleBST_Cursor() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	_ = tree.Cursor()

	// Output:
}

func ExampleBST_Remove() {
	tree := make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{K: "a", V: "alpha"})
	tree.Insert(testEntry{K: "b", V: "bravo"})
	tree.Insert(testEntry{K: "c", V: "charlie"})
	tree.Insert(testEntry{K: "d", V: "delta"})

	tree.Remove("b")
	fmt.Printf("exampleBS.Remove(): %v\n", tree)

	// Output:
	// exampleBS.Remove(): [{a alpha} {c charlie} {d delta}]
}
