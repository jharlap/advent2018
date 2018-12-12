package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	bb, err := ioutil.ReadFile("../inputs/12.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	g := NewGame()
	ll := bytes.Split(bb, []byte("\n"))
	for _, l := range ll {
		if len(l) < 10 {
			continue
		}
		i, o := l[:5], l[9]
		if o == '#' {
			g.AddRule(string(i), o)
		}
	}
	//                 1         2         3
	//       0         0         0         0
	//20: .#....##....#####...#######....#.#..##.
	//-2+3+4+9+10+11+12+13+17+18+19+20+21+22+23+28+30+33+34

	in := "##.#..########..##..#..##.....##..###.####.###.##.###...###.##..#.##...#.#.#...###..###.###.#.#"
	s := State{
		State: []byte(in),
	}
	for i := 1; i <= 50000000000; i++ {
		s = g.NextGeneration(s)
		if i == 20 || i == 50000000000 || i%1000000 == 0 {
			fmt.Println("i:", i, "state:", string(s.State), "off:", s.Offset, "score:", s.Score())
		}
	}
}

type Game struct {
	rules map[string]byte
}

type State struct {
	State  []byte
	Offset int
}

func NewGame() *Game {
	return &Game{
		rules: make(map[string]byte),
	}
}

func (g *Game) AddRule(in string, out byte) {
	g.rules[in] = out
}

func (g *Game) NextGeneration(s State) State {
	var r State
	nibbles := s.Nibbles()
	r.Offset = s.Offset - 2
	for _, w := range nibbles {
		c := g.rules[string(w)]
		if c != '#' {
			r.State = append(r.State, '.')
		} else {
			r.State = append(r.State, '#')
		}
	}
	return trim(r)
}

func (s State) Nibbles() [][]byte {
	r := make([][]byte, 0, len(s.State)+8)
	r = append(r, []byte{'.', '.', '.', '.', s.State[0]})
	r = append(r, []byte{'.', '.', '.', s.State[0], s.State[1]})
	r = append(r, []byte{'.', '.', s.State[0], s.State[1], s.State[2]})
	r = append(r, []byte{'.', s.State[0], s.State[1], s.State[2], s.State[3]})

	for i := 0; i < len(s.State)-4; i++ {
		r = append(r, s.State[i:i+5])
	}

	j := len(s.State)
	r = append(r, []byte{s.State[j-4], s.State[j-3], s.State[j-2], s.State[j-1], '.'})
	r = append(r, []byte{s.State[j-3], s.State[j-2], s.State[j-1], '.', '.'})
	r = append(r, []byte{s.State[j-2], s.State[j-1], '.', '.', '.'})
	r = append(r, []byte{s.State[j-1], '.', '.', '.', '.'})

	return r
}

func trim(s State) State {
	bb := bytes.TrimLeft(s.State, ".")
	s.Offset += len(s.State) - len(bb)
	s.State = bytes.TrimRight(bb, ".")
	return s
}

func (s State) Score() int {
	var r int
	for i := 0; i < len(s.State); i++ {
		if s.State[i] == '#' {
			r += i + s.Offset
		}
	}
	return r
}
