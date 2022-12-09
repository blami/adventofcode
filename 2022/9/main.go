// Read input from stdin and simulate head and tail(s) movements. Count all
// unique coordinates tail(s) visit.

package main

import (
	"bufio"
	"log"
	//"math"
	"os"
	"strconv"
)

// Coordinates
type XY [2]int

// Chase given head with tail according to the rules.
func (t *XY) Chase(h XY) {
	// NOTE: uh-oh, didn't find better way to avoid zillion ifs
	d := XY{h[0] - t[0], h[1] - t[1]} // x,y difference between head and tail
	switch d {
	// orthogonal
	case XY{0, -2}: // up
		(*t)[1]--
	case XY{0, 2}: // down
		(*t)[1]++
	case XY{-2, 0}: // left
		(*t)[0]--
	case XY{2, 0}: // right
		(*t)[0]++
	// diagonal (there are two types for each + 2,2 distances for tails)
	case XY{1, -2}, XY{2, -1}, XY{2, -2}: // up right
		(*t)[0]++
		(*t)[1]--
	case XY{-1, -2}, XY{-2, -1}, XY{-2, -2}: // up left
		(*t)[0]--
		(*t)[1]--
	case XY{1, 2}, XY{2, 1}, XY{2, 2}: // down right
		(*t)[0]++
		(*t)[1]++
	case XY{-1, 2}, XY{-2, 1}, XY{-2, 2}: // down left
		(*t)[0]--
		(*t)[1]++
		// debug unhandled cases
		// NOTE: I used this to debug part 2 of excercise
		/*
			default:
				if math.Abs(float64(d[0])) >= 2 || math.Abs(float64(d[1])) >= 2 {
					log.Fatal("unhandled: ", d)
				}
		*/
	}
}

func main() {

	// head and tail positions (x,y), head is tail 0
	t := [10]XY{} // tails (knots)
	h := &t[0]

	tup := []map[XY]int{{t[0]: 1}, {t[len(t)-1]: 1}} // tail unique positions 0-first, 1-ninth

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		d := byte(s.Text()[0]) // direction of head
		// XXX: I had my solution rejected and this was the bug, everything was
		// right logic wise but here I did (sleep deprivation)
		// strings.Split(s.Text() , "") which would split string not by " "
		// column but would 1) correctly yield direction, 2) split any 2 or
		// more digit number in 2nd column and take only first digit 3) NOT
		// BREAK EXAMPLE AS IT USES ONLY SINGLE DIGITS... Geez.
		n, _ := strconv.Atoi(s.Text()[2:]) // number of steps

		// step by step
		for n > 0 {
			// move head
			switch d {
			case 'U':
				h[1]++
			case 'D':
				h[1]--
			case 'L':
				h[0]--
			case 'R':
				h[0]++
			}

			// chase with tails
			for i := 1; i < len(t); i++ {
				t[i].Chase(t[i-1])
			}

			tup[0][t[1]] += 1
			tup[1][t[len(t)-1]] += 1

			n--
		}

	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(len(tup[0]))
	log.Print(len(tup[1]))
}
