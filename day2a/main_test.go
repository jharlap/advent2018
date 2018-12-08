package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWordToCounts(t *testing.T) {
	cc := []struct {
		word   string
		counts map[rune]int
	}{
		{"abcde", map[rune]int{'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1}},
		{"abbde", map[rune]int{'a': 1, 'b': 2, 'd': 1, 'e': 1}},
		{"abbbe", map[rune]int{'a': 1, 'b': 3, 'e': 1}},
		{"abbbcdee", map[rune]int{'a': 1, 'b': 3, 'c': 1, 'd': 1, 'e': 2}},
	}

	for _, c := range cc {
		t.Run(c.word, func(t *testing.T) {
			m := wordToCounts(c.word)
			if !reflect.DeepEqual(c.counts, m) {
				t.Errorf("got %v expected %v", m, c.counts)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	cc := []struct {
		ww []string
		c  int
	}{
		{[]string{"abcdef", "bababc", "abbcde", "abcccd", "aabcdd", "abcdee", "ababab"}, 12},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			v := checksum(c.ww)
			if v != c.c {
				t.Errorf("got %v expected %v", v, c.c)
			}
		})
	}
}
