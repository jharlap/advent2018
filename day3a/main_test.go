package main

import (
	"fmt"
	"testing"
)

func TestIntersection(t *testing.T) {
	cc := []struct {
		a, b       Rect
		intersects bool
		exp        Rect
	}{
		{Rect{1, 2, 4, 2}, Rect{0, 3, 2, 2}, true, Rect{1, 3, 1, 1}},
		{Rect{3, 1, 4, 2}, Rect{0, 3, 2, 2}, false, Rect{}},
		{Rect{1, 2, 4, 2}, Rect{3, 1, 4, 2}, true, Rect{3, 2, 2, 1}},
		{Rect{3, 1, 4, 2}, Rect{1, 2, 4, 2}, true, Rect{3, 2, 2, 1}},
		{Rect{0, 0, 9, 9}, Rect{2, 2, 3, 2}, true, Rect{2, 2, 3, 2}},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			ok, r := intersection(c.a, c.b)
			if ok != c.intersects {
				t.Errorf("intersects: got %v expected %v", ok, c.intersects)
			}
			if r != c.exp {
				t.Errorf("intersection: got %v expected %v", r, c.exp)
			}
		})
	}
}
func TestIntersects(t *testing.T) {
	cc := []struct {
		a, b       Rect
		intersects bool
	}{
		{Rect{1, 3, 4, 4}, Rect{3, 1, 4, 4}, true},
		{Rect{1, 3, 4, 4}, Rect{5, 5, 2, 2}, false},
		{Rect{3, 1, 4, 4}, Rect{1, 3, 4, 4}, true},
		{Rect{0, 0, 9, 9}, Rect{2, 2, 2, 2}, true},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			r := intersects(c.a, c.b)
			if r != c.intersects {
				t.Errorf("got %v expected %v", r, c.intersects)
			}
		})
	}
}
