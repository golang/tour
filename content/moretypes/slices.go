// +build OMIT

package main

import "fmt"

func main() {
	primes := []int{2, 3, 5, 7, 11, 13}
	fmt.Println("primes ==", primes)

	for i := 0; i < len(primes); i++ {
		fmt.Printf("primes[%d] == %d\n", i, primes[i])
	}
}
