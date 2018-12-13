package main

import (
	"fmt"
	"testing"
)

func TestTurn(t *testing.T) {
	cc := []struct {
		dir      XY
		nextTurn Turn
		exp      XY
	}{
		{DirUp, TurnLeft, DirLeft},
		{DirUp, TurnRight, DirRight},
		{DirUp, TurnStraight, DirUp},
		{DirDown, TurnLeft, DirRight},
		{DirDown, TurnRight, DirLeft},
		{DirDown, TurnStraight, DirDown},
		{DirRight, TurnLeft, DirUp},
		{DirRight, TurnRight, DirDown},
		{DirRight, TurnStraight, DirRight},
		{DirLeft, TurnLeft, DirDown},
		{DirLeft, TurnRight, DirUp},
		{DirLeft, TurnStraight, DirLeft},
	}

	for _, c := range cc {
		t.Run(fmt.Sprintf("D_%d_%d_T%d", c.dir.X, c.dir.Y, c.nextTurn), func(t *testing.T) {
			v := NewCart(XY{}, c.dir)
			v.nextTurn = c.nextTurn

			v.Turn()

			if v.Dir.X != c.exp.X || v.Dir.Y != c.exp.Y {
				t.Errorf("got %v expected %v", v.Dir, c.exp)
			}
		})
	}
}

func TestReadMap(t *testing.T) {
	in := `/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   
`
	m := MapFrom([]byte(in))
	cc := []struct {
		y, x int
		exp  Track
	}{
		{0, 0, TrackCurveBottomLeftTopRight},
		{0, 1, TrackHorizontal},
		{0, 2, TrackHorizontal},
		{0, 3, TrackHorizontal},
		{0, 4, TrackCurveTopLeftBottomRight},
		{1, 0, TrackVertical},
		{1, 1, TrackEmpty},
		{1, 2, TrackEmpty},
		{1, 3, TrackEmpty},
		{1, 4, TrackVertical},
		{2, 0, TrackVertical},
		{2, 1, TrackEmpty},
		{2, 2, TrackCurveBottomLeftTopRight},
		{2, 3, TrackHorizontal},
		{2, 4, TrackIntersection},
		{2, 5, TrackHorizontal},
		{2, 6, TrackHorizontal},
		{2, 7, TrackIntersection},
		{2, 8, TrackHorizontal},
		{2, 9, TrackCurveTopLeftBottomRight},
		{2, 10, TrackEmpty},
		{2, 11, TrackEmpty},
		{2, 12, TrackVertical},
	}
	for _, c := range cc {
		act := m.TrackAt(XY{c.x, c.y})
		if act != c.exp {
			t.Errorf("%d,%d got %v expected %v", c.x, c.y, act, c.exp)
		}
	}

	carts := m.Carts()
	if len(carts) != 2 {
		t.Errorf("carts got %d expected 2", len(carts))
	}
}

func TestFirstCrash(t *testing.T) {
	in := `/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   
`
	expected := XY{7, 3}
	m := MapFrom([]byte(in))
	p := m.FirstCrash(nil)
	if !expected.Equals(p) {
		t.Errorf("got %v expected %v", p, expected)
	}
}

func TestLastCart(t *testing.T) {
	in := `
/>-<\  
|   |  
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/
`
	expected := XY{6, 4}
	m := MapFrom([]byte(in))
	p := m.LastCart()
	if !expected.Equals(*p) {
		t.Errorf("got %v expected %v", p, expected)
	}

}

func debugTick(i int, m *Map) bool {
	fmt.Println(i, "--------------")
	fmt.Println(m)
	return true
}
