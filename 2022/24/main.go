package main

import (
	"bufio"
	"log"
	"os"
)

// XY type alias.
type XY [2]int

// Direction deltas, convenient map for moving blizzards...
var D map[byte]XY = map[byte]XY{
	'#': XY{0, 0}, // convenient for BFS wait
	'^': XY{0, -1},
	'>': XY{1, 0},
	'v': XY{0, 1},
	'<': XY{-1, 0},
}

// BFS to find least minutes from b to e in m of w*h size, starting at time t.
func bfs(m map[XY]byte, w int, h int, b XY, e XY, t int) int {
	q := [][3]int{[3]int{b[0], b[1], t}} // queue, state is expedition X, Y, time
	v := map[[3]int]bool{}
	for len(q) > 0 {
		c := q[0]
		q = q[1:] // pop

	DL:
		for _, d := range D {
			n := [3]int{c[0] + d[0], c[1] + d[1], c[2] + 1} // new state
			if n[0] == e[0] && n[1] == e[1] {               // end
				return n[2]
			}
			if _, ok := v[n]; ok {
				continue
			}
			if ch, ok := m[XY{n[0], n[1]}]; !ok || ch == '#' { // walls, outside map
				continue
			}
			// blizzards
			if n[0] > 0 && n[1] > 0 && n[0] < w-1 && n[1] < h-1 { // BUG: only check where blizzards move!
				for ch, d := range D {
					// check if in time n[2] will be at XY{n[0], n[1]} blizzard
					// rectangle min is 1,1 and max is w-2,h-2
					chx := (n[0] - (d[0] * n[2]) - 1) % (w - 2)
					chy := (n[1] - (d[1] * n[2]) - 1) % (h - 2)
					if chx < 0 {
						chx += (w - 2)
					}
					if chy < 0 {
						chy += (h - 2)
					}
					if m[XY{chx + 1, chy + 1}] == ch {
						continue DL
					}
				}
			}
			v[n] = true
			q = append(q, n)
		}
	}
	return 0
}

func main() {
	m := make(map[XY]byte) // map
	w := 0
	h := 0

	s := bufio.NewScanner(os.Stdin)
	y := 0
	for s.Scan() {
		l := s.Text()
		for x, ch := range l {
			switch ch {
			case '^', 'v', '<', '>', '#', '.': // storing also . to simplify bounds check
				m[XY{x, y}] = byte(ch)
			}
		}
		if w == 0 {
			w = len(l)
		}
		h = y + 1
		y++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	b := XY{1, 0}
	e := XY{w - 2, h - 1}

	log.Print(bfs(m, w, h, b, e, 0))                                         // part 1
	log.Print(bfs(m, w, h, b, e, bfs(m, w, h, e, b, bfs(m, w, h, b, e, 0)))) // part 2
}
