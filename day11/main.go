package main

import (
	"fmt"
)

func main() {
	grid := 7672
	var mp, mx, my, ms int
	s := 3
	for x := 1; x <= 300-s; x++ {
		for y := 1; y <= 300-s; y++ {
			p := areaPowerLevel(x, y, 3, grid)
			if p > mp {
				mp = p
				mx = x
				my = y
			}
		}
	}
	fmt.Printf("3x3: %d,%d p: %d\n", mx, my, mp)

	for s := 1; s < 300; s++ {
		for x := 1; x <= 300-s; x++ {
			for y := 1; y <= 300-s; y++ {
				p := areaPowerLevel(x, y, s, grid)
				if p > mp {
					mp = p
					mx = x
					my = y
					ms = s
				}
			}
		}
	}
	fmt.Printf("%d,%d,%d p: %d\n", mx, my, ms, mp)
}

func powerLevel(x, y, grid int) int {
	id := x + 10
	p := (id*y + grid) * id
	p = (p / 100) % 10
	return p - 5
}

func areaPowerLevel(x, y, size, grid int) int {
	var p int
	for cx := x; cx < x+size; cx++ {
		for cy := y; cy < y+size; cy++ {
			p += powerLevel(cx, cy, grid)
		}
	}
	return p
}
