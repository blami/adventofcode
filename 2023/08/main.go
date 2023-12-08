// Find number of steps from AAA to ZZZ and from all **A to all **Z.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Instr struct {
	L string
	R string
}

// TODO: Move gcd and lcm to utilities
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(x ...int) int {
	// BUG: len(x) >= 2
	r := x[0] * x[1] / gcd(x[0], x[1])
	for i := 0; i < len(x[2:]); i++ {
		r = lcm(r, x[i+2])
	}
	return r
}

func main() {
	dir := ""
	instr := map[string]Instr{}
	start := "AAA"
	end := "ZZZ"

	starta := []string{}

	s := bufio.NewScanner(os.Stdin)
	lin := 0
	for s.Scan() {
		l := s.Text()
		if lin == 0 {
			dir = l
		}
		if lin > 1 {
			l = strings.ReplaceAll(l, " = ", " ")
			l = strings.ReplaceAll(l, "(", "")
			l = strings.ReplaceAll(l, ")", "")
			l = strings.ReplaceAll(l, ",", " ")
			lf := strings.Split(l, " ")
			instr[lf[0]] = Instr{L: lf[1], R: lf[3]}
			// OOPS!
			//if lin == 2 {
			//	start = lf[0]
			//}
			//end = lf[0]

			// part 2
			if lf[0][2] == 'A' {
				starta = append(starta, lf[0])
			}
		}
		lin++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	dirp := 0
	st := 0
	cur := start
	for cur != end {

		switch dir[dirp] {
		case 'R':
			cur = instr[cur].R
		case 'L':
			cur = instr[cur].L
		default:
			panic("oops")
		}

		st++
		// next direction
		dirp++
		if dirp > len(dir)-1 {
			dirp = 0
		}
	}
	fmt.Println(st)

	// part 2
	// NOTE: was too lazy to mix both parts together...
	st = 0
	dirp = 0
	curz := make([]string, len(starta))
	copy(curz, starta)

	// total number of steps is lcm of individual step numbers from **A to **Z
	// in each lane.
	stz := make([]int, len(curz))
	for {
		st++
		for i, c := range curz {
			switch dir[dirp] {
			case 'R':
				curz[i] = instr[c].R
			case 'L':
				curz[i] = instr[c].L
			default:
				panic("oops")
			}
			if curz[i][2] == 'Z' {
				if stz[i] == 0 {
					stz[i] = st
				}
			}
		}

		done := true
		// BUG: ==0 assumption
		for i := range stz {
			if stz[i] == 0 {
				done = false
			}
		}
		if done {
			break
		}

		// next direction
		dirp++
		if dirp > len(dir)-1 {
			dirp = 0
		}
	}

	fmt.Println(lcm(stz...))
}
