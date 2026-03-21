package bst

import "fmt"

var exampleCursor *Cursor[testEntry]

func ExampleCursor() {
	exampleCursor = exampleBST.Cursor()
}

func ExampleCursor_Seek() {
	val, ok := exampleCursor.Seek("c")
	fmt.Printf("cursor.Seek(%q): %v / %v\n", "c", val, ok)

	// Output:
	// cursor.Seek("c"): {c charlie} / true
}

func ExampleCursor_Prev() {
	val, ok := exampleCursor.Prev()
	fmt.Printf("cursor.Prev(): %v / %v\n", val, ok)

	// Output:
	// cursor.Prev(): {b bravo} / true
}

func ExampleCursor_Next() {
	val, ok := exampleCursor.Next()
	fmt.Printf("cursor.Next(): %v / %v\n", val, ok)

	// Output:
	// cursor.Next(): {c charlie} / true
}
