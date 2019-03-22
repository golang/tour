// +build OMIT

package main

import "fmt"

func main() {
	q := []rune{'a', 'b', 'c', 'd', 'e', 'f'}
	printChars(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	s := []struct {
		c rune
		b bool
	}{
		{'a', true},
		{'b', false},
		{'c', true},
		{'d', true},
		{'e', false},
		{'f', true},
	}
	fmt.Println(s)
}

func printChars(letters []rune ) {
	for i := range letters {
		fmt.Printf("%c", letters[i])
	}
	fmt.Println()
}
