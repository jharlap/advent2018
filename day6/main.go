package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	vv, err := readInput("../inputs/6.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rect := minRect(vv)
	img := label(rect, vv)
	ignore := edgeValues(rect, img)
	counts := count(img)

	var lk, lv int
	for k, v := range counts {
		if ignore[k] {
			continue
		}

		if v > lv {
			lv = v
			lk = k
		}
	}
	fmt.Println("k:", lk, "v:", lv)

	maxDist := 10000
	img = labelInRange(rect, maxDist, vv)
	counts = count(img)
	fmt.Println("region:", counts[1])
}

func count(vv [][]int) map[int]int {
	r := make(map[int]int)
	for _, v := range vv {
		for _, l := range v {
			r[l]++
		}
	}
	return r
}

func edgeValues(rect Vertex, vv [][]int) map[int]bool {
	m := make(map[int]bool)
	for i := 0; i < rect.x; i++ {
		m[vv[0][i]] = true
		m[vv[rect.y-1][i]] = true
	}
	for i := 0; i < rect.y; i++ {
		m[vv[i][0]] = true
		m[vv[i][rect.x-1]] = true
	}
	return m
}

func labelInRange(rect Vertex, maxDist int, vv []Vertex) [][]int {
	ll := make([][]int, rect.y)
	for i := 0; i < rect.y; i++ {
		ll[i] = make([]int, rect.x)
	}

	for y := 0; y < rect.y; y++ {
		for x := 0; x < rect.x; x++ {

			c := Vertex{x: x, y: y}
			var s int
			for _, v := range vv {
				s += dist(c, v)
			}

			if s < maxDist {
				ll[y][x] = 1
			}
		}
	}
	return ll
}

func label(rect Vertex, vv []Vertex) [][]int {
	ll := make([][]int, rect.y)
	for i := 0; i < rect.y; i++ {
		ll[i] = make([]int, rect.x)
	}

	maxDist := rect.x*rect.y + 1
	for y := 0; y < rect.y; y++ {
		for x := 0; x < rect.x; x++ {
			d := maxDist
			l := -2

			c := Vertex{x: x, y: y}
			for i, v := range vv {
				cd := dist(c, v)
				if cd < d {
					d = cd
					l = i
				} else if cd == d {
					l = -1
				}

				if false && x == 2 && y == 3 {
					fmt.Println("c:", c, "i:", i, "v:", v, "cd:", cd, "d:", d, "l:", l)
				}
			}

			ll[y][x] = l
		}
	}
	return ll
}

func minRect(vv []Vertex) Vertex {
	var r Vertex
	for _, v := range vv {
		if v.x > r.x {
			r.x = v.x
		}
		if v.y > r.y {
			r.y = v.y
		}
	}
	r.x++
	r.y++
	return r
}

type Vertex struct {
	x, y int
}

func dist(a, b Vertex) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func readInput(filename string) ([]Vertex, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filename, err)
	}

	var vv []Vertex
	for _, l := range strings.Split(string(bb), "\n") {
		var v Vertex
		_, err := fmt.Sscanf(l, "%d, %d", &v.x, &v.y)
		if err != nil {
			return nil, fmt.Errorf("error scanning line '%s': %v", l, err)
		}
		vv = append(vv, v)
	}
	return vv, nil
}
