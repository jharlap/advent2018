package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	ll, err := readLines("../inputs/19.txt")
	if err != nil {
		fmt.Println("error reading instructions:", err)
		os.Exit(1)
	}

	pcr, ins, err := parseProgram(ll)
	if err != nil {
		fmt.Println("error parsing program:", err)
		os.Exit(1)
	}

	// step 1
	cpu := NewCPU()
	cpu.PCR = pcr
	runProgram(cpu, ins, nil)
	fmt.Println("Unoptimized Step 1 R0:", cpu.Reg[0])

	r0 := optimizedProgram(0)
	fmt.Println("Step 1 R0:", r0)

	r0 = fullyOptimizedProgram(0)
	fmt.Println("Step 1 R0:", r0)

	// step 2
	r0 = fullyOptimizedProgram(1)
	fmt.Println("Step 2 R0:", r0)

	//r0 = optimizedProgram(1)
	//fmt.Println("Semi-optimized Step 2 R0:", r0)
}

func fullyOptimizedProgram(R0 int) int {
	R4 := 893
	if R0 == 1 {
		R4 = 10551293
	}

	R0 = 0
	for i := 1; i <= R4; i++ {
		if R4%i == 0 {
			R0 += i
		}
	}
	return R0
}

func factor(x int) []int {
	var ff []int
	for i := 1; i < x/2; i++ {
		if x%i == 0 {
			f := x / i
			ff = append(ff, f, i)
		}
	}
	return ff
}

func optimizedProgram(R0 int) int {
	var R1, R2, R3, R4, i int

	// R5 = R5+16 // 0: addi 5 16 5
	goto L17

L1:
	R2 = 1 // 1: seti 1 1 2

L2:
	R1 = 1 // 2: seti 1 8 1

L3:
	R3 = R2 * R1 // 3: mulr 2 1 3

	if R3 == R4 {

		// R3 = R3 == R4 // 4: eqrr 3 4 3
		// R5 = R3+R5 // 5: addr 3 5 5
		// R5++ // 6: addi 5 1 5
		R0 = R2 + R0 // 7: addr 2 0 0
	}

	R1 = R1 + 1 // 8: addi 1 1 1

	if i > 0 && i%100000000 == 0 {
		fmt.Println(i, R0, R1, R2, R3, R4)
	}
	i++

	if R1 <= R4 {
		// R3 = R1 > R4 // 9: gtrr 1 4 3
		// R5 = R5 + R3 // 10: addr 5 3 5
		// R5 = 2 // 11: seti 2 6 5
		goto L3
	}

	R2 = R2 + 1 // 12: addi 2 1 2

	if R2 <= R4 {
		// R3 = R2 > R4 // 13: gtrr 2 4 3
		// R5 = R3 + R5 // 14: addr 3 5 5
		// R5 = 1 // 15: seti 1 2 5
		goto L2
	}

	//R5 = R5 * R5 // 16: mulr 5 5 5
	//L16:
	return R0

L17:
	//R4 = R4 + 2  // 17: addi 4 2 4
	//R4 = R4 * R4 // 18: mulr 4 4 4
	//R4 = 19 * R4 // 19: mulr 5 4 4
	//R4 = R4 * 11 // 20: muli 4 11 4
	//R3 = R3 + 2  // 21: addi 3 2 3
	//R3 = R3 * 22 // 22: mulr 3 5 3
	//R3 = R3 + 13 // 23: addi 3 13 3
	//R4 = R4 + R3 // 24: addr 4 3 4

	R4 = 893 // R4 = (R4+2)*(R4+2)*19*11 + (R3+2)*22 + 13

	// R5 = R5 + R0 // 25: addr 5 0 5
	//goto ????
	if R0 == 0 {
		// R5 = 0 // 26: seti 0 8 5
		goto L1
	} else if R0 != 1 {
		fmt.Println("L25! R0=", R0)
		return -1
	}

	//R3 = 27      // 27: setr 5 5 3
	//R3 = R3 * 28 // 28: mulr 3 5 3
	//R3 = 29 + R3 // 29: addr 5 3 3
	//R3 = 30 * R3 // 30: mulr 5 3 3
	//R3 = R3 * 14 // 31: muli 3 14 3
	//R3 = R3 * 32 // 32: mulr 3 5 3
	//R4 = R4 + R3 // 33: addr 4 3 4

	R4 = 10551293 //R4 = R4 + ((27*28+29)*30*14)*32

	R0 = 0 // 34: seti 0 9 0

	//R5 = 0       // 35: seti 0 9 5
	goto L1

}

func parseProgram(input string) (int, []instruction, error) {
	var (
		pcr int
		ins []instruction
	)
	for i, l := range strings.Split(input, "\n") {
		if len(l) == 0 {
			continue
		}

		if i == 0 {
			pc, err := parseIP(l)
			if err != nil {
				return 0, nil, fmt.Errorf("error parsing IP: %s err: %v", l, err)
			}
			pcr = pc
			continue
		}

		in, err := parseInstruction(l)
		if err != nil {
			return 0, nil, fmt.Errorf("error parsing instruction: %s err: %v", l, err)
		}

		ins = append(ins, in)
	}
	return pcr, ins, nil
}

func runProgram(cpu *CPU, ins []instruction, debug io.Writer) {
	for isValidPC(cpu.PC, ins) {
		cpu.Reg[cpu.PCR] = cpu.PC
		if debug != nil {
			fmt.Fprintf(debug, "ip=%d [%d, %d, %d, %d, %d, %d] ", cpu.PC, cpu.Reg[0], cpu.Reg[1], cpu.Reg[2], cpu.Reg[3], cpu.Reg[4], cpu.Reg[5])
		}
		cpu.Exec(ins[cpu.PC].op, ins[cpu.PC].a, ins[cpu.PC].b, ins[cpu.PC].c)
		if debug != nil {
			fmt.Fprintf(debug, "%s %d %d %d [%d, %d, %d, %d, %d, %d]\n", ins[cpu.PC].op, ins[cpu.PC].a, ins[cpu.PC].b, ins[cpu.PC].c, cpu.Reg[0], cpu.Reg[1], cpu.Reg[2], cpu.Reg[3], cpu.Reg[4], cpu.Reg[5])
		}
		cpu.PC = cpu.Reg[cpu.PCR]
		cpu.PC++
	}
}
func readLines(filename string) (string, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(bb), nil
}

func isValidPC(pc int, ins []instruction) bool {
	return pc >= 0 && pc < len(ins)
}

type CPU struct {
	Reg   []int
	instr map[string]op
	PC    int
	PCR   int
}

func NewCPU() *CPU {
	return &CPU{
		Reg: make([]int, 6),
		instr: map[string]op{
			"addr": addr,
			"addi": addi,
			"mulr": mulr,
			"muli": muli,
			"banr": banr,
			"bani": bani,
			"borr": borr,
			"bori": bori,
			"setr": setr,
			"seti": seti,
			"gtir": gtir,
			"gtri": gtri,
			"gtrr": gtrr,
			"eqir": eqir,
			"eqri": eqri,
			"eqrr": eqrr,
		},
	}
}

func (cpu *CPU) Exec(o string, a, b, c int) {
	cpu.instr[o](a, b, c, cpu.Reg)
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

type instruction struct {
	op      string
	a, b, c int
}

func parseIP(l string) (int, error) {
	var v int
	n, err := fmt.Sscanf(l, "#ip %d", &v)
	if err != nil {
		return -1, fmt.Errorf("error parsing '%s': %v", l, err)
	}
	if n != 1 {
		return -1, fmt.Errorf("error parsing '%s': expected N 1 got %d", l, n)
	}
	return v, nil
}

func parseInstruction(l string) (instruction, error) {
	var in instruction

	n, err := fmt.Sscanf(l, "%s %d %d %d", &in.op, &in.a, &in.b, &in.c)
	if err != nil {
		return in, fmt.Errorf("error parsing '%s': %v", l, err)
	}
	if n != 4 {
		return in, fmt.Errorf("error parsing '%s': expected N %d got %d", l, 4, n)
	}

	return in, nil
}
