// +build no-build OMIT

package main

import "golang.org/x/tour/tree"

// Walk prochází strojm t a posílá všechny hodnoty
// ze stromu do kanálu ch.
func Walk(t *tree.Tree, ch chan int)

// Same ověřuje jestli stromy
// t1 a t2 obsahují stejné hodnoty.
func Same(t1, t2 *tree.Tree) bool

func main() {
}
