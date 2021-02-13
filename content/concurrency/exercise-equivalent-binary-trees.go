// +build no-build OMIT

package main

import "golang.org/x/tour/tree"

// Walk przechodzi przez drzewko t oraz wysyła wszystkie wartości
// z drzewka do kanału ch.
func Walk(t *tree.Tree, ch chan int)

// Same określa czy drzewko t1 oraz t2
// posiadają te same wartości.
func Same(t1, t2 *tree.Tree) bool

func main() {
}
