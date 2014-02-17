// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"math"
)

const delta = 1e-10

func Sqrt(x float64) float64 {
	z := x
	for {
		n := z - (z*z-x)/(2*z)
		if math.Abs(n-z) < delta {
			break
		}
		z = n
	}
	return z
}

func main() {
	const x = 2
	mine, theirs := Sqrt(x), math.Sqrt(x)
	fmt.Println(mine, theirs, mine-theirs)
}
