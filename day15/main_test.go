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

func TestGameMoveUnit(t *testing.T) {
	in := `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######
`
	g := GameFrom([]byte(in))
	var u *Unit
	for _, v := range g.Units {
		if v.X == 2 && v.Y == 1 {
			u = v
			break
		}
	}

	moved := g.moveUnit(u)
	if !moved {
		t.Error("did not move")
	}

	if u.X == 2 && u.Y == 1 {
		t.Errorf("thinks it moved but didn't: %v\n", u)
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

	cc := []struct {
		in, ex         string
		round, outcome int
	}{
		{"#######\n#G..#E#\n#E#E.E#\n#G.##.#\n#...#E#\n#...E.#\n#######\n",
			"#######\n#...#E#\n#E#...#\n#.E##.#\n#E..#E#\n#.....#\n#######\n",
			37, 36334},
		{"#######\n#.G...#\n#...EG#\n#.#.#G#\n#..G#E#\n#.....#\n#######\n",
			"#######\n#.G...#\n#...EG#\n#.#.#G#\n#..G#E#\n#.....#\n#######\n",
			47, 27730},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := GameFrom([]byte(c.in))
			for g.RunRound() {
			}
			act := g.Outcome()
			if act != c.outcome {
				t.Errorf("outcome got %d expected %d", act, c.outcome)
			}
			if g.Round != c.round {
				t.Errorf("round got %d expected %d", g.Round, c.round)
			}
		})
	}
}

func TestPathCalculatorDist(t *testing.T) {
	in := `#####
#..##
#G.E#
#G#.#
#G..#
#####`
	exp := `-----
-32--
-21--
---1-
-432-
-----
`

	g := GameFrom([]byte(in))
	p := PathCalculatorFrom(g)
	start := g.Map.XY.ToIndex(3, 2)
	img := NewBytesImage(5, 6)
	img.ForEach(func(x, y int, _ byte) {
		vi := g.Map.XY.ToIndex(x, y)
		v := p.Dist(start, vi)
		vs := strconv.Itoa(v)
		if v > 9 {
			vs = "-"
		}
		img.Set(x, y, []byte(vs)[0])
	})

	if string(img.data) != exp {
		t.Errorf("got %d\n%v\nexpected %d\n%v\n", len(img.data), string(img.data), len(exp), exp)
	}
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
