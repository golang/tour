// +build OMIT

package main

import "fmt"

func main() {
	fmt.Println("liczę")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("zrobione")
}
