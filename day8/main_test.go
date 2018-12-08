package main

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	cc := []struct {
		in          []int
		nodes       int
		metadataSum int
	}{
		{[]int{2, 3, 0, 3, 10, 11, 12, 1, 1, 0, 1, 99, 2, 1, 1, 2}, 4, 138},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			root, n, err := tree(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != len(c.in) {
				t.Errorf("error: tree did not read all the input")
			}

			nn := countNodes(root)
			if nn != c.nodes {
				t.Errorf("node count got %d expected %d", nn, c.nodes)
			}

			ms := metadataSum(root)
			if ms != c.metadataSum {
				t.Errorf("metadata sum got %d expected %d", ms, c.metadataSum)
			}
		})
	}
}

func countNodes(root Node) int {
	n := 1
	for _, c := range root.Children {
		n += countNodes(c)
	}
	return n
}

func TestTreeValue(t *testing.T) {
	cc := []struct {
		in    []int
		value int
	}{
		{[]int{2, 3, 0, 3, 10, 11, 12, 1, 1, 0, 1, 99, 2, 1, 1, 2}, 66},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			root, n, err := tree(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != len(c.in) {
				t.Errorf("error: tree did not read all the input")
			}

			v := treeValue(root)
			if v != c.value {
				t.Errorf("got %d expected %d", v, c.value)
			}
		})
	}
}
