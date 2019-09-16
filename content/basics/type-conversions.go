// +build OMIT

package main

import (
	"fmt"
	"math"
)

func main() {
	var x, y int = 4, 5
	var f float64 = math.Sqrt(float64(x*x + y*y))
	fmt.Println(f)
	var z uint = uint(f)
	fmt.Println(x, y, z)
}
