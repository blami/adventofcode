// Read input from stdin and store heightmap in 2D matrix. Find number of trees
// that can be seen from border and also best scenic score (product of visible
// trees in all directions).

package main

import (
	"bufio"
	"log"
	"os"
)

// A heightmap.
type Map [][]int

// Look in all directions and return in how many directions the forest border
// can be seen and also scenic score.
func (m *Map) Look(x, y int) (int, int) {
	h := (*m)[y][x]
	v := 4
	s := [4]int{} // scenic score components 0-up, 1-down, 2-left, 3-right

	// NOTE: this is so ugly and reminds me of similar purpose function I wrote
	// in Visual Basic for one of my early game attempt.

	// up
	for yy := y - 1; yy >= 0; yy-- {
		s[0]++
		if (*m)[yy][x] >= h {
			v--
			break
		}
	}
	// down
	for yy := y + 1; yy < len(*m); yy++ {
		s[1]++
		if (*m)[yy][x] >= h {
			v--
			break
		}
	}
	// left
	for xx := x - 1; xx >= 0; xx-- {
		s[2]++
		if (*m)[y][xx] >= h {
			v--
			break
		}
	}
	// right
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

	var m Map

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var r []int
		for _, b := range s.Bytes() {
			r = append(r, int(b-48)) // no strconv needed, 0=48..9=57
		}
		m = append(m, r)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	w := len(m[0])
	h := len(m)

	// go tree by tree and try look from that tree at border, if possible count
	// it; on the way also count tree's scenic score
	v := 0   // trees visible from border
	hss := 0 // highest scenic score
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			vv, ss := m.Look(x, y)
			if vv >= 1 { // visible from at least one border
				v += 1
			}
			if ss > hss { // highest scenic score
				hss = ss
			}
		}
	}

	log.Print(v)
	log.Print(hss)
}
