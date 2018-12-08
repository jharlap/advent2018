package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMinRect(t *testing.T) {
	cc := []struct {
		vv []Vertex
		r  Vertex
	}{
		{[]Vertex{{1, 1}, {1, 6}, {8, 3}, {3, 4}, {5, 5}, {8, 9}}, Vertex{9, 10}},
		{[]Vertex{}, Vertex{1, 1}},
		{[]Vertex{{2, 3}}, Vertex{3, 4}},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			r := minRect(c.vv)
			if !reflect.DeepEqual(r, c.r) {
				t.Errorf("got %v expected %v", r, c.r)
			}
		})
	}
}

func TestLabelInRange(t *testing.T) {
	cc := []struct {
		vv  []Vertex
		rng int
		ll  [][]int
	}{
		{
			[]Vertex{{1, 1}, {1, 6}, {8, 3}, {3, 4}, {5, 5}, {8, 9}},
			32,
			[][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // .........
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // .A.......
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // .........
				{0, 0, 0, 1, 1, 1, 0, 0, 0}, // ...###..C
				{0, 0, 1, 1, 1, 1, 1, 0, 0}, // ..#D###..
				{0, 0, 1, 1, 1, 1, 1, 0, 0}, // ..###E#..
				{0, 0, 0, 1, 1, 1, 0, 0, 0}, // .B.###...
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // .........
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // .........
				{0, 0, 0, 0, 0, 0, 0, 0, 0}, // ........F
			},
		},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			rect := minRect(c.vv)
			r := labelInRange(rect, c.rng, c.vv)
			if !reflect.DeepEqual(r, c.ll) {
				t.Errorf("delta: %v", delta(rect, r, c.ll))
			}
		})
	}
}

func TestLabel(t *testing.T) {
	cc := []struct {
		vv []Vertex
		ll [][]int
	}{
		{
			[]Vertex{{1, 1}, {1, 6}, {8, 3}, {3, 4}, {5, 5}, {8, 9}},
			[][]int{
				{0, 0, 0, 0, 0, -1, 2, 2, 2},  // aaaaa.ccc
				{0, 0, 0, 0, 0, -1, 2, 2, 2},  // aAaaa.ccc
				{0, 0, 0, 3, 3, 4, 2, 2, 2},   // aaaddeccc
				{0, 0, 3, 3, 3, 4, 2, 2, 2},   // aadddeccC
				{-1, -1, 3, 3, 3, 4, 4, 2, 2}, // ..dDdeecc
				{1, 1, -1, 3, 4, 4, 4, 4, 2},  // bb.deEeec
				{1, 1, 1, -1, 4, 4, 4, 4, -1}, // bBb.eeee.
				{1, 1, 1, -1, 4, 4, 4, 5, 5},  // bbb.eeeff
				{1, 1, 1, -1, 4, 4, 5, 5, 5},  // bbb.eefff
				{1, 1, 1, -1, 5, 5, 5, 5, 5},  // bbb.ffffF
			},
		},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			rect := minRect(c.vv)
			r := label(rect, c.vv)
			if !reflect.DeepEqual(r, c.ll) {
				t.Errorf("delta: %v", delta(rect, r, c.ll))
			}
		})
	}
}

func TestEdgeValues(t *testing.T) {
	cc := []struct {
		vv []Vertex
		ll map[int]bool
	}{
		{
			[]Vertex{{1, 1}, {1, 6}, {8, 3}, {3, 4}, {5, 5}, {8, 9}},
			map[int]bool{-1: true, 0: true, 1: true, 2: true, 5: true},
		},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			rect := minRect(c.vv)
			ll := label(rect, c.vv)
			r := edgeValues(rect, ll)
			if !reflect.DeepEqual(r, c.ll) {
				t.Errorf("got: %v expected: %v", r, c.ll)
			}
		})
	}
}

type Diff struct {
	x, y, a, b int
}

func delta(r Vertex, a, b [][]int) []Diff {
	var d []Diff
	for y := 0; y < r.y; y++ {
		for x := 0; x < r.x; x++ {
			if a[y][x] != b[y][x] {
				d = append(d, Diff{x: x, y: y, a: a[y][x], b: b[y][x]})
			}
		}
	}
	return d
}
