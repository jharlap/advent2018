package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	bb, err := ioutil.ReadFile("../inputs/13.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	m := MapFrom(bb)
	crash := m.FirstCrash(nil)
	fmt.Println("first crash:", crash)

	m = MapFrom(bb)
	pos := m.LastCart()
	fmt.Println("Last pos:", pos)
}

type XY struct {
	X, Y int
}

func (p XY) Equals(v XY) bool {
	return p.X == v.X && p.Y == v.Y
}

type Turn int

const (
	TurnLeft Turn = iota
	TurnStraight
	TurnRight
)

var (
	DirEmpty  = XY{0, 0}
	DirUp     = XY{0, -1}
	DirDown   = XY{0, 1}
	DirLeft   = XY{-1, 0}
	DirRight  = XY{1, 0}
	turnTheta = []int{-1, 0, 1}
)

type Cart struct {
	Pos      XY
	Dir      XY
	Crashed  bool
	nextTurn Turn
}

func NewCart(p XY, d XY) *Cart {
	return &Cart{
		Pos: p,
		Dir: d,
	}
}

func (c *Cart) Turn() {
	if c.nextTurn != TurnStraight {
		sinTheta := turnTheta[c.nextTurn]
		c.Dir.X, c.Dir.Y = -c.Dir.Y*sinTheta, c.Dir.X*sinTheta
	}
	c.nextTurn = (c.nextTurn + 1) % 3
}

func (c *Cart) stringerChar() byte {
	if c.Crashed {
		return 'X'
	}
	switch c.Dir {
	case DirUp:
		return '^'
	case DirDown:
		return 'v'
	case DirRight:
		return '>'
	case DirLeft:
		return '<'
	}
	return '?'
}

type Track int

const (
	TrackEmpty Track = iota
	TrackHorizontal
	TrackVertical
	TrackIntersection
	TrackCurveTopLeftBottomRight // \
	TrackCurveBottomLeftTopRight // /
)

type Map struct {
	m [][]Track
	c []*Cart
}

func MapFrom(in []byte) *Map {
	bb := bytes.Split(in, []byte("\n"))
	var maxWidth int
	for _, row := range bb {
		w := len(row)
		if w > maxWidth {
			maxWidth = w
		}
	}

	m := make([][]Track, len(bb))
	var c []*Cart
	for y := 0; y < len(bb); y++ {
		m[y] = make([]Track, maxWidth)
		for x := 0; x < len(bb[y]); x++ {
			t, cartDir := TrackFrom(bb[y][x])
			m[y][x] = t

			if cartDir != DirEmpty {
				c = append(c, NewCart(XY{x, y}, cartDir))
			}
		}
	}
	return &Map{
		m: m,
		c: c,
	}
}

func (m *Map) String() string {
	buf := new(bytes.Buffer)
	for y, row := range m.m {
		for x, t := range row {
			var b byte
			for _, c := range m.c {
				if c.Pos.Equals(XY{x, y}) {
					b = c.stringerChar()
					break
				}
			}
			if b == 0 {
				switch t {
				case TrackCurveBottomLeftTopRight:
					b = '/'
				case TrackCurveTopLeftBottomRight:
					b = '\\'
				case TrackHorizontal:
					b = '-'
				case TrackVertical:
					b = '|'
				case TrackIntersection:
					b = '+'
				default:
					b = ' '
				}
			}
			buf.WriteByte(b)
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (m *Map) LastCart() *XY {
	var i int
	for len(m.c) > 1 {
		crash := m.Tick()
		if false && i > 3771 {
			fmt.Println(m)
		}
		if crash != nil {
			var r []*Cart
			for _, c := range m.c {
				if !c.Crashed {
					r = append(r, c)
				}
			}
			m.c = r
		}
		i++
	}

	fmt.Println("last cart ended on i:", i-1)
	if len(m.c) == 1 {
		return &m.c[0].Pos
	}
	return nil
}

func (m *Map) FirstCrash(shouldContinue func(int, *Map) bool) XY {
	var (
		i     int
		crash *XY
	)
	for crash == nil {
		crash = m.Tick()
		if shouldContinue != nil && !shouldContinue(i, m) {
			return XY{-1, -1}
		}
		i++
	}
	return *crash
}

// Tick advances the map by one tick and returns the coordinate of the first crash if one occurred
func (m *Map) Tick() *XY {
	var crashPos *XY

	sort.Slice(m.c, func(i, j int) bool {
		a := m.c[i].Pos
		b := m.c[j].Pos
		if a.Y < b.Y {
			return true
		} else if a.Y == b.Y && a.X < b.X {
			return true
		}
		return false
	})
	for i, c := range m.c {
		c.Pos = XY{
			X: c.Pos.X + c.Dir.X,
			Y: c.Pos.Y + c.Dir.Y,
		}

		for j, oc := range m.c {
			if i == j {
				continue
			}

			if c.Pos.Equals(oc.Pos) {
				crashPos = &c.Pos
				c.Crashed = true
				oc.Crashed = true
			}
		}

		t := m.TrackAt(c.Pos)
		switch t {
		case TrackIntersection:
			c.Turn()
		case TrackCurveBottomLeftTopRight:
			switch c.Dir {
			case DirDown:
				c.Dir = DirLeft
			case DirUp:
				c.Dir = DirRight
			case DirLeft:
				c.Dir = DirDown
			case DirRight:
				c.Dir = DirUp
			}
		case TrackCurveTopLeftBottomRight:
			switch c.Dir {
			case DirDown:
				c.Dir = DirRight
			case DirUp:
				c.Dir = DirLeft
			case DirLeft:
				c.Dir = DirUp
			case DirRight:
				c.Dir = DirDown
			}
		}
	}

	return crashPos
}

func (m *Map) Carts() []*Cart {
	return m.c
}

func (m Map) TrackAt(pos XY) Track {
	return m.m[pos.Y][pos.X]
}

func TrackFrom(b byte) (Track, XY) {
	switch b {
	case '-':
		return TrackHorizontal, DirEmpty
	case '|':
		return TrackVertical, DirEmpty
	case '/':
		return TrackCurveBottomLeftTopRight, DirEmpty
	case '\\':
		return TrackCurveTopLeftBottomRight, DirEmpty
	case '+':
		return TrackIntersection, DirEmpty
	case '>':
		return TrackHorizontal, DirRight
	case '<':
		return TrackHorizontal, DirLeft
	case 'v':
		return TrackVertical, DirDown
	case '^':
		return TrackVertical, DirUp
	}
	return TrackEmpty, DirEmpty
}
