package main

import (
	"fmt"
	"testing"
)

func TestScore(t *testing.T) {
	t.Skip()
	cc := []struct {
		players int
		marbles int
		score   int
	}{
		{9, 25, 32},
		{10, 1618, 8317},
		{13, 7999, 146373},
		{17, 1104, 2764},
		{21, 6111, 54718},
		{30, 5807, 37305},
	}

	for _, c := range cc {
		t.Run(fmt.Sprintf("p %d m %d", c.players, c.marbles), func(t *testing.T) {
			s := score(c.players, c.marbles)
			if s != c.score {
				t.Errorf("got %d expected %d", s, c.score)
			}
		})
	}
}

func TestTakeTurn(t *testing.T) {
	expPos := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 19, 24, 25}
	g := newGame(5, 26)
	for i := 0; i < 26; i++ {
		p := g.takeTurn()
		if p != expPos[i] {
			t.Errorf("pos %d got %d expected %d", i, p, expPos[i])
		}
	}
}
