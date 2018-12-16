package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestGameFrom(t *testing.T) {
	in := `#####
#..##
#G.E#
#G#.#
#..G#
#####
`

	g := GameFrom([]byte(in))
	type pt struct {
		x, y   int
		marker byte
	}
	units := []pt{
		{1, 2, 'G'},
		{1, 3, 'G'},
		{3, 4, 'G'},
		{3, 2, 'E'},
	}
	if len(units) != len(g.Units) {
		t.Errorf("n got %d expected %d", len(g.Units), len(units))
	}
	for _, p := range units {
		t.Run(fmt.Sprintf("%d_%d", p.x, p.y), func(t *testing.T) {
			var found bool
			for _, u := range g.Units {
				if u.X == p.x && u.Y == p.y {
					found = true
					if p.marker != u.Marker {
						t.Errorf("marker got %v expected %v", u.Marker, p.marker)
					}
				}
			}
			if !found {
				t.Errorf("not found %v", p)
			}
		})
	}
}

func TestGameEdges(t *testing.T) {
	in := `#####
#..##
#.EG#
###.#
#####
`

	g := GameFrom([]byte(in))
	type uv struct {
		sx, sy, dx, dy int
	}
	exp := []uv{
		{1, 1, 2, 1},
		{1, 1, 1, 2},
		{2, 1, 1, 1},
		{2, 1, 2, 2},
		{1, 2, 1, 1},
		{1, 2, 2, 2},
		{2, 2, 2, 1},
		{2, 2, 1, 2},
		{3, 2, 3, 3},
		{3, 3, 3, 2},
	}
	r := g.Edges()
	if len(r) != len(exp) {
		t.Errorf("n got %d expected %d (%v)", len(r), len(exp), r)
	}

	for i := range exp {
		sx, sy := g.Map.XY.FromIndex(r[i].U)
		dx, dy := g.Map.XY.FromIndex(r[i].V)
		act := uv{sx, sy, dx, dy}
		if !reflect.DeepEqual(act, exp[i]) {
			t.Errorf("%d got %v expected %v", i, act, exp[i])
		}
	}
}

func TestXYEncoderToIndex(t *testing.T) {
	in := `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
`
	g := GameFrom([]byte(in))
	var i int
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			r := g.Map.XY.ToIndex(x, y)
			if r != i {
				t.Errorf("%d,%d got %d expected %d", x, y, r, i)
			}
			i++
		}
	}
}

func TestXYEncoderFromIndex(t *testing.T) {
	in := `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
`
	g := GameFrom([]byte(in))
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			i := g.Map.XY.ToIndex(x, y)
			rx, ry := g.Map.XY.FromIndex(i)
			if rx != x || ry != y {
				t.Errorf("%d got %d,%d expected %d,%d", i, rx, ry, x, y)
			}
		}
	}
}

func TestGameOutcome(t *testing.T) {
	if testing.Short() {
		t.Skip("testing in short mode")
	}

	in := `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
`
	g := GameFrom([]byte(in))
	for i := 0; g.RunRound(); i++ {
		if i < 3 {
			fmt.Println(i, "-----------")
			fmt.Println(string(g.Map.data))
			for _, u := range g.Units {
				fmt.Println(*u)
			}
		}
		if i > 50 {
			t.Fail()
			return
		}
	}
	exp := 27730
	act := g.Outcome()
	if act != exp {
		t.Errorf("got %d expected %d", act, exp)
	}
}

func TestPathCalculatorDist(t *testing.T) {
	in := `#####
#..##
#G.E#
#G#.#
#...#
#####`

	g := GameFrom([]byte(in))
	p := PathCalculatorFrom(g)
	start := g.Map.XY.ToIndex(3, 2)
	img := NewBytesImage(5, 6)
	fmt.Println(img.data)
	img.ForEach(func(x, y int, _ byte) {
		vi := g.Map.XY.ToIndex(x, y)
		v := p.Dist(start, vi)
		if y == 2 {
			fmt.Println(x, y, start, vi, v)
		}
		vs := strconv.Itoa(v)
		if v > 9 {
			vs = "-"
		}
		img.Set(x, y, []byte(vs)[0])
	})
	fmt.Println(string(g.Map.data))
	fmt.Println(string(img.data))
	t.Fail()
}

func TestBytesImageFor4Neighbors(t *testing.T) {
	in := `#####
#..##
#G.E#
#G#.#
#...#
#####
`
	bi := BytesImageFrom([]byte(in))

	verify := func(x, y int, exp string) {
		var nn []byte
		bi.For4Neighbors(x, y, func(px, py int, b byte) {
			nn = append(nn, b)
		})
		if string(nn) != exp {
			t.Errorf("(%d,%d) got %v expected %v", x, y, string(nn), exp)
		}
	}

	verify(0, 0, "##")
	verify(1, 1, "##.G")
	verify(4, 5, "##")
}

func TestBytesImageFor8Neighbors(t *testing.T) {
	in := `#####
#..##
#G.E#
#G#.#
#...#
#####
`
	bi := BytesImageFrom([]byte(in))

	verify := func(x, y int, exp string) {
		var nn []byte
		bi.For8Neighbors(x, y, func(px, py int, b byte) {
			nn = append(nn, b)
		})
		if string(nn) != exp {
			t.Errorf("(%d,%d) got %v expected %v", x, y, string(nn), exp)
		}
	}

	verify(0, 0, "##.")
	verify(1, 1, "####.#G.")
	verify(4, 5, ".##")
}
