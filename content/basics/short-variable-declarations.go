// +build OMIT

package main

import "fmt"

func main() {
	var i, j int = 1, 2
	k := 3
	c, python, java := true, false, "ne!"

	fmt.Println(i, j, k, c, python, java)
}
