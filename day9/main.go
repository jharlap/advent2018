package main

import (
	"fmt"
)

func main() {
	s := score(400, 71864)
	fmt.Println("winning score:", s)
	s = score(400, 7186400)
	fmt.Println("winning score 100x:", s)
}

func score(players, lastMarbleScore int) int {
	marbles := lastMarbleScore + 1
	g := newGame(players, marbles)
	for i := 0; i < marbles; i++ {
		g.takeTurn()
	}
	return g.topScore()
}

type game struct {
	ring       []int
	scores     map[int]int
	cur        int
	nextMarble int
	players    int
	nextPlayer int
}

func newGame(players, marbles int) game {
	return game{
		players: players,
		scores:  make(map[int]int),
		ring:    make([]int, 0, marbles),
	}
}

func (g *game) topScore() int {
	var s int
	for _, v := range g.scores {
		if v > s {
			s = v
		}
	}
	return s
}

func (g *game) takeTurn() int {
	p := g.nextPlayer
	m := g.nextMarble
	g.nextMarble++
	g.nextPlayer = (g.nextPlayer + 1) % g.players

	if m > 0 && m%23 == 0 {
		r := normalizeIndex(g.cur-7, len(g.ring))
		g.scores[p] += m + g.ring[r]
		g.remove(r)
		g.cur = r
		return g.cur
	}

	i := g.normalizeInsertionIndex(g.cur + 2)
	g.insert(m, i)
	g.cur = i
	return g.cur
}

func normalizeIndex(i, curLen int) int {
	for i < curLen {
		i += curLen
	}
	if i >= curLen {
		i = i % curLen
	}
	return i
}

func (g *game) normalizeInsertionIndex(i int) int {
	cl := len(g.ring)
	if cl <= 1 {
		return cl
	}

	i = normalizeIndex(i, cl)
	if i == 0 {
		i = cl
	}
	return i
}

func (g *game) remove(i int) int {
	m := g.ring[i]

	a := g.ring
	a = a[:i+copy(a[i:], a[i+1:])]
	g.ring = a

	//g.ring = append(g.ring[:i], g.ring[i+1:]...)
	return m
}

func (g *game) insert(m, i int) {
	if len(g.ring) == i {
		g.ring = append(g.ring, m)
		return
	}
	g.ring = append(g.ring, 0)
	copy(g.ring[i+1:], g.ring[i:])
	g.ring[i] = m
}
