// Read input from stdin, parse the initial state of crates on stacks and
// execute instructions to move by both CrateMover9000 (one by one) and
// CrateMover9001 (multiple at once).

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Reverse given slice of bytes.
func reverse(sl []byte) []byte {
	rsl := make([]byte, len(sl))
	for i, v := range sl {
		rsl[len(sl)-i-1] = v
	}
	return rsl
}

// Stacks of crates.
type Stacks [][]byte

// Initialize stacks from the header string.
// NOTE: This assumes always 9 stacks, all crates single letter formatted as in
// examples. Strict spacing is expected. Ugly but works.
func NewStacks(hdr []string) Stacks {
	stacks := Stacks{}
	// crate position maps
	cp := map[int]int{1: 0, 5: 1, 9: 2, 13: 3, 17: 4, 21: 5, 25: 6, 29: 7, 33: 8}

	for range [9]struct{}{} {
		stacks = append(stacks, []byte{})
	}
	// populate stacks with crates, topmost crate first
	for i, l := range hdr {
		if i >= len(hdr)-1 {
			break
		}
		for j, c := range []byte(l) {
			n, ok := cp[j]
			if ok && c != ' ' {
				stacks[n] = append(stacks[n], c)
			}
		}
	}

	return stacks
}

// Return top of the stack crate letters (left to right).
func (s *Stacks) Top() string {
	r := ""
	for _, v := range *s {
		if len(v) > 0 {
			r += string([]byte{v[0]})
		}
	}
	return r
}

// Move n crates from stack to stack, either one by one (CrateMover 9000), or
// keeping their original order (CrateMover 9001).
func (s *Stacks) Move(n, from, to int, keep bool) {
	buf := append([]byte(nil), (*s)[from][:n]...) // need to copy to avoid reuse of underlying array
	if !keep {
		buf = reverse(buf)
	}
	(*s)[from] = (*s)[from][n:]
	(*s)[to] = append(buf, (*s)[to]...)
}

func main() {
	var stacks9000, stacks9001 Stacks

	hdr := []string{} // stacks visualization header
	body := false     // body parsing flag
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if !body {
			if l == "" {
				body = true // switch to parsing body
				// parsing is insignificant toil here...
				stacks9000 = NewStacks(hdr)
				stacks9001 = NewStacks(hdr)
				continue
			}
			hdr = append(hdr, l)
			continue
		}

		instr := strings.Split(l, " ")
		n, _ := strconv.Atoi(instr[1])
		from, _ := strconv.Atoi(instr[3])
		to, _ := strconv.Atoi(instr[5])

		// damned elves, they index things from 1
		stacks9000.Move(n, from-1, to-1, false)
		stacks9001.Move(n, from-1, to-1, true)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(stacks9000.Top())
	log.Print(stacks9001.Top())
}
