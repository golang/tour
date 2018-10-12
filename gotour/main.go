// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
)

func main() {
	io.WriteString(os.Stderr, "golang.org/x/tour/gotour has moved to golang.org/x/tour\n")
	os.Exit(1)
}
