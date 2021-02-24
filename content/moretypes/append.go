// +build OMIT

package main

import "fmt"

func main() {
	var s []int
	printSlice(s)

	// append działa na wycinkach nilowych.
	s = append(s, 0)
	printSlice(s)

	// Wycinek powiększa się w miarę potrzeb.
	s = append(s, 1)
	printSlice(s)

	// Możemy dodać więcej niż jeden element w danym czasie.
	s = append(s, 2, 3, 4)
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
