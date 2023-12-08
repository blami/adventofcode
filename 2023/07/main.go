package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
	//"sync"
	//"unicode"
)

func same(in string) int {
	cnt := map[rune]int{}
	for _, c := range in {
		cnt[c]++
	}
	r := 0
	for _, v := range cnt {
		if v > r {
			r = v
		}
	}
	return r
}

// 7 - five of a kind
// 6 - four of a kind
// 5 - full house
// 4 - three of kind
// 3 - two pair
// 2 - one pair
// 1 - high card
// 0 - nothing
func hand(in string) int {
	// count cards
	cnt := map[rune]int{}
	for _, c := range in {
		cnt[c]++
	}

	// five of a kind
	if len(cnt) == 1 {
		return 7
	}
	// four of a kind
	// full house
	if len(cnt) == 2 {
		for _, v := range cnt {
			if v == 4 {
				return 6
			}
			if v == 3 {
				return 5
			}
		}
	}
	// three of a kind
	// two pair
	if len(cnt) == 3 {
		p := 0
		for _, v := range cnt {
			if v == 3 {
				return 4
			}
			if v == 2 {
				p++
			}
		}
		if p == 2 {
			return 3
		}
	}
	// one pair
	if len(cnt) == 4 {
		p := 0
		for _, v := range cnt {
			if v == 2 {
				p++
			}
		}
		if p == 1 {
			return 2
		}
	}
	// high card
	if len(cnt) == 5 {
		return 1
	}
	return 0
}

func card(in rune) int {
	for v, c := range []rune{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'Q', 'J', 'K', 'A'} {
		if c == in {
			return v + 1
		}
	}
	return 0
}

func ncard(in rune) int {
	for v, c := range []rune{'J', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'Q', 'K', 'A'} {
		if c == in {
			return v + 1
		}
	}
	return 0
}

func cmp(a Hand, b Hand) bool {
	h1 := hand(a.Hand)
	h2 := hand(b.Hand)
	if h1 == h2 {
		for i, _ := range a.Hand {
			c1 := card(rune(a.Hand[i]))
			c2 := card(rune(b.Hand[i]))
			//fmt.Println(a, b, string(a.Hand[i]), string(b.Hand[i]), c1, c2)
			if c1 == c2 {
				continue
			}
			if c1 < c2 {
				return true
			}
			if c1 > c2 {
				return false
			}
		}
	}
	if h1 < h2 {
		return true
	}
	if h1 > h2 {
		return false
	}
	panic("oops")
}

func cmpnc(a Hand, b Hand) bool {
	h1 := hand(a.Hand)
	h2 := hand(b.Hand)
	if h1 == h2 {
		for i, _ := range a.Hand {
			c1 := ncard(rune(a.Hand[i]))
			c2 := ncard(rune(b.Hand[i]))
			//fmt.Println(a, b, string(a.Hand[i]), string(b.Hand[i]), c1, c2)
			if c1 == c2 {
				continue
			}
			if c1 < c2 {
				return true
			}
			if c1 > c2 {
				return false
			}
		}
	}
	if h1 < h2 {
		return true
	}
	if h1 > h2 {
		return false
	}
	panic("oops")
}

func ncmp(a Hand, b Hand) bool {
	h1 := hand(a.Hand)
	h2 := hand(b.Hand)
	if h1 == h2 {
		for i := range a.Hand {
			c1 := ncard(rune(a.Hand[i]))
			c2 := ncard(rune(b.Hand[i]))
			if a.Orig != "" {
				c1 = ncard(rune(a.Orig[i]))
			}
			if b.Orig != "" {
				c2 = ncard(rune(b.Orig[i]))
			}
			fmt.Println(a, b, string(a.Hand[i]), string(b.Hand[i]), c1, c2)
			if c1 == c2 {
				continue
			}
			if c1 < c2 {
				return true
			}
			if c1 > c2 {
				return false
			}
		}
	}
	if h1 < h2 {
		return true
	}
	if h1 > h2 {
		return false
	}
	fmt.Println(a.Hand, b.Hand, h1, h2)
	panic("oops")
	//return false
}

// TODO: move this to common utilities (generic?)
func combs(set []rune, l int) [][]rune {
	var r [][]rune
	var cur []rune

	var bt func(si int)
	bt = func(si int) {
		if len(cur) == l {
			tc := make([]rune, l)
			copy(tc, cur)
			r = append(r, tc)
			return
		}

		for i := si; i < len(set); i++ {
			cur = append(cur, set[i])
			bt(i)
			cur = cur[:len(cur)-1]
		}
	}

	bt(0)
	return r
}

func jokers(h Hand) Hand {
	cards := []rune{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'Q', 'K', 'A'}

	n := strings.Count(h.Hand, "J")
	if n == 0 {
		return h
	}
	is := []int{}
	for i := range h.Hand {
		if h.Hand[i] == 'J' {
			is = append(is, i)
		}
	}
	hh := []rune(h.Hand)
	l := []Hand{}
	for _, comb := range combs(cards, n) {
		for j, i := range is {
			hh[i] = comb[j]
		}
		l = append(l, Hand{string(hh), h.Bid, h.Hand})
	}

	// sort and find best
	// need to always use old sort
	slices.SortFunc(l, cmpnc)
	fmt.Println(l)
	return l[len(l)-1]
}

type Hand struct {
	Hand string
	Bid  int
	Orig string
}

func main() {
	hands := []Hand{}
	nhands := []Hand{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		ll := strings.Split(l, " ")
		b, _ := strconv.Atoi(ll[1])
		h := Hand{ll[0], b, ""}
		hands = append(hands, h)
		nhands = append(nhands, jokers(h))
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(hands, cmp)
	slices.SortFunc(nhands, ncmp)

	win := 0
	nwin := 0
	for rank, hand := range hands {
		//fmt.Println(hand.Hand, rank + 1, "*", hand.Bid)
		win += hand.Bid * (rank + 1)
	}
	for rank, hand := range nhands {
		fmt.Println(hand.Hand, hand.Orig, rank+1, "*", hand.Bid)
		nwin += hand.Bid * (rank + 1)
	}

	fmt.Println(win)
	fmt.Println(nwin)
	//fmt.Println(jokers(Hand{"KTJJT", 1, ""}))
	//
	fmt.Println(combs([]rune{'a', 'b', 'c'}, 3))
}
