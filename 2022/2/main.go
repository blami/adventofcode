// Read input from stdin and calculate score for each round of rock, papper,
// scissors game.

package main

import (
	"bufio"
	"log"
	"os"
)

// Score the round as sum of own shape and outcome of round according to rules.
// NOTE: Works with 1,2,3 normalized values.
func score(p1, p2 byte) int {
	s := 0

	// score own shape
	s += int(p2)

	// score the match
	switch {
	// draws
	case p1 == p2:
		s += 3
	// wins
	case (p2-1 == p1) || (p2 == 1 && p1 == 3):
		s += 6
		// loses (no need to run through this branch as score is 0
		//case (p2 == p1 - 1) || (p1 == 1 && p2 == 3):
	}

	return s
}

// Pick shape given what p1 plays so that rounds ends up with expected outcome
// e. NOTE: Works with 1,2,3 normalized values.
func pick(p1 byte, e byte) byte {
	var x byte = 0
	switch e {
	// draw
	case 2:
		return p1
	// win
	case 3:
		x = p1 + 1
	// lose
	case 1:
		x = p1 - 1
	}
	if x > 3 {
		return 1
	}
	if x < 1 {
		return 3
	}
	return x
}

func main() {
	sum1 := 0 // total score for part 1
	sum2 := 0 // total score for part 2

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		in := s.Bytes()
		// normalize inputs, both ABC and XYZ are sequentially going bytes with
		// ABC starting at 65 and XYZ at 88, subtracting that -1 gives
		// comparable 1,2,3 sequences
		sum1 += score(in[0]-64, in[2]-87)

		// for part 2 the second column is round outcome rather than
		// shape so pick() right shape to play to that expected outcome
		sum2 += score(in[0]-64, pick(in[0]-64, in[2]-87))
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(sum1)
	log.Print(sum2)
}
