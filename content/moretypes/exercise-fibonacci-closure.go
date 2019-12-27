// +build no-build OMIT

package main

import "fmt"

// fibonacci je funkce která vrací
// funkci která vrací int.
func fibonacci() func() int {
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
