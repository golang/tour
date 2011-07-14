// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pic

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"fmt"
)

func Show(f func(int, int)[][]uint8) {
	const (	
		dx = 256
		dy = 256
	)
	data := f(dx, dy)
	m := image.NewNRGBA(dx, dy)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			v := data[y][x]
			m.Pix[y*dy+x] = image.NRGBAColor{v, v, 255, 255}
		}
	}
	ShowImage(m)
}

func ShowImage(m image.Image) {
	var buf bytes.Buffer
	png.Encode(&buf, m)
	enc := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(enc, buf.Bytes())
	fmt.Println("IMAGE:" + string(enc))
}
