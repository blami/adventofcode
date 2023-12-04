// Find sum of points of winning cards in original deck and count all cards in
// won cards deck.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Id     int
	Wns    []int
	Ons    []int
	Orig   bool
	Copies int
}

func main() {
	sum := 0
	cnt := 0
	deck := []Card{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		// poor man's input parsing, grab id, winning and own numbers
		id, _ := strconv.Atoi(strings.Fields(l)[1][0 : len(strings.Fields(l)[1])-1])
		delim := strings.Index(l, ":")
		l = l[delim+1:]

		wns := []int{}
		ons := []int{}
		inoc := false
		for _, t := range strings.Fields(l) {
			if t == "|" {
				inoc = true
				continue
			}
			n, _ := strconv.Atoi(t)
			if !inoc {
				wns = append(wns, n)
			} else {
				ons = append(ons, n)
			}
		}

		deck = append(deck, Card{
			Id:     id,
			Wns:    wns,
			Ons:    ons,
			Orig:   true,
			Copies: 1, // part 2
		})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// loop over deck until there are no changes
	ch := true
	for ch {
		ch = false
		for i := range deck {
			//debug("", i, deck)

			pts := 0
			first := true
			ws := 0

			// deplete all copies of card
			// part 2
			cps := deck[i].Copies
			if cps == 0 {
				continue
			}
			deck[i].Copies = 0

			// fmt.Printf("card %d\n", deck[i].Id)
			for _, wn := range deck[i].Wns {
				for _, on := range deck[i].Ons {
					if on == wn {
						if deck[i].Orig {
							if first {
								pts += 1
								first = false
							} else {
								pts = pts * 2
							}
						}
						ws += 1
					}
				}
			}
			deck[i].Orig = false

			// do copying (add copies to card)
			//fmt.Printf("  %d wins, copying:\n", ws)
			for j := i + 1; j < i+1+ws; j++ {
				//fmt.Printf(" %d", deck[j].Id)
				for k := 0; k < len(deck); k++ {
					if deck[k].Id == deck[j].Id {
						deck[k].Copies += cps
						ch = true
					}
				}
			}
			//fmt.Println()

			cnt += cps
			sum += pts
		}
	}

	fmt.Println(sum)
	fmt.Println(cnt)
}

// func debug(s string, I int, deck []Card) {
// 	fmt.Println(s, "---")
// 	for i, c := range deck {
// 		x := ' '
// 		if i == I { x = '>' }
// 		fmt.Printf("%c[%d] card %d: copies: %d\n", x, i, c.Id, c.Copies)
// 	}
// }
