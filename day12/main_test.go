package main

import (
	"testing"
)

func TestScore(t *testing.T) {
	cc := []struct {
		s   string
		o   int
		exp int
	}{
		{".....", 0, 0},
		{"#..#", -2, -1},
		{"##.#..#####", 0, 44},
		{"##.#..#####", 1, 52},
		{"##.#..#####", -1, 36},
		{".#....##....#####...#######....#.#..##.", -3, 325},
	}

	for i, c := range cc {
		s := State{
			State:  []byte(c.s),
			Offset: c.o,
		}
		act := s.Score()
		if act != c.exp {
			t.Errorf("case %d got: %d expected %d", i, act, c.exp)
		}
	}
}

func TestNextGeneration(t *testing.T) {
	g := NewGame()
	g.AddRule("...##", '#')
	g.AddRule("..#..", '#')
	g.AddRule(".#...", '#')
	g.AddRule(".#.#.", '#')
	g.AddRule(".#.##", '#')
	g.AddRule(".##..", '#')
	g.AddRule(".####", '#')
	g.AddRule("#.#.#", '#')
	g.AddRule("#.###", '#')
	g.AddRule("##.#.", '#')
	g.AddRule("##.##", '#')
	g.AddRule("###..", '#')
	g.AddRule("###.#", '#')
	g.AddRule("####.", '#')

	gens := []State{
		State{[]byte("#..#.#..##......###...###"), 0},
		State{[]byte("#...#....#.....#..#..#..#"), 0},
		State{[]byte("##..##...##....#..#..#..##"), 0},
		State{[]byte("#.#...#..#.#....#..#..#...#"), -1},
		State{[]byte("#.#..#...#.#...#..#..##..##"), 0},
		State{[]byte("#...##...#.#..#..#...#...#"), 1},
		State{[]byte("##.#.#....#...#..##..##..##"), 1},
		State{[]byte("#..###.#...##..#...#...#...#"), 0},
		State{[]byte("#....##.#.#.#..##..##..##..##"), 0},
		State{[]byte("##..#..#####....#...#...#...#"), 0},
		State{[]byte("#.#..#...#.##....##..##..##..##"), -1},
		State{[]byte("#...##...#.#...#.#...#...#...#"), 0},
		State{[]byte("##.#.#....#.#...#.#..##..##..##"), 0},
		State{[]byte("#..###.#....#.#...#....#...#...#"), -1},
		State{[]byte("#....##.#....#.#..##...##..##..##"), -1},
		State{[]byte("##..#..#.#....#....#..#.#...#...#"), -1},
		State{[]byte("#.#..#...#.#...##...#...#.#..##..##"), -2},
		State{[]byte("#...##...#.#.#.#...##...#....#...#"), -1},
		State{[]byte("##.#.#....#####.#.#.#...##...##..##"), -1},
		State{[]byte("#..###.#..#.#.#######.#.#.#..#.#...#"), -2},
		State{[]byte("#....##....#####...#######....#.#..##"), -2},
	}

	s := gens[0]
	for i := 1; i < len(gens); i++ {
		r := g.NextGeneration(s)
		if string(r.State) != string(gens[i].State) {
			t.Errorf("gen %d got %s expected %s", i, string(r.State), gens[i].State)
		}
		if r.Offset != gens[i].Offset {
			t.Errorf("gen %d got %d expected %d", i, r.Offset, gens[i].Offset)
		}
		s = r
	}
}

func TestTrim(t *testing.T) {
	cc := []struct {
		in, out State
	}{
		{State{[]byte("#"), 0}, State{[]byte("#"), 0}},
		{State{[]byte("....#"), -2}, State{[]byte("#"), 2}},
		{State{[]byte("..#.."), -2}, State{[]byte("#"), 0}},
	}

	for _, c := range cc {
		r := trim(c.in)
		if string(r.State) != string(c.out.State) {
			t.Errorf("%s: got %s expected %s", string(c.in.State), string(r.State), c.out.State)
		}
		if r.Offset != c.out.Offset {
			t.Errorf("%s: got %d expected %d", string(c.in.State), r.Offset, c.out.Offset)
		}
	}
}
