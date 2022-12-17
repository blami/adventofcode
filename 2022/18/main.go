// Read input from stdin. Explore lava cube space and calculate total surface
// and outer facing surface of droplet.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// XYZ alias.
type XYZ [3]int

// Return all sides of cube.
func sides(c XYZ) [6]XYZ {
	return [6]XYZ{
		XYZ{c[0] + 1, c[1], c[2]},
		XYZ{c[0] - 1, c[1], c[2]},
		XYZ{c[0], c[1] + 1, c[2]},
		XYZ{c[0], c[1] - 1, c[2]},
		XYZ{c[0], c[1], c[2] + 1},
		XYZ{c[0], c[1], c[2] - 1},
	}
}

// Is given side in cubes slice?
func in(cs *[]XYZ, s XYZ) bool {
	for _, c := range *cs {
		if c == s {
			return true
		}
	}
	return false
}

// Is the side within droplet boundary?
func bound(s XYZ, min, max int) bool {
	for _, a := range s {
		if a < min-1 || a > max+1 {
			return false
		}
	}
	return true
}

func main() {
	cs := []XYZ{} // cubes

	suc := 0 // connected surface
	max := 0
	min := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		c := XYZ{}
		for i, v := range strings.Split(s.Text(), ",") {
			c[i], _ = strconv.Atoi(v)
			if c[i] > max {
				max = c[i]
			}
			if c[i] < min {
				min = c[i]
			}
		}

		// place the cube, if any of existing cubes neighbors with cube remove
		// 2 sides from total surface for each neighboring side
		for _, ec := range cs {
			if c[1] == ec[1] && c[2] == ec[2] && (c[0]+1 == ec[0] || c[0]-1 == ec[0]) {
				suc += 2
			}
			if c[0] == ec[0] && c[2] == ec[2] && (c[1]+1 == ec[1] || c[1]-1 == ec[1]) {
				suc += 2
			}
			if c[0] == ec[0] && c[1] == ec[1] && (c[2]+1 == ec[2] || c[2]-1 == ec[2]) {
				suc += 2
			}
		}
		cs = append(cs, c)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// bfs visit the droplet from outside (considering min,min,min as safe
	// outside for now). outside surface is then sum of visited sides that
	// belong to cubes.
	v := []XYZ{}
	q := []XYZ{XYZ{min, min, min}} // queue
	suo := 0                       // outside surface
	for len(q) > 0 {
		c := q[len(q)-1]
		q = q[:len(q)-1] // pop (BUG: made stupid mistake to pop from begining)
		if !in(&v, c) {
			v = append(v, c)
		} // add c to seen
		for _, s := range sides(c) {
			if !in(&cs, s) && !in(&v, s) && bound(s, min, max) {
				q = append(q, s)
			}
		}
	}
	for _, c := range cs {
		for _, s := range sides(c) {
			if in(&v, s) {
				suo++
			} // sum visited cube sides
		}
	}

	log.Print((len(cs) * 6) - suc) // part 1
	log.Print(suo)                 // part 2
}
