package main

import (
	"fmt"
	"testing"
)

func TestPowerLevel(t *testing.T) {
	cc := []struct {
		x, y, grid int
		exp        int
	}{
		{3, 5, 8, 4},
		{122, 79, 57, -5},
		{217, 196, 39, 0},
		{101, 153, 71, 4},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			p := powerLevel(c.x, c.y, c.grid)
			if p != c.exp {
				t.Errorf("got %d expected %d", p, c.exp)
			}
		})
	}
}

func TestAreaPowerLevel(t *testing.T) {
	cc := []struct {
		x, y, size, grid int
		exp              int
	}{
		{33, 45, 3, 18, 29},
		{21, 61, 3, 42, 30},
		{90, 269, 16, 18, 113},
		{232, 251, 12, 42, 119},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			p := areaPowerLevel(c.x, c.y, c.size, c.grid)
			if p != c.exp {
				t.Errorf("got %d expected %d", p, c.exp)
			}
		})
	}
}
