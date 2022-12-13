// Read input from stdin. Search shortest path between S and E and then between
// E and 'a' level.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"image"
)

// X,Y alias.
type XY [2]int

// Flag to produce gif image (slooow)
var vis bool

// Debug prints of current state.
// NOTE: I missed out very important sentence in riddle assignment which led me
// to extensive debugging of correct program...
func debugDump(m [][]int, v map[XY]bool, d map[XY]int) {
	var dm, vm, dim string

	for y, l := range m {
		for x, h := range l {
			ch := string(byte(h) + 'a')
			dm += ch

			if _, ok := v[XY{x, y}]; ok {
				ch = "."
			}
			vm += ch

			di, _ := d[XY{x, y}]
			dim += fmt.Sprintf("%3d ", di)
		}
		dm += "\n"
		vm += "\n"
		dim += "\n"
	}

	fmt.Println(dm)
	fmt.Println(vm)
	fmt.Println(dim)
}

// Bread-first-search on map m from starting point sp to either end point ep or
// break on any tile that matches height eph (if not -1). Use cross function to
// validate crossing from current to new position. Returns distance in steps.
// NOTE: Probably would do better with a* or Dijkstra but I don't know them
// from top of my head...
func bfs(m [][]int, sp, ep XY, eph int, cross func(m [][]int, cp, np XY) bool) int {
	v := map[XY]bool{} // map of XY:visited
	di := map[XY]int{} // map of XY:distance

	var imgs []*image.Paletted

	v[sp] = true // set start as visited
	di[sp] = 0   // set distance to S as 0

	h := len(m)
	w := len(m[0])

	q := []XY{sp} // queue
	d := [2][4]int{
		[4]int{-1, 1, 0, 0},
		[4]int{0, 0, 1, -1},
	}

	for len(q) > 0 {
		// pop current position from q
		cp := q[0]
		q = q[1:]

		if vis {
			imgs = append(imgs, render(m, v, sp, ep, cp))
		}

		// if eph is > -1 break out and return
		if eph > -1 && m[cp[1]][cp[0]] == eph {
			if vis {
				saveGif("out2.gif", imgs)
			}
			return di[cp]
		}

		// test all directions
		for i := 0; i < 4; i++ {
			np := XY{cp[0] + d[0][i], cp[1] + d[1][i]}

			if np[0] < 0 || np[0] >= w || np[1] < 0 || np[1] >= h { // border check
				continue
			}
			if _, ok := v[np]; ok { // visited check
				continue
			}
			if cross(m, cp, np) { // tile cross check
				v[np] = true
				di[np] = di[cp] + 1
				q = append(q, np)
			}
			//debugDump(m, v, di)
		}
	}

	if vis {
		saveGif("out1.gif", imgs)
	}

	return di[ep]
}

func main() {

	m := [][]int{} // map data [Y][X] normalized to a=1 z=26
	sp := XY{}     // start position
	ep := XY{}     // end position

	vis = os.Getenv("DEBUG") != ""

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		b := s.Bytes() // read one line of map

		m = append(m, make([]int, 0))
		y := len(m) - 1
		for x, ch := range b {
			if ch == 'S' {
				ch = 'a' // S equals to a in terms of height
				sp = XY{x, y}
			}
			if ch == 'E' {
				ch = 'z' // E equals to z in terms of height
				ep = XY{x, y}
			}
			m[y] = append(m[y], int(ch)-'a')
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// go from S to E one height level at time
	// DAMN This also means that the elevation of the destination square can be
	// much lower than the elevation of your current square.
	log.Print(bfs(m, sp, ep, -1, func(m [][]int, cp, np XY) bool {
		return m[np[1]][np[0]]-m[cp[1]][cp[0]] <= 1
	}))
	// go from E to any a
	log.Print(bfs(m, ep, XY{}, 0 /* a */, func(m [][]int, cp, np XY) bool {
		return m[np[1]][np[0]]-m[cp[1]][cp[0]] >= -1
	}))
}
