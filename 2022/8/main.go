// Read input from stdin as heightmap row by row. Find number of trees that can
// be seen from border and also highest scenic score (mul of number of trees
// seen from tree in all directions).
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Type alias for heightmap
type Map [][]int

// Look in all directions from tree and return in how many directions the
// forest border can be seen and also how far one can see (scenic score).
func (m *Map) Look(x, y int) (int, int) {
	h := (*m)[y][x]
	v := 4
	s := [4]int{} // scenic score components 0-up, 1-down, 2-left, 3-right

	// lol. this is so ugly and reminds me of similar purpose function I wrote
	// in Visual Basic for one of my early game attempt.

	for yy := y - 1; yy >= 0; yy-- {
		s[0]++
		if (*m)[yy][x] >= h {
			v--
			break
		}
	}
	// check if the column to bottom is all lower
	for yy := y + 1; yy < len(*m); yy++ {
		s[1]++
		if (*m)[yy][x] >= h {
			v--
			break
		}
	}
	// check if the column to left is all lower
	for xx := x - 1; xx >= 0; xx-- {
		s[2]++
		if (*m)[y][xx] >= h {
			v--
			break
		}
	}
	// check if the column to right is all lower
	for xx := x + 1; xx < len((*m)[0]); xx++ {
		s[3]++
		if (*m)[y][xx] >= h {
			v--
			break
		}
	}
	return v, s[0] * s[1] * s[2] * s[3]
}

func main() {

	var m Map // map of forest

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var r []int
		for _, t := range strings.Split(s.Text(), "") {
			tv, _ := strconv.Atoi(t)
			r = append(r, tv)
		}
		m = append(m, r)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// map dimensions
	w := len(m[0])
	h := len(m)

	// go tree by tree and try look from that tree at border, if possible count
	// it; on the way also count tree's scenic score
	v := 0
	hss := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			vd, ss := m.Look(x, y)
			if vd >= 1 { // visible from at least one border
				v += 1
			}
			if ss > hss {
				hss = ss
			}
		}
	}

	log.Print(v)
	log.Print(hss)
}
