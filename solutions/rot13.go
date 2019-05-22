// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"io"
	"os"
	"strings"
)

func rot13(b byte) byte {
	const alphabetLength = 26
	var a byte
	switch {
	case 'a' <= b && b <= 'z':
		a = 'a'
	case 'A' <= b && b <= 'Z':
		a = 'A'
	default:
		return b
	}
	return (b-a+13)%alphabetLength + a
}

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = rot13(p[i])
	}
	return
}

func main() {
	s := strings.NewReader(
		"Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
