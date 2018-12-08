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
	fmt.Println("checksum:", checksum(ww))
}

func checksum(ww []string) int {
	var two, three int
	for _, w := range ww {
		c := wordToCounts(w)
		if hasExactlyN(c, 2) {
			two++
		}
		if hasExactlyN(c, 3) {
			three++
		}
	}
	return two * three
}

func wordToCounts(w string) map[rune]int {
	m := make(map[rune]int)
	rr := []rune(w)
	for _, r := range rr {
		m[r]++
	}

	return m
}

func hasExactlyN(m map[rune]int, n int) bool {
	for _, v := range m {
		if v == n {
			return true
		}
	}
	return false
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
