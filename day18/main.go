package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	bb, err := ioutil.ReadFile("../inputs/18.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	img := BytesImageFrom(bb)
	fmt.Println("value at minute 10:", valueOf(stateAtMinute(img, 10)))
	fmt.Println("value at minute 999:", valueOf(stateAtMinute(img, 999)))

	fmt.Println("predicted at 1000000000:", predictedValueAt(1000000000))
	fmt.Println("value at minute 1000000000:", valueOf(stateAtMinute(img, 1000000000)))
}

const (
	Trees      = '|'
	Lumberyard = '#'
	Empty      = '.'
)

func predictedValueAt(i int) int {
	if i%1000 != 0 {
		return -1
	}
	return []int{193120, 219349, 202515, 190740, 210630, 207172, 199758}[(i/1000)%7]
}

func stateAtMinute(img *BytesImage, m int) *BytesImage {
	start := time.Now()
	for i := 0; i < m; i++ {
		img = advanceMinuteMemo(img)
		if i%1000000 == 0 {
			fmt.Println(time.Since(start).Seconds(), i, valueOf(img), predictedValueAt(i))
		}
	}
	return img
}

var memoData map[string]*BytesImage

func advanceMinuteMemo(in *BytesImage) *BytesImage {
	if memoData == nil {
		memoData = make(map[string]*BytesImage)
	}

	if res, ok := memoData[string(in.data)]; ok {
		return res
	}

	res := advanceMinute(in)
	memoData[string(in.data)] = res
	return res
}

func advanceMinute(in *BytesImage) *BytesImage {
	out := NewBytesImage(in.width-1, in.height)
	in.ForEach(func(x, y int, b byte) {
		switch b {
		case Empty:
			if countNeighbourhoodOfType(in, x, y, Trees) >= 3 {
				out.Set(x, y, Trees)
			} else {
				out.Set(x, y, Empty)
			}
		case Trees:
			if countNeighbourhoodOfType(in, x, y, Lumberyard) >= 3 {
				out.Set(x, y, Lumberyard)
			} else {
				out.Set(x, y, Trees)
			}
		case Lumberyard:
			if shouldStayLumberyard(in, x, y) {
				out.Set(x, y, Lumberyard)
			} else {
				out.Set(x, y, Empty)
			}
		}
	})
	return out
}

func shouldStayLumberyard(in *BytesImage, x, y int) bool {
	var nt, nl int
	in.For8Neighbors(x, y, func(_, _ int, n byte) {
		if n == Trees {
			nt++
		} else if n == Lumberyard {
			nl++
		}
	})
	return nl > 0 && nt > 0
}

func countNeighbourhoodOfType(in *BytesImage, x, y int, t byte) int {
	var cnt int
	in.For8Neighbors(x, y, func(_, _ int, n byte) {
		if n == t {
			cnt++
		}
	})
	return cnt
}

func valueOf(img *BytesImage) int {
	n := bytes.Count(img.data, []byte{Trees})
	m := bytes.Count(img.data, []byte{Lumberyard})
	return n * m
}
