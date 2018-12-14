package main

import (
	"fmt"
	"testing"
)

func TestTenAfterN(t *testing.T) {
	cc := []struct {
		n   int
		exp string
	}{
		{9, "5158916779"},
		{5, "0124515891"},
		{18, "9251071085"},
		{2018, "5941429882"},
	}

	for _, c := range cc {
		t.Run(fmt.Sprintf("ticks %d", c.n), func(t *testing.T) {
			r := TenAfterN(c.n)
			rs := intsToString(r)
			if rs != c.exp {
				t.Errorf("got %s expected %s", rs, c.exp)
			}
		})
	}
}

func TestFindTickForRecipes(t *testing.T) {
	cc := []struct {
		in  []int
		exp int
	}{
		{[]int{5, 1, 5, 8, 9}, 9},
		{[]int{0, 1, 2, 4, 5}, 5},
		{[]int{9, 2, 5, 1, 0}, 18},
		{[]int{9, 2, 5, 1, 0, 7}, 18},
		{[]int{5, 9, 4, 1, 4}, 2018},
	}

	for _, c := range cc {
		t.Run(fmt.Sprintf("input %d", c.in), func(t *testing.T) {
			r := FindTickForRecipes(c.in)
			if r != c.exp {
				t.Errorf("got %d expected %d", r, c.exp)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	cc := []struct {
		a, b []int
		exp  bool
	}{
		{[]int{1}, []int{1, 2}, false},
		{[]int{1, 2}, []int{1, 2}, true},
		{[]int{1, 2}, []int{1, 3}, false},
		{[]int{3, 2}, []int{1, 2}, false},
	}
	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			r := matches(c.a, c.b)
			if r != c.exp {
				t.Errorf("got %v expected %v", r, c.exp)
			}
		})
	}
}
