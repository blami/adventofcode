package main

import (
	"bufio"
	"log"
	"os"
)

// XY type alias.
type XY [2]int

// A map of current elf positions.
type Map struct {
	t, l, b, r int       // map boundaries (may go of initially as we start at 0)
	elves      map[XY]XY // map of elves at their current positions to their proposal
}

// Render map. For debugging purposes.
func (m *Map) render() {
	for y := m.t; y <= m.b; y++ {
		l := ""
		for x := m.l; x <= m.r; x++ {
			ch := '.'
			if _, ok := m.elves[XY{x, y}]; ok {
				ch = '#'
			}
			l += string(ch)
		}
		log.Printf("%3d %s", y, l)
	}
	log.Printf("    %d", m.l)
}

// Shift slice of 3 XYs to left r times.
func shl(sl [][3]XY, r int) [][3]XY {
	r = r % len(sl)
	sl = append(sl[r:], sl[0:r]...)
	return sl
}

// Propose move from p. Return new position as a resolution. In case there's no
// movement at all, just return p.
func (m *Map) propose(p XY, r int) XY {
	dd := [][3]XY{ // direction differences from p
		[3]XY{XY{0, -1}, XY{1, -1}, XY{-1, -1}}, // N, NE, NW
		[3]XY{XY{0, 1}, XY{1, 1}, XY{-1, 1}},    // S, SE, SW
		[3]XY{XY{-1, 0}, XY{-1, -1}, XY{-1, 1}}, // W, NW, SW
		[3]XY{XY{1, 0}, XY{1, -1}, XY{1, 1}},    // E, NE, SE
	}
	// BUG: read three times, code once; direction lookup changes with each
	// round (shifts left)
	dd = shl(dd, r)

	// BUG: read twice, code once; elf first looks around
	lone := true
	for _, lo := range []XY{XY{0, -1}, XY{1, -1}, XY{1, 0}, XY{1, 1}, XY{0, 1}, XY{-1, 1}, XY{-1, 0}, XY{-1, -1}} {
		//log.Printf(" lonecheck elf %v %v lo=%v", p, XY{p[0] + lo[0], p[1] + lo[1]}, lo)
		if _, ok := m.elves[XY{p[0] + lo[0], p[1] + lo[1]}]; ok {
			//log.Printf("elf %v is not alone", p)
			lone = false
			break
		}
	}
	//if lone { log.Printf("elf %v is alone", p) }
	if !lone {
	M:
		for d := range dd {
			for _, dxy := range dd[d] {
				if _, ok := m.elves[XY{p[0] + dxy[0], p[1] + dxy[1]}]; ok {
					//log.Printf("elf %v can't go %d", p, d)
					continue M // nop, try another direction
				}
			}
			//log.Printf("elf %v proposing %v", p, XY{p[0] + dd[d][0][0], p[1] + dd[d][0][1]})
			return XY{p[0] + dd[d][0][0], p[1] + dd[d][0][1]}
		}
	}
	//log.Printf("elf %v not proposing", p)
	return p // not moving
}

// Is there more than one proposal to go to XY?
func (m *Map) only(p XY) bool {
	cnt := 0
	for _, pp := range m.elves {
		if pp == p {
			cnt++
			if cnt > 1 {
				return false
			}
		}
	}
	// cnt == 0 doushimasho?
	return true
}

// Do one round.
func (m *Map) round(r int) int {
	for xy := range m.elves {
		m.elves[xy] = m.propose(xy, r)
	}
	nelves := make(map[XY]XY) // need new map to place elves at new positions
	m.l = 0                   // squeeze the map so round expands it to smallest possible rect
	m.r = 0
	m.t = 0
	m.b = 0
	mvs := 0 // count moves in round
	for xy := range m.elves {
		nxy := xy
		if m.only(m.elves[xy]) && xy != m.elves[xy] {
			//log.Printf("elf %v moving to %v", xy, m.elves[xy])
			nxy = m.elves[xy] // move to proposed location
			mvs++
		} else {
			//log.Printf("elf %v not moving", xy)
			nxy = xy
		}
		nelves[nxy] = nxy
		if nxy[0] < m.l { // stretch the map boundaries as we go
			m.l = nxy[0]
		}
		if nxy[0] > m.r {
			m.r = nxy[0]
		}
		if nxy[1] < m.t {
			m.t = nxy[1]
		}
		if nxy[1] > m.b {
			m.b = nxy[1]
		}
	}
	m.elves = nelves
	return mvs
}

// Distance between a and b.
func dist(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

// Calculate number of free tiles on the map.
func (m *Map) free() int {
	return ((dist(m.r, m.l) + 1) * (dist(m.b, m.t) + 1)) - len(m.elves)
}

func main() {

	m := Map{}
	m.elves = make(map[XY]XY) // constructors are for kids...
	m.t, m.l = 10, 10         // just to make things nicer in begining

	s := bufio.NewScanner(os.Stdin)
	y := 0
	for s.Scan() {
		l := s.Text()
		for x, ch := range l {
			if ch == '#' {
				m.elves[XY{x, y}] = XY{x, y}
				// ugly but works
				if x < m.l {
					m.l = x
				}
				if x > m.r {
					m.r = x
				}
				if y < m.t {
					m.t = y
				}
				if y > m.b {
					m.b = y
				}
			}
		}
		y++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	r := 0 // round counter
	for {
		mv := m.round(r)
		if r == 9 {
			log.Print(m.free())
		}
		if r > 9 && mv == 0 { // at least ten rounds to solve part 1
			log.Print(r + 1)
			break
		}
		/*
			if r%100 == 0 {
				log.Printf("r=%d not stable", r)
			}
			log.Printf("round %d", r + 1)
			m.render()
		*/
		r++
	}
}
