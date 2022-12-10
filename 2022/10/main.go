// Read input from stdin, interpret the instructions; oputput sum of signal
// strengths at given cycle counts and render 40px scanlines.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	X := 1  // register X
	cc := 0 // cycle counter
	ss := 0 // signal strength

	pmp := -20          // previous measuring point
	crt := []string{""} // display output
	sl := 0             // scanline

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := strings.Split(s.Text(), " ")
		ins := l[0] // instruction
		arg := 0    // argument

		c := 0 // current instruction cycle count
		switch ins {
		case "addx":
			arg, _ = strconv.Atoi(l[1])
			c = 2
		case "noop":
			c = 1
		}

		// do cycles
		for c > 0 {
			cc++ // global cycle counter
			c--

			// signal strength
			//log.Printf("cc=%d c=%d ins=%s arg=%d X=%d", cc, c, ins, arg, X)
			if cc == pmp+40 {
				//log.Printf("measure point cc=%d X=%d", cc, X)
				ss += X * cc
				pmp = cc
			}

			// crt
			if X > len(crt[sl])-2 && X < len(crt[sl])+2 {
				crt[sl] += "#"
			} else {
				crt[sl] += "."
			}
			if cc%40 == 0 {
				crt = append(crt, "")
				sl++
			}

			// fun happens after cycle
			if c == 0 {
				X += arg
			}
		}

	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(ss)
	for i := 0; i < len(crt)-1; i++ {
		log.Print(crt[i])
	}
}
