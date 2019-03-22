// +build OMIT

package main

import "fmt"

func main() {
	s := []rune{'a', 'b', 'c', 'd', 'e', 'f'}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)
}

func printSlice(s []rune) {
	fmt.Printf("len=%d cap=%d %c\n", len(s), cap(s), s)
}
