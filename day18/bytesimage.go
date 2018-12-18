package main

import (
	"bytes"
)

const Newline = '\n'

type BytesImage struct {
	data          []byte
	width, height int
	XY            XYEncoder
}

func NewBytesImage(w, h int) *BytesImage {
	bi := &BytesImage{
		data:   make([]byte, (w+1)*h),
		width:  w + 1,
		height: h,
		XY:     XYEncoder{Width: w + 1, Height: h},
	}

	for i := w; i < len(bi.data); i += w + 1 {
		bi.data[i] = Newline
	}
	return bi
}

func BytesImageFrom(data []byte) *BytesImage {
	w := bytes.IndexRune(data, '\n') + 1
	h := bytes.Count(data, []byte{byte('\n')}) + 1
	return &BytesImage{
		data:   data[:],
		width:  w,
		height: h,
		XY:     XYEncoder{Width: w, Height: h},
	}
}

func (im *BytesImage) Set(x, y int, b byte) {
	i := im.XY.ToIndex(x, y)
	if i < 0 || i > len(im.data) {
		return
	}
	if im.data[i] == Newline {
		return
	}
	im.data[i] = b
}

func (im *BytesImage) ValueAt(x, y int) byte {
	i := im.XY.ToIndex(x, y)
	if i < 0 || i >= len(im.data) {
		return 0
	}
	b := im.data[i]
	if b == Newline {
		return 0
	}
	return b
}

func (im *BytesImage) ForEach(f func(x, y int, b byte)) {
	for i := range im.data {
		if im.data[i] == Newline {
			continue
		}

		x, y := im.XY.FromIndex(i)
		f(x, y, im.data[i])
	}
}

func (im *BytesImage) ForRow(y int, f func(x, y int, b byte)) {
	for x := 0; x < im.width; x++ {
		v := im.ValueAt(x, y)
		if v > 0 {
			f(x, y, v)
		}
	}
}

func (im *BytesImage) ForFirstUpFrom(x, y int, f func(x, y int, b byte) bool) {
	for cy := y; cy >= 0; cy-- {
		v := im.ValueAt(x, cy)
		if v > 0 {
			f(x, cy, v)
			return
		}
	}
}

func (im *BytesImage) ForEachDownFrom(x, y int, f func(x, y int, b byte) bool) {
	var done bool
	for cy := y; !done && cy <= im.height; cy++ {
		v := im.ValueAt(x, cy)
		if v > 0 {
			done = f(x, cy, v)
		}
	}
}

func (im *BytesImage) For4Neighbors(x, y int, f func(x, y int, b byte)) {
	nn := []struct{ x, y int }{{x, y - 1}, {x - 1, y}, {x + 1, y}, {x, y + 1}}
	for _, n := range nn {
		if n.x < 0 || n.y < 0 || n.x >= im.width || n.y >= im.height {
			continue
		}
		v := im.ValueAt(n.x, n.y)
		if v > 0 {
			f(n.x, n.y, v)
		}
	}
}

func (im *BytesImage) For8Neighbors(x, y int, f func(x, y int, b byte)) {
	for j := max(y-1, 0); j <= min(y+1, im.height-1); j++ {
		for i := max(x-1, 0); i <= min(x+1, im.width-1); i++ {
			if i == x && j == y {
				continue
			}
			v := im.ValueAt(i, j)
			if v > 0 {
				f(i, j, v)
			}
		}
	}
}

func (im *BytesImage) String() string {
	return string(im.data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type XYEncoder struct {
	Width, Height int
}

func (e XYEncoder) ToIndex(x, y int) int {
	return y*e.Width + x
}

func (e XYEncoder) FromIndex(i int) (int, int) {
	y := i / e.Width
	x := i - (y * e.Width)
	return x, y
}
