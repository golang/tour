// +build OMIT

package main

import (
	"code.google.com/p/go-tour/pic"
	"image"
)

type Image struct{}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
