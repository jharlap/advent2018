package main

import (
	"testing"
)

func TestCollapse(t *testing.T) {
	cc := []struct {
		in, ex string
	}{
		{"a", "a"},
		{"aA", ""},
		{"Aa", ""},
		{"bAa", "b"},
		{"dabAcCaCBAcCcaDA", "dabCBAcaDA"},
	}
	for _, c := range cc {
		t.Run(c.in, func(t *testing.T) {
			r := collapse(c.in)
			if r != c.ex {
				t.Errorf("got %s expected %s", r, c.ex)
			}
		})
	}
}
