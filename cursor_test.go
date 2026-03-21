package bst

import "fmt"

var exampleCursor *Cursor[testEntry]

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
