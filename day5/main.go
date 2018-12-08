package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	bb, err := ioutil.ReadFile("../inputs/5.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	s := string(bb)
	fmt.Println("len:", len(collapse(s)))
	b := len(s)
	for i := 'A'; i <= 'Z'; i++ {
		f := filter(s, i)
		l := len(collapse(f))
		if l < b {
			b = l
		}
	}
	fmt.Println("best len:", b)
}

func filter(s string, c rune) string {
	const Aa = 32
	return strings.Map(func(r rune) rune {
		if r == c || r == c+Aa {
			return -1
		}
		return r
	}, s)
}

func collapse(s string) string {
	const Aa = 32
	in := []byte(s)
	removed := -1
	for removed != 0 {
		removed = 0
		if len(in) < 2 {
			return string(in)
		}
		for i := len(in) - 2; i >= 0; i-- {
			if len(in) < i+2 {
				return string(in)
			}
			a := int(in[i])
			b := int(in[i+1])
			d := a - b
			if d == Aa || d == -Aa {
				removed++
				in = append(in[:i], in[i+2:]...)
			}
		}
	}
	return string(in)
}
