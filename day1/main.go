package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	ff, err := readFreqs("../inputs/1.txt")
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	f := findFreq(ff)
	fmt.Println("Freq:", f)
}

func readFreqs(filename string) ([]int, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input: %v", err)
	}
	defer r.Close()

	var ff []int
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		if i, err := strconv.Atoi(t); err != nil {
			return nil, fmt.Errorf("error atoi %v: %v", t, err)
		} else {
			ff = append(ff, i)
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}
	return ff, nil
}
func findFreq(ff []int) int {
	seen := make(map[int]bool)
	seen[0] = true
	var f int
	for {
		for _, v := range ff {
			f += v
			if seen[f] {
				return f
			}
			seen[f] = true
		}
	}
}
