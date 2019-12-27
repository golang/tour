// +build OMIT

package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Dvakrát %v je %v\n", v, v*2)
	case string:
		fmt.Printf("%q je %v bajtů dlouhé\n", v, len(v))
	default:
		fmt.Printf("Neznám typ %T!\n", v)
	}
}

func main() {
	do(21)
	do("ahoj")
	do(true)
}
