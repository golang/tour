// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pic implements functions that
// display pictures on the Go playground.
package pic // import "golang.org/x/tour/pic"

import (
	"bufio"
	"encoding/base64"
	"image"
	"image/png"
	"io"
	"os"
)

// Show displays a picture defined by the function f
// when executed on the Go Playground.
//
// f should return a slice of length dy,
// each element of which is a slice of dx
// 8-bit unsigned int. The integers are
// interpreted as bluescale values,
// where the value 0 means full blue,
// and the value 255 means full white.
func Show(f func(dx, dy int) [][]uint8) {
	const (
		dx = 256
		dy = 256
	)
	data := f(dx, dy)
	m := image.NewNRGBA(image.Rect(0, 0, dx, dy))
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			v := data[y][x]
			i := y*m.Stride + x*4
			m.Pix[i] = v
			m.Pix[i+1] = v
			m.Pix[i+2] = 255
			m.Pix[i+3] = 255
		}
	}
	ShowImage(m)
}

// ShowImage displays the image m
// when executed on the Go Playground.
func ShowImage(m image.Image) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	io.WriteString(w, "IMAGE:")
	b64 := base64.NewEncoder(base64.StdEncoding, w)
	err := (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(b64, m)
	if err != nil {
		panic(err)
	}
	b64.Close()
	io.WriteString(w, "\n")
}
