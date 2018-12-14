package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := 165061
	fmt.Println("10 after", input, "is", intsToString(TenAfterN(input)))

	t := FindTickForRecipes([]int{1, 6, 5, 0, 6, 1})
	fmt.Println("found input at tick:", t)
}

type Scoreboard struct {
	Scores     []int
	Elf1, Elf2 int
}

func StartingScoreboard() *Scoreboard {
	return &Scoreboard{
		Scores: []int{3, 7},
		Elf1:   0,
		Elf2:   1,
	}
}

func (s *Scoreboard) Tick() {
	s.Scores = append(s.Scores, combineScores(s.Scores[s.Elf1], s.Scores[s.Elf2])...)
	s.Elf1 = (s.Elf1 + 1 + s.Scores[s.Elf1]) % len(s.Scores)
	s.Elf2 = (s.Elf2 + 1 + s.Scores[s.Elf2]) % len(s.Scores)
}

func TenAfterN(n int) []int {
	s := StartingScoreboard()
	for len(s.Scores) < n+10 {
		s.Tick()
	}
	return s.Scores[n : n+10]
}

func FindTickForRecipes(seq []int) int {
	seqLen := len(seq)
	s := StartingScoreboard()
	var lastLen int
	for {
		s.Tick()
		newLen := len(s.Scores)
		if newLen < seqLen {
			continue
		}
		for i := max(lastLen-seqLen, 0); i+seqLen < newLen; i++ {
			cur := s.Scores[i : i+seqLen]
			//fmt.Printf("i: %d lastLen: %d newLen: %d seq: %s score: %s matches: %v\n", i, lastLen, newLen, intsToString(seq), intsToString(cur), matches(seq, cur))
			if matches(seq, cur) {
				return i
			}
		}
		lastLen = newLen
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func matches(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func combineScores(a, b int) []int {
	v := a + b
	if v > 9 {
		return []int{v / 10, v % 10}
	}
	return []int{v}
}

func intsToString(ii []int) string {
	var ss []string
	for _, i := range ii {
		ss = append(ss, strconv.Itoa(i))
	}
	return strings.Join(ss, "")
}
