// +build OMIT

package main

import "fmt"

func main() {
	var s []int
	printSlice(s)

	// append funguje na nil výřezech.
	s = append(s, 0)
	printSlice(s)

	// Výřez roste jak je potřeba.
	s = append(s, 1)
	printSlice(s)

	// Můžeme přidávat víc než jeden element najednou.
	s = append(s, 2, 3, 4)
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
