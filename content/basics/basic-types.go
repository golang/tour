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
	fmt.Printf("Typ: %T Hodnota: %v\n", ToBe, ToBe)
	fmt.Printf("Typ: %T Hodnota: %v\n", MaxInt, MaxInt)
	fmt.Printf("Typ: %T Hodnota: %v\n", z, z)
}
