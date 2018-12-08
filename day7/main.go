package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	g, err := readInput("../inputs/7.txt")
	if err != nil {
		fmt.Println("error reading input:", err)
		os.Exit(1)
	}

	steps, err := solveDeps(g)
	if err != nil {
		fmt.Println("error solving deps:", err)
		os.Exit(1)
	}

	fmt.Println("Steps:", strings.Join(steps, ""))

	dur, err := solveDepsTimed(g, 15, 60)
	if err != nil {
		fmt.Println("error solving deps timed:", err)
		os.Exit(1)
	}

	fmt.Println("Duration:", dur)
}

type graph struct {
	nodes []string
	adj   map[string][]string
}

func solveDeps(g graph) ([]string, error) {
	var steps []string
	done := make(map[string]bool)
	for len(done) < len(g.nodes) {
		todo := doable(g, done, nil, 1)
		if len(todo) == 0 {
			return nil, fmt.Errorf("error: impossible after %v", steps)
		}
		steps = append(steps, todo...)
		for _, s := range todo {
			done[s] = true
		}
	}
	return steps, nil
}

func solveDepsTimed(g graph, workers int, delay int) (int, error) {
	done := make(map[string]bool)
	inProgress := make(map[string]int)
	var t int
	for len(done) < len(g.nodes) {
		for k := range inProgress {
			inProgress[k]--
			if inProgress[k] <= 0 {
				done[k] = true
				delete(inProgress, k)
			}
		}
		todo := doable(g, done, inProgress, workers-len(inProgress))
		for _, s := range todo {
			inProgress[s] = timeFor(s) + delay
		}
		//fmt.Println("t:", t, "done:", keys(done), "doing:", inProgress)
		t++
	}
	return t - 1, nil
}

func keys(m map[string]bool) []string {
	var r []string
	for k := range m {
		r = append(r, k)
	}
	return r
}

func timeFor(s string) int {
	return int(s[0]) - 'A' + 1
}

func doable(g graph, done map[string]bool, doing map[string]int, workers int) []string {
	var todo []string
	for _, n := range g.nodes {
		if done[n] {
			continue
		}
		ok := true
		for _, d := range g.adj[n] {
			ok = ok && done[d]
		}
		if ok && doing[n] == 0 {
			todo = append(todo, n)
		}
	}
	sort.Strings(todo)
	return todo[:min(workers, len(todo))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func nodes(adj map[string][]string) []string {
	m := make(map[string]bool)
	for k, v := range adj {
		m[k] = true
		for _, s := range v {
			m[s] = true
		}
	}

	var ss []string
	for k := range m {
		ss = append(ss, k)
	}
	return ss
}

func readInput(filename string) (graph, error) {
	f, err := os.Open(filename)
	if err != nil {
		return graph{}, fmt.Errorf("error opening file '%s': %v", filename, err)
	}
	defer f.Close()

	vv := make(map[string][]string)
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		var n, d string
		_, err := fmt.Sscanf(l, "Step %s must be finished before step %s can begin.", &d, &n)
		if err != nil {
			return graph{}, fmt.Errorf("error parsing line '%s': %v", l, err)
		}

		if !contains(vv[n], d) {
			vv[n] = append(vv[n], d)
		}
	}

	nn := nodes(vv)
	return graph{nodes: nn, adj: vv}, nil
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
