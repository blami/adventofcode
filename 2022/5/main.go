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

// An alias for representing stack of crates
type Stacks [][]byte

func (s Stacks) Topmost() string {
	r := ""
	for _, v := range s {
		if len(v) > 0 {
			r += string([]byte{v[0]})
		}
	}
	return r
}

// Initialize stacks from header string.
// NOTE: This assumes always 9 stacks, all crates single letter formatted as in
// examples. Strict spacing is expected.
// NOTE: I hate this but it had to be done so that elves can continue...
func initStacks(hdr []string) Stacks {
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

// Reverse given slice of bytes (for moving one by one).
func reverse(s []byte) []byte {
	r := make([]byte, len(s))
	for i, v := range s {
		r[len(s)-i-1] = v
	}
	return r
}

// Move n crates from stack to stack in stacks, either one by one (CrateMover
// 9000), or keeping their original order (CrateMover 9001).
func move(stacks Stacks, n, from, to int, keep bool) Stacks {
	buf := append([]byte(nil), stacks[from][:n]...) // need to copy to avoid reuse of underlying array
	if !keep {
		buf = reverse(buf)
	}
	stacks[from] = stacks[from][n:]
	stacks[to] = append(buf, stacks[to]...)
	return stacks
}

func main() {
	// stacks aren't that big so no pointers needed (fired)
	var stacks9000, stacks9001 Stacks

	hdr := []string{} // stacks visualization header
	body := false     // body parsing flag
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if !body {
			if l == "" {
				body = true                  // switch to parsing body
				stacks9000 = initStacks(hdr) // initialize stacks
				stacks9001 = initStacks(hdr) // initialize again to avoid writing deepcopy lol
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
		stacks9000 = move(stacks9000, n, from-1, to-1, false)
		stacks9001 = move(stacks9001, n, from-1, to-1, true)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(stacks9000.Topmost())
	log.Print(stacks9001.Topmost())
}
