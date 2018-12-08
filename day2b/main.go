package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	ww, err := readInput("../inputs/2.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for i, a := range ww {
		for j, b := range ww {
			if i == j {
				continue
			}

			ar := []rune(a)
			br := []rune(b)
			m := matches(ar, br)
			if countMismatches(m) == 1 {
				w := removeMismatch(ar, m)
				fmt.Printf("a: %s\nb: %s\nw: %s\n", a, b, w)
				return
			}
		}
	}
	fmt.Println("No word found")
}

func removeMismatch(w []rune, m []bool) string {
	var r []rune
	for i := 0; i < len(w); i++ {
		if m[i] {
			r = append(r, w[i])
		}
	}
	return string(r)
}

func matches(a, b []rune) []bool {
	if len(a) != len(b) {
		return nil
	}

	var r []bool
	for i := 0; i < len(a); i++ {
		r = append(r, a[i] == b[i])
	}
	return r
}

func countMismatches(m []bool) int {
	var r int
	for _, b := range m {
		if !b {
			r++
		}
	}
	return r
}

func readInput(filename string) ([]string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input: %v", err)
	}
	defer r.Close()

	var l []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		l = append(l, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}
	return l, nil
}
