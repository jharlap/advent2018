package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	ll, err := readInput("../inputs/4.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	sort.Strings(ll)

	gi := periods(ll)
	fmt.Println("guards:", len(gi))

	{
		var max GuardInfo
		for _, g := range gi {
			if g.Duration > max.Duration {
				max = g
			}
		}

		fmt.Println(max)

		m, _ := bestMinute(max.Periods)
		fmt.Println("best minute:", m)
		fmt.Println("result:", max.ID*m)
	}

	{
		var bm, bc, bg int
		for _, g := range gi {
			m, c := bestMinute(g.Periods)
			if c > bc {
				bm = m
				bc = c
				bg = g.ID
			}
		}
		fmt.Println("g:", bg, "m:", bm, "c:", bc, "r:", bg*bm)
	}
}

func bestMinute(pp []Period) (int, int) {
	m := make([]int, 60)
	var bm, bc int
	for _, p := range pp {
		for i := p.SleepMin; i < p.WakeMin; i++ {
			m[i]++
			if m[i] > bc {
				bc = m[i]
				bm = i
			}
		}
	}

	return bm, bc
}

type GuardInfo struct {
	ID       int
	Duration int
	Periods  []Period
}

type Period struct {
	SleepMin int
	WakeMin  int
}

func periods(ll []string) map[int]GuardInfo {
	beginShift := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}\] Guard #(\d+) begins shift`)
	sleep := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:(\d{2})\] falls asleep`)
	wake := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:(\d{2})\] wakes up`)

	var (
		g    int
		s, w int
		r    = make(map[int]GuardInfo)
	)
	for _, l := range ll {
		if len(l) == 0 {
			continue
		}

		if mm := beginShift.FindStringSubmatch(l); mm != nil {
			g = atoi(mm[1])
		} else if mm := sleep.FindStringSubmatch(l); mm != nil {
			s = atoi(mm[1])
		} else if mm := wake.FindStringSubmatch(l); mm != nil {
			w = atoi(mm[1])

			gi := r[g]
			gi.ID = g
			gi.Duration += w - s
			gi.Periods = append(gi.Periods, Period{SleepMin: s, WakeMin: w})
			r[g] = gi
		} else {
			fmt.Println("No match!", l)
		}
	}
	return r
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func readInput(filename string) ([]string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input: %v", err)
	}
	defer r.Close()

	var rr []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		rr = append(rr, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}
	return rr, nil
}
