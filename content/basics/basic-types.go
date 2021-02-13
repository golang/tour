// +build OMIT

package main

import (
	"fmt"
	"math/cmplx"
)

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

func main() {
	fmt.Printf("Typ: %T Wartość: %v\n", ToBe, ToBe)
	fmt.Printf("Typ: %T Wartość: %v\n", MaxInt, MaxInt)
	fmt.Printf("Typ: %T Wartość: %v\n", z, z)
}
