// +build no-build OMIT

package main

import "fmt"

// fibonacci to funkcja która zwraca
// funkcję która zwraca int.
func fibonacci() func() int {
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
