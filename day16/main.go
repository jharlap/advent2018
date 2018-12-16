package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	ll, err := readLines("../inputs/16_samples.txt")
	if err != nil {
		fmt.Println("error reading samples:", err)
		os.Exit(1)
	}

	var n int
	od := newOpDecoder()
	for i := 0; i < len(ll); i += 4 {
		s, err := parseSample(ll[i : i+3])
		if err != nil {
			fmt.Println("Error parsing sample:", ll[i:i+3], "err:", err)
			os.Exit(1)
		}

		cands := candidateOps(s)
		if len(cands) >= 3 {
			n++
		}

		if s.instr[0] == 90 {
			for _, c := range cands {
				fmt.Println("9 could be", c.i)
			}
		}
		od.addCandidates(s.instr[0], cands)
	}

	fmt.Println("Samples behaving like 3+ ops:", n)

	fmt.Println("od:", od)

	cpu := CPUFrom(od)

	// figured out on paper
	/*
		cpu := CPU{
			Reg: make([]int, 4),
			instr: []opNamed{
				allOps[2],
				allOps[14],
				allOps[8],
				allOps[15],
				allOps[12],
				allOps[3],
				allOps[6],
				allOps[5],
				allOps[0],
				allOps[4],
				allOps[13],
				allOps[10],
				allOps[1],
				allOps[11],
				allOps[9],
				allOps[7],
			},
		}
	*/

	ll, err = readLines("../inputs/16_prog.txt")
	if err != nil {
		fmt.Println("error reading samples:", err)
		os.Exit(1)
	}

	for _, l := range ll {
		if len(l) == 0 {
			continue
		}
		in, err := parseInstruction(l)
		if err != nil {
			fmt.Println("error parsing instruction:", l, "err:", err)
			os.Exit(1)
		}

		cpu.Exec(in[0], in[1], in[2], in[3])
	}

	fmt.Println("R0:", cpu.Reg[0])
}

func readLines(filename string) ([]string, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return strings.Split(string(bb), "\n"), nil
}

type CPU struct {
	Reg   []int
	instr []opNamed
}

func CPUFrom(od *opDecoder) CPU {
	return CPU{
		Reg:   make([]int, 4),
		instr: od.solve(),
	}
}

func (cpu *CPU) Exec(o, a, b, c int) {
	cpu.instr[o].op(a, b, c, cpu.Reg)
}

type op func(a, b, c int, in []int)

// addr (add register) stores into register C the result of adding register A and register B.
func addr(a, b, c int, o []int) {
	o[c] = o[a] + o[b]
}

// addi (add immediate) stores into register C the result of adding register A and value B.
func addi(a, b, c int, o []int) {
	o[c] = o[a] + b
}

// mulr (multiply register) stores into register C the result of multiplying register A and register B.
func mulr(a, b, c int, o []int) {
	o[c] = o[a] * o[b]
}

// muli (multiply immediate) stores into register C the result of multiplying register A and value B.
func muli(a, b, c int, o []int) {
	o[c] = o[a] * b
}

// banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
func banr(a, b, c int, o []int) {
	o[c] = o[a] & o[b]
}

// bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.
func bani(a, b, c int, o []int) {
	o[c] = o[a] & b
}

// borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
func borr(a, b, c int, o []int) {
	o[c] = o[a] | o[b]
}

// bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.
func bori(a, b, c int, o []int) {
	o[c] = o[a] | b
}

// setr (set register) copies the contents of register A into register C. (Input B is ignored.)
func setr(a, _, c int, o []int) {
	o[c] = o[a]
}

// seti (set immediate) stores value A into register C. (Input B is ignored.)
func seti(a, _, c int, o []int) {
	o[c] = a
}

// gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
func gtir(a, b, c int, o []int) {
	if a > o[b] {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

// gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
func gtri(a, b, c int, o []int) {
	if o[a] > b {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

// gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.
func gtrr(a, b, c int, o []int) {
	if o[a] > o[b] {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

// eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
func eqir(a, b, c int, o []int) {
	if a == o[b] {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

// eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
func eqri(a, b, c int, o []int) {
	if o[a] == b {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

// eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.

func eqrr(a, b, c int, o []int) {
	if o[a] == o[b] {
		o[c] = 1
	} else {
		o[c] = 0
	}
}

type sample struct {
	in    []int
	instr []int
	out   []int
}

/*
parseSample parses an input block like:

Before: [3, 2, 1, 1]
9 2 1 2
After:  [3, 2, 2, 1]
*/
func parseSample(lines []string) (sample, error) {
	var (
		s   sample
		err error
	)

	s.in, err = parse(lines[0], "Before: [%d, %d, %d, %d]")
	if err != nil {
		return sample{}, fmt.Errorf("error parsing before '%s': %v", lines[0], err)
	}

	s.instr, err = parse(lines[1], "%d %d %d %d")
	if err != nil {
		return sample{}, fmt.Errorf("error parsing instr '%s': %v", lines[1], err)
	}

	s.out, err = parse(lines[2], "After: [%d, %d, %d, %d]")
	if err != nil {
		return sample{}, fmt.Errorf("error parsing after '%s': %v", lines[2], err)
	}

	return s, nil
}

func candidateOps(s sample) []opNamed {
	var cands []opNamed
	for _, o := range allOps {
		reg := make([]int, len(s.in))
		copy(reg, s.in)

		o.op(s.instr[1], s.instr[2], s.instr[3], reg)

		if eq(reg, s.out) {
			cands = append(cands, o)
		}
	}
	return cands
}

func eq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type opNamed struct {
	op   op
	name string
	i    int
}

var allOps = []opNamed{
	{addr, "addr", 0},
	{addi, "addi", 1},
	{mulr, "mulr", 2},
	{muli, "muli", 3},
	{banr, "banr", 4},
	{bani, "bani", 5},
	{borr, "borr", 6},
	{bori, "bori", 7},
	{setr, "setr", 8},
	{seti, "seti", 9},
	{gtir, "gtir", 10},
	{gtri, "gtri", 11},
	{gtrr, "gtrr", 12},
	{eqir, "eqir", 13},
	{eqri, "eqri", 14},
	{eqrr, "eqrr", 15},
}

type opDecoder struct {
	h map[int]map[int]bool
}

func newOpDecoder() *opDecoder {
	od := &opDecoder{
		h: make(map[int]map[int]bool),
	}
	for i := range allOps {
		od.h[i] = make(map[int]bool)
		for j := range allOps {
			od.h[i][j] = true
		}
	}
	return od
}

func (od *opDecoder) clone() *opDecoder {
	r := newOpDecoder()
	for k := range od.h {
		for j, v := range od.h[k] {
			r.h[k][j] = v
		}
	}
	return r
}

func (od *opDecoder) addCandidates(in int, cands []opNamed) {
	if _, ok := od.h[in]; !ok {
		od.h[in] = make(map[int]bool)
	}

	possible := make(map[int]bool)
	for _, c := range cands {
		possible[c.i] = true
	}
	for i := range od.h[in] {
		od.h[in][i] = od.h[in][i] && possible[i]
	}
}

func (od *opDecoder) opsForInstruction(in int) []opNamed {
	var r []opNamed
	for i, v := range od.h[in] {
		if v {
			r = append(r, allOps[i])
		}
	}

	return r
}

func (od *opDecoder) markClaimed(in int) {
	for i := range od.h {
		od.h[i][in] = false
	}
}

func (od *opDecoder) solve() []opNamed {
	ops := make([]opNamed, len(allOps))
	known := make(map[int]bool)

	for j := range allOps {
		if known[j] {
			continue
		}
		for i := range allOps {
			if known[i] {
				continue
			}

			cands := od.opsForInstruction(i)
			if len(cands) == 1 {
				ops[i] = cands[0]
				known[i] = true
				od.markClaimed(cands[0].i)
			}
		}
	}
	return ops
}

func (od *opDecoder) String() string {
	buf := new(bytes.Buffer)
	for i := range od.h {
		fmt.Fprintf(buf, "Pos %d: ", i)
		for k, v := range od.h[i] {
			if v {
				fmt.Fprintf(buf, "%d ", k)
			}
		}
		fmt.Fprintln(buf)
	}
	return buf.String()
}

func parseInstruction(l string) ([]int, error) {
	return parse(l, "%d %d %d %d")
}

func parse(l, format string) ([]int, error) {
	var a, b, c, d int
	n, err := fmt.Sscanf(l, format, &a, &b, &c, &d)
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s': %v", l, err)
	}
	if n != 4 {
		return nil, fmt.Errorf("error parsing '%s': only decoded %d", l, n)
	}

	return []int{a, b, c, d}, nil
}
