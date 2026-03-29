package bst

import (
	"fmt"
	"log"
)

func ExampleBST() {
	tree := make(BST[testEntry], 0, 1024)

	_ = tree

	// Output:
}

func ExampleBST_Insert() {
	tree := exampleBSTWithEntries()

	fmt.Printf("exampleBST: %v\n", tree)

	// Output:
	// exampleBST: [{a alpha} {b bravo} {c charlie} {d delta}]
}

func ExampleBST_Get() {
	tree := exampleBSTWithEntries()

	val, ok := tree.Get("a")
	fmt.Printf("exampleBST.Get(%q): %v / %v\n", "a", val, ok)

	// Output:
	// exampleBST.Get("a"): {a alpha} / true
}

func ExampleBST_ForEach() {
	tree := exampleBSTWithEntries()

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
	tree := exampleBSTWithEntries()

	_ = tree.Cursor()

	// Output:
}

func ExampleBST_Remove() {
	tree := exampleBSTWithEntries()

	tree.Remove("b")
	fmt.Printf("exampleBS.Remove(): %v\n", tree)

	// Output:
	// exampleBS.Remove(): [{a alpha} {c charlie} {d delta}]
}

func exampleBSTWithEntries() (tree BST[testEntry]) {
	tree = make(BST[testEntry], 0, 4)
	tree.Insert(testEntry{key: "a", value: "alpha"})
	tree.Insert(testEntry{key: "b", value: "bravo"})
	tree.Insert(testEntry{key: "c", value: "charlie"})
	tree.Insert(testEntry{key: "d", value: "delta"})
	return tree
}
