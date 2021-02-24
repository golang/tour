// +build OMIT

package main

import "fmt"

func main() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Wytnij wycinek by nadać mu długość zerową.
	s = s[:0]
	printSlice(s)

	// Powiększ długość wycinka.
	s = s[:4]
	printSlice(s)

	// Usuń z wycinka pierwsze dwie wartości.
	s = s[2:]
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
