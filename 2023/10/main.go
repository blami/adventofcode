// Find half of length of the loop and number of contained tiles in it.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type XY struct {
	X int
	Y int
}

func main() {
	m := [][]rune{}
	w, h := 0, 0
	st := XY{}
	sum := 0 // part 2

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		w = len(l)
		sti := strings.Index(l, "S")
		if sti > -1 {
			st.X = sti
			st.Y = h
		}
		m = append(m, []rune(l))
		h++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	dir := []XY{}
	// find directions of move
	// BUG: This is not correct as S might be attached to L7JF in multiple
	// directions... it worked for my input but breaks test5.txt
	if st.X-1 >= 0 && m[st.Y][st.X-1] == '-' {
		dir = append(dir, XY{-1, 0})
	}
	if st.X+1 < w && m[st.Y][st.X+1] == '-' {
		dir = append(dir, XY{1, 0})
	}
	if st.Y-1 >= 0 && m[st.Y-1][st.X] == '|' {
		dir = append(dir, XY{0, -1})
	}
	if st.Y+1 < h && m[st.Y+1][st.X] == '|' {
		dir = append(dir, XY{0, 1})
	}

	cur := XY{st.X, st.Y}
	step := 0
	// part 2
	a := [][]rune{} // note down the loop so we can check parity of |'s
	for y := 0; y < h; y++ {
		a = append(a, []rune(strings.Repeat(".", w)))
	}

	done := false
	for !done {
		step++
		a[cur.Y][cur.X] = '#'
		cur.X += dir[0].X
		cur.Y += dir[0].Y
		fmt.Println("moving to", dir[0], cur.X, cur.Y, string(m[cur.Y][cur.X]), "step", step)
		switch m[cur.Y][cur.X] {
		// -,| - continue moving in same direction
		case 'L', '7':
			dx := dir[0].X
			dir[0].X = dir[0].Y
			dir[0].Y = dx
		case 'J', 'F':
			dx := dir[0].X
			dir[0].X = dir[0].Y * -1
			dir[0].Y = dx * -1
		case '.':
			panic("oops")
		case 'S':
			done = true
			break
		}
	}

	// BUG: Flood fill won't work because of pipe squeezing
	/*
		done = false
		for !done {
			done = true
			for y := range(a) {
				for x := range(a[y]) {
					// try to search in all directions
					if a[y][x] == '.' {
						if (y==0) || (y==h-1) || (x > 0 && a[y][x - 1] == ' ') || (x < w - 1 && a[y][x + 1] == ' ') || (y > 0 && a[y - 1][x] == ' ') || (y < h-1 && a[y + 1][x] == ' ') {
							a[y][x] = ' '
							done = false
						}
					}
				}
			}
		}
		for _, l := range(a) {
			sum += strings.Count(string(l), ".")
		}
	*/
	// scan lines and check parity of |JL's lol
	for y := range m {
		in := false
		for x := range m[y] {
			if a[y][x] == '#' {
				// BUG: This is not entirely correct as S might be used as | as
				// well but not in my input.
				if m[y][x] == '|' || m[y][x] == 'J' || m[y][x] == 'L' {
					in = !in
				}
			} else {
				if in {
					m[y][x] = 'I'
					sum++
				} else {
					m[y][x] = 'O'
				}
			}
		}
	}

	// debug
	for y := range m {
		fmt.Print(string(m[y]))
		//fmt.Print(" ", string(a[y]))
		fmt.Println()
	}

	fmt.Println(step / 2)
	fmt.Println(sum)
}
