package main

import (
	"strings"
	"testing"
)

func TestValueAtMinute(t *testing.T) {
	in := `.#.#...|#.
.....#|##|
.|..|...#.
..|#.....#
#.#|||#|#|
...#.||...
.|....|...
||...#|.#|
|.||||..|.
...#.|..|.`

	img := BytesImageFrom([]byte(in))

	exp := 1147
	v := valueOf(stateAtMinute(img, 10))
	if v != exp {
		t.Errorf("got %d expected %d", v, exp)
	}
}

func TestAdvanceMinute(t *testing.T) {
	in := `.#.#...|#.
.....#|##|
.|..|...#.
..|#.....#
#.#|||#|#|
...#.||...
.|....|...
||...#|.#|
|.||||..|.
...#.|..|.`

	exp := `.......##.
......|###
.|..|...#.
..|#||...#
..##||.|#|
...#||||..
||...|||..
|||||.||.|
||||||||||
....||..|.`

	img := BytesImageFrom([]byte(in))
	v := advanceMinute(img)
	if strings.TrimSpace(v.String()) != exp {
		t.Errorf("got %s expected %s", v, exp)
	}
}
