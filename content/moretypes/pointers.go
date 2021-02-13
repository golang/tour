// +build OMIT

package main

import "fmt"

func main() {
	i, j := 42, 2701

	p := &i         // wkaźnik do i
	fmt.Println(*p) // przeczytaj i poprzez wkaźnik
	*p = 21         // ustaw wartość i poprzez wkaźnik
	fmt.Println(i)  // zobacz nową wartość i

	p = &j         // wskaźnik do j
	*p = *p / 37   // podziel j za pomocą wskaźnika
	fmt.Println(j) // zobacz nowa wartość j
}
