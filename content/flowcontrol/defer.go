// +build OMIT

package main

import "fmt"

func main() {
	defer fmt.Println("world")

	fmt.Println("hello")
}
