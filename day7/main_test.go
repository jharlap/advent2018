package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSolveDeps(t *testing.T) {
	adj := map[string][]string{
		"A": {"C"},
		"F": {"C"},
		"B": {"A"},
		"D": {"A"},
		"E": {"B", "D", "F"},
	}
	exp := "CABDFE"
	g := graph{adj: adj, nodes: nodes(adj)}
	steps, err := solveDeps(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	act := strings.Join(steps, "")
	if exp != act {
		t.Errorf("got: %s expected: %s", act, exp)
	}
}

func TestSolveDepsTimed(t *testing.T) {
	adj := map[string][]string{
		"A": {"C"},
		"F": {"C"},
		"B": {"A"},
		"D": {"A"},
		"E": {"B", "D", "F"},
	}
	exp := 15
	g := graph{adj: adj, nodes: nodes(adj)}
	act, err := solveDepsTimed(g, 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp != act {
		t.Errorf("got: %d expected: %d", act, exp)
	}
}

func TestTimeFor(t *testing.T) {
	cc := []struct {
		s string
		i int
	}{
		{"A", 1},
		{"Z", 26},
	}
	for _, c := range cc {
		t.Run(c.s, func(t *testing.T) {
			a := timeFor(c.s)
			if a != c.i {
				t.Errorf("got: %d expected: %d", a, c.i)
			}
		})
	}
}
func TestDoable(t *testing.T) {
	cc := []struct {
		workers int
		adj     map[string][]string
		done    map[string]bool
		doing   map[string]int
		exp     string
	}{
		{1, map[string][]string{"A": {"B", "C"}}, nil, nil, "B"},
		{2, map[string][]string{"A": {"B", "C"}}, nil, nil, "BC"},
		{2, map[string][]string{"A": {"B", "C"}}, nil, map[string]int{"C": 1}, "B"},
		{1, map[string][]string{"A": {"B", "C"}, "B": {"C"}}, nil, nil, "C"},
		{2, map[string][]string{"A": {"B", "C"}, "B": {"C"}}, nil, nil, "C"},
		{1, map[string][]string{"A": {"B", "C"}, "B": {"C"}}, map[string]bool{"C": true}, nil, "B"},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			g := graph{adj: c.adj, nodes: nodes(c.adj)}
			r := doable(g, c.done, c.doing, c.workers)
			act := strings.Join(r, "")
			if act != c.exp {
				t.Errorf("got: %s expected: %s", act, c.exp)
			}
		})
	}
}

func TestNodes(t *testing.T) {

}
