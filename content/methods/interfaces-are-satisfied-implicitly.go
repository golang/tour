// +build OMIT

package main

import "fmt"

type I interface {
	M()
}

type T struct {
	S string
}

// Tato metoda znamená že typ T implementuje interface I,
// ale nikde nemusíme explicitně deklarovat, že to tak je.
func (t T) M() {
	fmt.Println(t.S)
}

func main() {
	var i I = T{"hello"}
	i.M()
}
