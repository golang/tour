// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"math/cmplx"
)

const delta = 1e-10

func Cbrt(x complex128) complex128 {
	z := x
	for {
		n := z - (z*z*z-x)/(3*z*z)
		if cmplx.Abs(n-z) < delta {
			break
		}
		z = n
	}
	return z
}

func main() {
	const x = 2
	mine, theirs := Cbrt(x), cmplx.Pow(x, 1./3.)
	fmt.Println(mine, theirs, mine-theirs)
}
