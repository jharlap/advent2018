package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in, err := readInput("../inputs/8.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	root, _, err := tree(in)
	if err != nil {
		fmt.Println("Error building tree:", err)
		os.Exit(1)
	}

	fmt.Println("metadata sum:", metadataSum(root))

	fmt.Println("value:", treeValue(root))
}

type Node struct {
	Children []Node
	Metadata []int
}

// tree builds a tree from the input, returning the root node and number of input values read
func tree(in []int) (Node, int, error) {
	if len(in) < 2 {
		return Node{}, 0, fmt.Errorf("error: invalid input to tree (too few inputs): %v", in)
	}

	var root Node
	nr := 2
	nc := in[0]
	nm := in[1]

	for i := 0; i < nc; i++ {
		c, n, err := tree(in[nr:])
		if err != nil {
			return Node{}, 0, fmt.Errorf("error build child: %v", err)
		}
		nr += n
		root.Children = append(root.Children, c)
	}

	root.Metadata = in[nr : nr+nm]
	nr += nm
	return root, nr, nil
}

func metadataSum(root Node) int {
	var m int
	for _, c := range root.Children {
		m += metadataSum(c)
	}
	for _, v := range root.Metadata {
		m += v
	}
	return m
}

func treeValue(root Node) int {
	if len(root.Children) == 0 {
		var v int
		for _, i := range root.Metadata {
			v += i
		}
		return v
	}

	var value int
	nChildren := len(root.Children)
	for _, n := range root.Metadata {
		if n > nChildren {
			continue // invalid indices value 0
		}
		n-- // 1-based indexing

		value += treeValue(root.Children[n])
	}
	return value
}

func readInput(filename string) ([]int, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input: %v", err)
	}
	defer r.Close()

	var l []int
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	for s.Scan() {
		w := s.Text()
		i, err := strconv.Atoi(w)
		if err != nil {
			return nil, fmt.Errorf("error parsing word '%s': %v", w, err)
		}
		l = append(l, i)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}
	return l, nil
}
