// +build OMIT

package main

import "fmt"

func main() {
	s := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k'}

	s = s[1:4]
	prynt(s)

	s = s[:2]
	prynt(s)

	s = s[1:]
	prynt(s)
}

func prynt(charslice []rune) {
	for i := range charslice {
		fmt.Printf("%c", charslice[i])
		}
	fmt.Print("\n")
	}
