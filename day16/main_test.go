package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestCandidateOps(t *testing.T) {
	cc := []struct {
		in []string
		ex []string
	}{
		{[]string{"Before: [3, 2, 1, 1]", "9 2 1 2", "After:  [3, 2, 2, 1]"}, []string{"addi", "mulr", "seti"}},
	}

	for i, c := range cc {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			s, err := parseSample(c.in)
			if err != nil {
				t.Errorf("unexpected error parsing sample: %v", err)
			}

			oo := candidateOps(s)
			var act []string
			for _, o := range oo {
				act = append(act, o.name)
			}
			sort.Strings(act)
			if !reflect.DeepEqual(act, c.ex) {
				t.Errorf("got %v expected %v", act, c.ex)
			}
		})
	}
}

func TestOpDecoder(t *testing.T) {
	od := newOpDecoder()
	od.addCandidates(0, []opNamed{allOps[0], allOps[1]})
	od.addCandidates(0, []opNamed{allOps[2], allOps[1]})
	od.addCandidates(1, []opNamed{allOps[0], allOps[3]})
	od.addCandidates(1, []opNamed{allOps[0], allOps[3], allOps[4]})

	ops := od.opsForInstruction(0)
	expect(t, true, verifyOpsContains(ops, allOps[1]), "0 1")
	expect(t, false, verifyOpsContains(ops, allOps[0]), "0 0")
	expect(t, false, verifyOpsContains(ops, allOps[2]), "0 2")
	expect(t, false, verifyOpsContains(ops, allOps[3]), "0 3")

	ops = od.opsForInstruction(1)
	expect(t, true, verifyOpsContains(ops, allOps[0]), "1 0")
	expect(t, true, verifyOpsContains(ops, allOps[3]), "1 3")
	expect(t, false, verifyOpsContains(ops, allOps[1]), "1 1")
	expect(t, false, verifyOpsContains(ops, allOps[2]), "1 2")
}

func expect(t *testing.T, act, exp bool, msg string) {
	if act != exp {
		t.Errorf("%s: got %v expected %v", msg, act, exp)
	}
}

func verifyOpsContains(ops []opNamed, op opNamed) bool {
	for _, o := range ops {
		if o.i == op.i {
			return true
		}
	}
	return false
}
