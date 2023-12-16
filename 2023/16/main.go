// Find number of lightbeam energized tiles, find maximum possible energized
// tiles possible when starting from any edge tile.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type XY [2]int

type Beam struct {
	XY   XY
	Dir  XY
	Done bool
}

type Tile struct {
	C rune
	E bool
}

func debug(m [][]Tile, b []Beam) {
	for y := range m {
		for x := range m[y] {
			c := m[y][x].C
			if m[y][x].E {
				c = '#'
			}

			for i := range b {
				if b[i].Done {
					continue
				}
				if b[i].XY[0] == x && b[i].XY[1] == y {
					switch b[i].Dir {
					case XY{-1, 0}:
						c = '<'
					case XY{1, 0}:
						c = '>'
					case XY{0, -1}:
						c = '^'
					case XY{0, 1}:
						c = 'v'
					}
				}
			}
			fmt.Print(string(c))
		}
		fmt.Print(" ")
		for x := range m {
			fmt.Print(string(m[y][x].C))
		}
		fmt.Println()
	}
	fmt.Println()
}

func n(m [][]Tile) int {
	r := 0
	for y := range m {
		for x := range m[y] {
			if m[y][x].E {
				r++
			}
		}
	}
	return r
}

type XYXY struct {
	SP  XY
	Dir XY
}

func main() {
	m := [][]Tile{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ml := []Tile{}
		for _, c := range s.Text() {
			ml = append(ml, Tile{c, false})
		}
		m = append(m, ml)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// race the beam
	n2 := 0
	n1 := 0
	// generate all starting points
	bs := []Beam{}

	for x := range m[0] {
		bs = append(bs,
			// top
			Beam{XY{x, 0}, XY{1, 0}, false},
			Beam{XY{x, 0}, XY{-1, 0}, false},
			Beam{XY{x, 0}, XY{0, 1}, false},
			// bottom
			Beam{XY{x, len(m) - 1}, XY{1, 0}, false},
			Beam{XY{x, len(m) - 1}, XY{-1, 0}, false},
			Beam{XY{x, len(m) - 1}, XY{0, -1}, false},
		)
	}
	for y := range m {
		// already covered by x
		if y == 0 || y == len(m)-1 {
			continue
		}
		bs = append(bs,
			// left
			Beam{XY{0, y}, XY{0, -1}, false},
			Beam{XY{0, y}, XY{0, 1}, false},
			Beam{XY{0, y}, XY{1, 0}, false},
			// right
			Beam{XY{len(m[0]) - 1, y}, XY{0, -1}, false},
			Beam{XY{len(m[0]) - 1, y}, XY{0, 1}, false},
			Beam{XY{len(m[0]) - 1, y}, XY{-1, 0}, false},
		)
	}

	for cf, b0 := range bs {
		// clean up
		for y := range m {
			for x := range m[0] {
				m[y][x].E = false
			}
		}

		b := []Beam{b0}
		// remember all tuples of starting XY + direction to handle loops
		sp := map[[2]XY]bool{{b0.XY, b0.Dir}: true}
		done := false
		for !done {
			// go over all beams
			for i := range b {
				/*
					if !b[i].Done {
						fmt.Println("beam", i, "xy", b[i].XY, "dir", b[i].Dir)
					}
				*/
				x, y := b[i].XY[0], b[i].XY[1]
				if x < 0 || x >= len(m[0]) || y < 0 || y >= len(m) {
					b[i].Done = true
				}
				if b[i].Done {
					continue
				}

				// ENERGIZE!
				m[y][x].E = true

				switch m[y][x].C {
				case '.':
					// just continue
				case '/':
					switch b[i].Dir {
					case XY{1, 0}: // from left
						b[i].Dir = XY{0, -1}
					case XY{-1, 0}: // from rignt
						b[i].Dir = XY{0, 1}
					case XY{0, 1}: // from top
						b[i].Dir = XY{-1, 0}
					case XY{0, -1}: // from bottom
						b[i].Dir = XY{1, 0}
					}
				case '\\':
					switch b[i].Dir {
					case XY{1, 0}: // from left
						b[i].Dir = XY{0, 1}
					case XY{-1, 0}: // from right
						b[i].Dir = XY{0, -1}
					case XY{0, 1}: // from top
						b[i].Dir = XY{1, 0}
					case XY{0, -1}: // from bottom
						b[i].Dir = XY{-1, 0}
					}
				case '|':
					// if moving from top, bottom just continue, otherwise split
					if b[i].Dir[0] != 0 {
						b[i].Done = true // finish this one and append two new ones
						if _, ok := sp[[2]XY{{x, y - 1}, {0, -1}}]; !ok {
							b = append(b, Beam{XY{x, y - 1}, XY{0, -1}, false}) // up
							sp[[2]XY{{x, y - 1}, {0, -1}}] = true
						}
						if _, ok := sp[[2]XY{{x, y + 1}, {0, 1}}]; !ok {
							b = append(b, Beam{XY{x, y + 1}, XY{0, 1}, false}) // down
							sp[[2]XY{{x, y + 1}, {0, 1}}] = true
						}
					}
				// copypasta anyone?
				case '-':
					// if moving from left, right just continue, otherwise split
					if b[i].Dir[1] != 0 {
						b[i].Done = true // finish this one and append two new ones
						if _, ok := sp[[2]XY{{x - 1, y}, {-1, 0}}]; !ok {
							b = append(b, Beam{XY{x - 1, y}, XY{-1, 0}, false}) // left
							sp[[2]XY{{x - 1, y}, {-1, 0}}] = true
						}
						if _, ok := sp[[2]XY{{x + 1, y}, {1, 0}}]; !ok {
							b = append(b, Beam{XY{x + 1, y}, XY{1, 0}, false}) // right
							sp[[2]XY{{x + 1, y}, {1, 0}}] = true
						}
					}
				}
				// if not done move
				if b[i].Done {
					continue
				}
				b[i].XY[0] += b[i].Dir[0]
				b[i].XY[1] += b[i].Dir[1]
			}
			// do we have any beams to move?
			done = true
			for i := range b {
				if !b[i].Done {
					done = false
				}
			}

		}

		nn := n(m)
		/*
			fmt.Println("cfg", cf, "n", nn)
			debug(m, []Beam{})
			fmt.Println()
		*/
		if cf == 0 {
			n1 = nn
		}
		if nn > n2 {
			n2 = nn
		}
	}

	fmt.Println(n1)
	fmt.Println(n2)
}
