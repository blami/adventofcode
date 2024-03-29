// Find half of length of the loop and number of contained tiles in it.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"image"
)

type XY struct {
	X int
	Y int
}

func main() {
	m := [][]rune{}
	w, h := 0, 0
	st := XY{}
	// part 2
	sum := 0
	// visualization makes run slow
	vis := os.Getenv("DEBUG") != ""
	vs := []int{1, 1}
	vss := 5
	var imgs []*image.Paletted

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
	// vis
	if vis {
		if w > 100 { // BIG input
			vs = []int{50, 5}
			vss = 2
		}
	}

	dir := []XY{}
	// find directions of move
	if st.X-1 >= 0 && (m[st.Y][st.X-1] == '-' || m[st.Y][st.X-1] == 'L') {
		dir = append(dir, XY{-1, 0})
	}
	if st.X+1 < w && (m[st.Y][st.X+1] == '-' || m[st.Y][st.X+1] == '7') {
		dir = append(dir, XY{1, 0})
	}
	if st.Y-1 >= 0 && (m[st.Y-1][st.X] == '|' || m[st.Y-1][st.X] == 'F') {
		dir = append(dir, XY{0, -1})
	}
	if st.Y+1 < h && (m[st.Y+1][st.X] == '|' || m[st.Y+1][st.X] == 'J') {
		dir = append(dir, XY{0, 1})
	}
	// part 2
	// fix for S being border too
	sb := true
	for _, d := range dir {
		if d.X != 0 {
			sb = false
		}
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
		// vis
		if vis && step%vs[0] == 0 {
			imgs = append(imgs, render(m, a, st, cur, XY{-1, -1}, vss))
		}

		a[cur.Y][cur.X] = '#'
		cur.X += dir[0].X
		cur.Y += dir[0].Y
		//fmt.Println("moving to", dir[0], cur.X, cur.Y, string(m[cur.Y][cur.X]), "step", step)
		switch m[cur.Y][cur.X] {
		// -,| - continue moving in same direction
		case 'L', '7':
			dir[0].X, dir[0].Y = dir[0].Y, dir[0].X
		case 'J', 'F':
			dir[0].X, dir[0].Y = dir[0].Y*-1, dir[0].X*-1
		case '.':
			panic("oops")
		case 'S':
			done = true
			break
		}
	}

	// scan lines and check parity of |JL's lol
	for y := range m {
		in := false
		for x := range m[y] {
			if a[y][x] == '#' {
				if m[y][x] == '|' || m[y][x] == 'J' || m[y][x] == 'L' || (sb && m[y][x] == 'S') {
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

			// vis
			if vis && (y%vs[1] == 0 || y == h-1) && x == w-1 {
				imgs = append(imgs, render(m, a, st, XY{-1, -1}, XY{x, y}, vss))
			}
		}
	}

	// debug
	/*
		for y := range m {
			fmt.Print(string(m[y]))
			//fmt.Print(" ", string(a[y]))
			fmt.Println()
		}
	*/
	if vis {
		saveGif("out.gif", imgs)
	}

	fmt.Println(step / 2)
	fmt.Println(sum)
}
