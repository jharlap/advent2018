package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestScore(t *testing.T) {
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
	exp := []int{0, 16, 8, 17, 4, 18, 19, 2, 24, 20, 25, 10, 21, 5, 22, 11, 1, 12, 6, 13, 3, 14, 7, 15}
	expPos := []int{0, 1, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15, 1, 3, 5, 7, 9, 11, 13, 6, 8, 10}
	g := newGame(5, 26)
	for i := 0; i < 26; i++ {
		preCur := g.cur
		p := g.takeTurn()
		if p != expPos[i] {
			t.Errorf("pos %d (%d) got %d expected %d", i, preCur, p, expPos[i])
		}
	}
	if !reflect.DeepEqual(exp, g.ring) {
		t.Errorf("got %v expected %v", g.ring, exp)
	}
}
