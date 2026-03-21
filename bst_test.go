package bst

import (
	"fmt"
	"log"
)

var exampleBST BST[testEntry]

func ExampleBST() {
	exampleBST = make(BST[testEntry], 0, 1024)

	fmt.Printf("exampleBST: %v\n", exampleBST)
	// Output:
	// exampleBST: []
}

func ExampleBST_Insert() {
	exampleBST.Insert(testEntry{key: "a", value: "alpha"})
	exampleBST.Insert(testEntry{key: "b", value: "bravo"})
	exampleBST.Insert(testEntry{key: "c", value: "charlie"})
	exampleBST.Insert(testEntry{key: "d", value: "delta"})

	fmt.Printf("exampleBST: %v\n", exampleBST)
	// Output:
	// exampleBST: [{a alpha}, {b bravo}, {c charlie}, {d delta}]
}

func ExampleBST_Get() {
	val, ok := exampleBST.Get("a")

	fmt.Printf("exampleBST.Get(%q): %v / %v\n", "a", val, ok)
	// Output:
	// tree.Get("a"): {a alpha} / true
}

func ExampleBST_ForEach() {
	if err := exampleBST.ForEach(func(te testEntry) error {
		fmt.Printf("exampleBST.ForEach(): %v\n", te)
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Output:
	// tree.ForEach(): {a alpha}
	// tree.ForEach(): {b bravo}
	// tree.ForEach(): {c charlie}
	// tree.ForEach(): {d delta}
}

func ExampleBST_Cursor() {
	cursor := exampleBST.Cursor()

	val, ok := cursor.Seek("b")
	fmt.Printf("cursor.Seek(%q): %v / %v\n", "b", val, ok)
	// Output:
	// cursor.Seek("b"): {b bravo} / true
}

func ExampleBST_Remove() {
	exampleBST.Remove("b")

	fmt.Printf("exampleBST: %v\n", exampleBST)
	// Output:
	// tree: [{a alpha} {c charlie}]
}
