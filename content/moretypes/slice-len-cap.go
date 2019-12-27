// +build OMIT

package main

import "fmt"

func main() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Rozřízne řez aby dostal nulovou délku.
	s = s[:0]
	printSlice(s)

	// Rozšíří jeho délku.
	s = s[:4]
	printSlice(s)

	// Vypustí jeho první dvě hodnoty.
	s = s[2:]
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
