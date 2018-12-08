package main

import (
	"fmt"
	"testing"
)

func TestFindFreq(t *testing.T) {
	cc := []struct {
		ff []int
		f  int
	}{
		{[]int{+1, -1}, 0},
		{[]int{+3, +3, +4, -2, -4}, 10},
		{[]int{-6, +3, +8, +5, -6}, 5},
		{[]int{+7, +7, -2, -7, -4}, 14},
	}
	for _, c := range cc {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			f := findFreq(c.ff)
			if c.f != f {
				t.Errorf("expected %v got %v", c.f, f)
			}
		})
	}
}
