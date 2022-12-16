package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
)

// XYZ alias.
type XYZ [3]int

func main() {
	cs := []XYZ{} // cubes

	sut := 0 // total surface
	suc := 0 // connected surface

	var b [2]XYZ // min and max coordinates found

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		c := XYZ{}
		for i, v := range strings.Split(s.Text(), ",") {
			c[i], _ = strconv.Atoi(v)
		}
		sut += 6

		// place the cube, if any of existing cubes neighbors with cube remove
		// 2 sides from total surface for each neighboring side
		for _, ec := range cs {
			if c[1] == ec[1] && c[2] == ec[2] && (c[0] + 1 == ec[0] || c[0] - 1 == ec[0]) { suc+=2 }
			if c[0] == ec[0] && c[2] == ec[2] && (c[1] + 1 == ec[1] || c[1] - 1 == ec[1]) { suc+=2 }
			if c[0] == ec[0] && c[1] == ec[1] && (c[2] + 1 == ec[2] || c[2] - 1 == ec[0]) { suc+=2 }
		}
		cs = append(cs, c)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(sut-suc) // part 1
}
