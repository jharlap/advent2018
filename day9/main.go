package main

import (
	"container/list"
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
	ring       *list.List
	cur        *list.Element
	scores     map[int]int
	nextMarble int
	players    int
	nextPlayer int
}

func newGame(players, marbles int) game {
	return game{
		players: players,
		scores:  make(map[int]int),
		ring:    list.New(),
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

	if g.ring.Len() == 0 {
		g.cur = g.ring.PushFront(0)
		return g.cur.Value.(int)
	} else if g.ring.Len() == 1 {
		g.cur = g.ring.PushBack(1)
		return g.cur.Value.(int)
	} else if g.ring.Len() == 2 {
		g.cur = g.ring.InsertBefore(2, g.cur)
		return g.cur.Value.(int)
	}

	if m%23 == 0 {
		g.scores[p] += m
		for i := 0; i < 7; i++ {
			g.movePrev()
		}
		newCur := g.cur.Next()
		if newCur == nil {
			newCur = g.ring.Front()
		}
		g.scores[p] += g.ring.Remove(g.cur).(int)
		g.cur = newCur
		return g.cur.Value.(int)
	}

	g.moveNext()
	g.cur = g.ring.InsertAfter(m, g.cur)
	return g.cur.Value.(int)
}

func (g *game) moveNext() {
	if g.cur.Next() == nil {
		g.cur = g.ring.Front()
	} else {
		g.cur = g.cur.Next()
	}
}

func (g *game) movePrev() {
	if g.cur.Prev() == nil {
		g.cur = g.ring.Back()
	} else {
		g.cur = g.cur.Prev()
	}
}
