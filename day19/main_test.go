package main

import (
	"bytes"
	"testing"
)

func TestRunProgram(t *testing.T) {
	pcr, ins, err := parseProgram(`#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5
`)
	if err != nil {
		t.Fatalf("unexpected error parsing program: %v", err)
	}
	exp := `ip=0 [0, 0, 0, 0, 0, 0] seti 5 0 1 [0, 5, 0, 0, 0, 0]
ip=1 [1, 5, 0, 0, 0, 0] seti 6 0 2 [1, 5, 6, 0, 0, 0]
ip=2 [2, 5, 6, 0, 0, 0] addi 0 1 0 [3, 5, 6, 0, 0, 0]
ip=4 [4, 5, 6, 0, 0, 0] setr 1 0 0 [5, 5, 6, 0, 0, 0]
ip=6 [6, 5, 6, 0, 0, 0] seti 9 0 5 [6, 5, 6, 0, 0, 9]
`

	cpu := NewCPU()
	cpu.PCR = pcr
	buf := new(bytes.Buffer)
	runProgram(cpu, ins, buf)
	if buf.String() != exp {
		t.Errorf("got %s expected %s", buf, exp)
	}
}
