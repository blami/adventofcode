// Read input from stdin

package main

import (
	"bufio"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"unicode"
)

// Find items shared between two containers (be it whole rucksacks or their
// compartments.
func shared(h1, h2 []byte) []byte {
	// NOTE: assuming always correct input
	slices.Sort(h1)
	slices.Sort(h2)
	s := make([]byte, 0)

	i := 0
	j := 0
	for i < len(h1) && j < len(h2) {
		switch {
		case h1[i] == h2[j]:
			s = append(s, h1[i])
			i++
			j++
		case h1[i] < h2[j]:
			i++
		default:
			j++
		}
	}

	return s
}

// Rank given item by priority
func rank(i byte) int {
	if unicode.IsUpper(rune(i)) {
		return int(i) - 38 // A-Z is 27-52
	}
	return int(i) - 96
}

func main() {

	sump := 0  // sum of shared item priorities
	sumpb := 0 // sum of badge priorities

	s := bufio.NewScanner(os.Stdin)
	g := make([][]byte, 3)
	n := 0
	for s.Scan() {
		r := s.Bytes()
		// NOTE: we know exactly one item is shared between compartments of
		// single rucksack so its safe to rank [0]
		sump += rank(shared(r[:len(r)/2], r[len(r)/2:])[0])

		g[n] = r
		n++
		if n == 3 {
			// NOTE: again, we know there's only one badge item
			sumpb += rank(shared(shared(g[0], g[1]), g[2])[0])
			n = 0
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(sump)
	log.Print(sumpb)
}
