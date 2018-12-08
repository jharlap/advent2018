package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rr, err := readInput("../inputs/3.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var ii []Rect
	for i, a := range rr {
		zeroIntersections := true
		for j, b := range rr {
			if i == j {
				continue
			}

			ok, r := intersection(a, b)
			if ok {
				zeroIntersections = false
				ii = append(ii, r)
			}
		}
		if zeroIntersections {
			fmt.Println("zero intersections:", a.id)
		}
	}

	m := make(map[int]bool)
	for _, r := range ii {
		for x := r.x; x < r.x+r.w; x++ {
			for y := r.y; y < r.y+r.h; y++ {
				k := x*1000000 + y
				m[k] = true
			}
		}
	}
	fmt.Println("Found intersections:", len(m))
}

type Rect struct {
	x, y, w, h int
	id         int
}

func intersection(a, b Rect) (bool, Rect) {
	if a.x+a.w <= b.x { // a is left of b
		return false, Rect{}
	} else if b.x+b.w <= a.x { // b is left of a
		return false, Rect{}
	} else if a.y+a.h <= b.y { // a is above b
		return false, Rect{}
	} else if b.y+b.h <= a.y { // b is above a
		return false, Rect{}
	}

	x := max(a.x, b.x)
	y := max(a.y, b.y)
	w := min(a.x+a.w, b.x+b.w) - x
	h := min(a.y+a.h, b.y+b.h) - y

	return true, Rect{x: x, y: y, w: w, h: h}
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

func intersects(a, b Rect) bool {
	if a.x+a.w <= b.x { // a is left of b
		return false
	} else if b.x+b.w <= a.x { // b is left of a
		return false
	} else if a.y+a.h <= b.y { // a is above b
		return false
	} else if b.y+b.h <= a.y { // b is above a
		return false
	}
	return true
}

func readInput(filename string) ([]Rect, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input: %v", err)
	}
	defer r.Close()

	var rr []Rect
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		var (
			r Rect
		)
		// read input like "#14 @ 851,648: 13x15"
		_, err := fmt.Sscanf(l, "#%d @ %d,%d: %dx%d", &r.id, &r.x, &r.y, &r.w, &r.h)
		if err != nil {
			return nil, fmt.Errorf("error parsing line '%s': %v", l, err)
		}
		rr = append(rr, r)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}
	return rr, nil
}
