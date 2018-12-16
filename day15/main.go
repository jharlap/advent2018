package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

func main() {
	in, err := ioutil.ReadFile("../inputs/15.txt")
	if err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	fmt.Println(string(in))
	g := GameFrom(in)
	fmt.Println(g.Units)
}

type Unit struct {
	X, Y   int
	Marker byte
	HP     int
	Power  int
}

func (u *Unit) attack(v *Unit) {
	v.HP -= u.Power
}

const (
	Empty = '.'
	Wall  = '#'
)

type Point struct {
	X, Y int
}

type Game struct {
	Map   *BytesImage
	Units []*Unit
	Round int
}

func GameFrom(in []byte) *Game {
	g := &Game{
		Map: BytesImageFrom(in),
	}

	g.Map.ForEach(func(x, y int, v byte) {
		if v == 'E' || v == 'G' {
			g.Units = append(g.Units, &Unit{
				Marker: v,
				X:      x,
				Y:      y,
				HP:     200,
				Power:  3,
			})
		}
	})

	return g
}

// Round runs a round and returns whether the round was completed or stopped partway
func (g *Game) RunRound() bool {
	unitTypes := make(map[byte]struct{})
	for _, u := range g.Units {
		unitTypes[u.Marker] = struct{}{}
	}
	if len(unitTypes) <= 1 {
		return false
	}

	sort.Slice(g.Units, func(i, j int) bool {
		return g.Units[i].Y < g.Units[j].Y || g.Units[i].X < g.Units[j].X
	})

	for _, u := range g.Units {
		if u.HP <= 0 {
			continue
		}

		inRange := g.unitsInRange(u)
		if len(inRange) == 0 {
			moved := g.moveUnit(u)
			if moved {
				inRange = g.unitsInRange(u)
			}
		}

		if len(inRange) > 0 {
			v := inRange[0]
			u.attack(v)
			if v.HP <= 0 {
				g.removeUnit(v)
			}
		}
	}

	g.Round++
	return true
}

func (g *Game) Outcome() int {
	var s int
	for _, u := range g.Units {
		s += u.HP
	}
	return g.Round * s
}

func (g *Game) removeUnit(u *Unit) {
	g.Map.Set(u.X, u.Y, Empty)
	for i := range g.Units {
		if g.Units[i] == u {
			g.Units = g.Units[:i+copy(g.Units[i:], g.Units[i+1:])]
			return
		}
	}
}

// moveUnit moves a unit closer to the nearest enemy, returning whether a move was possible or not
func (g *Game) moveUnit(u *Unit) bool {
	ui := g.Map.XY.ToIndex(u.X, u.Y)
	p := PathCalculatorFrom(g)
	var (
		minDist    int = PathMaxDist
		minUnitIdx int
	)
	for _, v := range g.Units {

		// ignore friends
		if u.Marker == v.Marker {
			continue
		}

		vi := g.Map.XY.ToIndex(v.X, v.Y)
		d := p.Dist(ui, vi)
		if d < minDist {
			minDist = d
			minUnitIdx = vi
		}
		//fmt.Printf("from %v to %v dist is %d\n", *u, *v, d)
	}

	if minDist == PathMaxDist {
		// no movement possible!
		fmt.Println("minDist is maxed!")
		return false
	}

	path := p.ShortestPath(ui, minUnitIdx)
	if len(path) == 0 {
		// no movement possible!
		fmt.Println("path is empty!")
		return false
	}

	nx, ny := g.Map.XY.FromIndex(path[0])
	//fmt.Printf("move from %d,%d to %d,%d\n", u.X, u.Y, nx, ny)
	g.Map.Set(u.X, u.Y, Empty)
	g.Map.Set(nx, ny, u.Marker)
	u.X, u.Y = nx, ny

	return true
}

func (g *Game) unitsInRange(u *Unit) []*Unit {
	nn := make(map[int]bool)
	g.Map.For4Neighbors(u.X, u.Y, func(x, y int, v byte) {
		if v != Empty && v != Wall {
			nn[g.Map.XY.ToIndex(x, y)] = true
		}
	})
	if len(nn) == 0 {
		return nil
	}

	var uu []*Unit
	for _, v := range g.Units {
		if nn[g.Map.XY.ToIndex(v.X, v.Y)] {
			uu = append(uu, v)
		}
	}
	return uu
}

const Newline = '\n'

type BytesImage struct {
	data          []byte
	width, height int
	XY            XYEncoder
}

func NewBytesImage(w, h int) *BytesImage {
	bi := &BytesImage{
		data:   make([]byte, (w+1)*h),
		width:  w + 1,
		height: h,
		XY:     XYEncoder{Width: w + 1, Height: h},
	}

	for i := w; i < len(bi.data); i += w + 1 {
		bi.data[i] = Newline
	}
	return bi
}

func BytesImageFrom(data []byte) *BytesImage {
	w := bytes.IndexRune(data, '\n') + 1
	h := bytes.Count(data, []byte{byte('\n')}) + 1
	return &BytesImage{
		data:   data[:],
		width:  w,
		height: h,
		XY:     XYEncoder{Width: w, Height: h},
	}
}

func (im *BytesImage) Set(x, y int, b byte) {
	i := im.XY.ToIndex(x, y)
	if i < 0 || i > len(im.data) {
		return
	}
	if im.data[i] == Newline {
		return
	}
	im.data[i] = b
}

func (im *BytesImage) ValueAt(x, y int) byte {
	i := im.XY.ToIndex(x, y)
	if i < 0 || i >= len(im.data) {
		return 0
	}
	b := im.data[i]
	if b == Newline {
		return 0
	}
	return b
}

func (im *BytesImage) ForEach(f func(x, y int, b byte)) {
	for i := range im.data {
		if im.data[i] == Newline {
			continue
		}

		x, y := im.XY.FromIndex(i)
		f(x, y, im.data[i])
	}
}

func (im *BytesImage) For4Neighbors(x, y int, f func(x, y int, b byte)) {
	nn := []struct{ x, y int }{{x, y - 1}, {x - 1, y}, {x + 1, y}, {x, y + 1}}
	for _, n := range nn {
		if n.x < 0 || n.y < 0 || n.x >= im.width || n.y >= im.width {
			continue
		}
		v := im.ValueAt(n.x, n.y)
		if v > 0 {
			f(n.x, n.y, v)
		}
	}
}

func (im *BytesImage) For8Neighbors(x, y int, f func(x, y int, b byte)) {
	for j := max(y-1, 0); j <= min(y+1, im.height-1); j++ {
		for i := max(x-1, 0); i <= min(x+1, im.width-1); i++ {
			if i == x && j == y {
				continue
			}
			v := im.ValueAt(i, j)
			if v > 0 {
				f(i, j, v)
			}
		}
	}
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

type XYEncoder struct {
	Width, Height int
}

func (e XYEncoder) ToIndex(x, y int) int {
	return y*e.Width + x
}

func (e XYEncoder) FromIndex(i int) (int, int) {
	y := i / e.Width
	x := i - (y * e.Width)
	return x, y
}

type Edge struct {
	U, V int
}

func (g *Game) Edges() []Edge {
	var ee []Edge
	g.Map.ForEach(func(ux, uy int, uv byte) {
		if uv == Wall {
			return
		}
		g.Map.For4Neighbors(ux, uy, func(vx, vy int, vv byte) {
			if vv == Wall || !(uv == Empty || vv == Empty) {
				return
			}
			ee = append(ee, Edge{
				U: g.Map.XY.ToIndex(ux, uy),
				V: g.Map.XY.ToIndex(vx, vy),
			})
		})
	})
	return ee
}

const PathMaxDist = math.MaxInt32
