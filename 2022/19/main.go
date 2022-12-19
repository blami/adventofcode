// Read input from stdin and find optimal blueprint for geode production as per
// rules.

package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const ( // fields in state
	O  = 0 // resources
	C  = 1
	B  = 2
	G  = 3
	RO = 4 // robots
	RC = 5
	RB = 6
	RG = 7
	T  = 8 // remaining time
)

// Find maximum of given values.
func max(a ...int) int {
	r := a[0]
	for _, n := range a {
		if n > r {
			r = n
		}
	}
	return r
}

// Check whether s is contained in sl.
// BUG: use map instead, this is too slow...
/*
func in(sl [][9]int, s [9]int) bool {
	for _, v := range sl {
		if v == s { return true }
	}
	return false
}
*/

// Maximize number of geodes given robot costs (tuples of 4 values) and time
// left.
// TODO: don't need int[4] (smaller stack) as prices have same structure (e.g.
// ore robot will always cost only ore)
// NOTE: THIS IS NOT GENERIC as it does not use all combinations of prices of
// each robot (only certain combinations of o,c,b and g components). Should
// still work for all inputs...
func maxg(o, c, b, g [4]int, t int) int {
	//log.Print(o, c, b, g, t) // print actual blueprint
	q := make([][9]int, 0)
	q = append(q, [9]int{0, 0, 0, 0, 1, 0, 0, 0, t}) // queue of states
	v := map[[9]int]bool{}                           // use map, or slice and take a nap...

	bg := 0
	for len(q) > 0 {
		s := q[0] // pop queue
		q = q[1:]

		bg = max(bg, s[G])
		if s[T] == 0 {
			continue
		} // skip states that reached t minutes

		com := max(o[0], c[0], b[0], g[0]) // max cost in ore
		if s[RO] >= com {
			s[RO] = com
		}
		if s[RC] >= b[1] {
			s[RC] = b[1]
		}
		if s[RB] >= g[2] {
			s[RB] = g[2]
		}
		if s[O] >= s[T]*com-s[RO]*(s[T]-1) {
			s[O] = s[T]*com - s[RO]*(s[T]-1)
		}
		if s[C] >= s[T]*b[1]-s[RC]*(s[T]-1) {
			s[C] = s[T]*b[1] - s[RC]*(s[T]-1)
		}
		if s[B] >= s[T]*g[2]-s[RB]*(s[T]-1) {
			s[B] = s[T]*g[2] - s[RB]*(s[T]-1)
		}

		if _, ok := v[s]; ok {
			continue
		} else {
			v[s] = true
		}

		// debug
		/*
			if len(v) % 10000 == 0 {
				log.Printf("%d bg=%d v%d q%d", s[T], bg, len(v), len(q))
				log.Print(q[len(q)-1])
			}
			if s[O] < 0 || s[C] < 0 || s[B] < 0 || s[G] < 0 {
				log.Fatal(s)
			}
		*/

		// mine
		q = append(q, [9]int{
			s[O] + s[RO],
			s[C] + s[RC],
			s[B] + s[RB],
			s[G] + s[RG],
			s[RO],
			s[RC],
			s[RB],
			s[RG],
			s[T] - 1,
		})
		// buy
		if s[O] >= o[0] { // buy ore robot
			q = append(q, [9]int{
				s[O] - o[0] + s[RO],
				s[C] + s[RC],
				s[B] + s[RB],
				s[G] + s[RG],
				s[RO] + 1,
				s[RC],
				s[RB],
				s[RG],
				s[T] - 1,
			})
		}
		if s[O] >= c[0] { // buy clay robot
			q = append(q, [9]int{
				s[O] - c[0] + s[RO],
				s[C] + s[RC],
				s[B] + s[RB],
				s[G] + s[RG],
				s[RO],
				s[RC] + 1,
				s[RB],
				s[RG],
				s[T] - 1,
			})
		}
		if s[O] >= b[0] && s[C] >= b[1] { // buy obsidian robot
			q = append(q, [9]int{
				s[O] - b[0] + s[RO],
				s[C] - b[1] + s[RC],
				s[B] + s[RB],
				s[G] + s[RG],
				s[RO],
				s[RC],
				s[RB] + 1,
				s[RG],
				s[T] - 1,
			})
		}
		if s[O] >= g[0] && s[B] >= g[2] { // buy geode robot
			q = append(q, [9]int{
				s[O] - g[0] + s[RO],
				s[C] + s[RC],
				s[B] - g[2] + s[RB],
				s[G] + s[RG],
				s[RO],
				s[RC],
				s[RB],
				s[RG] + 1,
				s[T] - 1,
			})
		}
	}

	return bg
}

func main() {

	bb := [2]int{0, 1} // best blueprint per rules

	re := regexp.MustCompile(`(?P<c>\d+ (ore|clay|obsidian))+`)
	s := bufio.NewScanner(os.Stdin)
	i := 1                // blueprint counter
	tr := map[string]int{ // translate materials to [4]int indices
		"ore":      0,
		"clay":     1,
		"obsidian": 2,
		"geode":    3,
	}
	for s.Scan() {
		// ugly input processing
		l := strings.Split(strings.Trim(strings.Split(s.Text(), ":")[1], "."), ".")
		c := [4][4]int{} // ore, clay, obsidian, geode robot costs in blueprint
		for j := range l {
			mtch := re.FindAllStringSubmatch(l[j], -1)
			c[j] = [4]int{} // cost in ore, clay, obsidian
			for k := range mtch {
				d := strings.Split(mtch[k][0], " ")
				c[j][tr[d[1]]], _ = strconv.Atoi(d[0])
			}
		}
		// solve that particular blueprint
		bb[0] += maxg(c[0], c[1], c[2], c[3], 24) * i
		if i < 4 {
			bb[1] *= maxg(c[0], c[1], c[2], c[3], 32) // * i (no quality)
		}
		i++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(bb[0])
	log.Print(bb[1])
}
