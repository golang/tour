package main

import (
	"image"
	"image/color"
	"tour/pic"
)

type Image [][]uint8

func (m Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (m Image) Bounds() image.Rectangle {
	dx, dy := len(m), 0
	if dx > 0 {
		dy = len(m[0])
	}

	return image.Rect(0, 0, dx, dy)
}

func (m Image) At(x, y int) color.Color {
	v := m[x][y]
	return color.RGBA{v, v, 255, 255}
}

func NewImage(x, y int, f func(int, int) uint8) Image {
	m := make([][]uint8, x)
	for i := range m {
		m[i] = make([]uint8, y)
		for j := range m[i] {
			m[i][j] = f(i, j)
		}
	}
	return m
}

func main() {
	m := NewImage(255, 255, func(a, b int) uint8 { return uint8(a * b) })
	pic.ShowImage(m)
}
