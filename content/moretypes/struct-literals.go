// +build OMIT

package main

import "fmt"

type Vertex struct {
	X, Y int
}

var (
	v1 = Vertex{1, 2}  // posiada typ Vertex
	v2 = Vertex{X: 1}  // Y:0 jest domniemane
	v3 = Vertex{}      // X:0 oraz Y:0
	p  = &Vertex{1, 2} // posiada typ *Vertex
)

func main() {
	fmt.Println(v1, p, v2, v3)
}
