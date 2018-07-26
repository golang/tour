// +build OMIT

package main

import "fmt"

func main() {
	defer fmt.Println("světe")

	fmt.Println("ahoj")
}
