// +build OMIT

package main

import "fmt"

func main() {
	letters := [7]rune{'a','b','c','d','e','f','g'}

	var slice []rune = letters[1:4]
	
	for i := range slice {
		fmt.Printf("%c", slice[i])
		}
}
