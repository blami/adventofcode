package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type XY [2]int

func debug(m [][]rune, h [][]XY, i int) {
	fmt.Println(i)
	for y := range(m) {
		for x := range(m[0]) {
			p := string(m[y][x])
			for _, st := range h[len(h)-1] {
				if x == st[0] && y == st[1] {
					p = "O"
				}
			}
			fmt.Print(p)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	m := [][]rune{}
	sy, sx := 0, -1
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		m = append(m, []rune(l))
		// find S
		if sx == -1 {
			sx = strings.Index(l, "S")
		}
		if sx == -1 {
			sy += 1
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	
	h := [][]XY{{{sx,sy}}}
	for i := 0; i < 64; i++ {
		// go over all last steps and record new steps
		ha := []XY{}
		for _, c := range h[len(h)-1] {
			for _, d := range []XY{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				x := c[0] + d[0]
				y := c[1] + d[1]

				if x < 0 || x >= len(m[0]) || y < 0 || y >= len(m) {
					continue
				}
				if m[y][x] == '.' {
					has := false
					for i := range(ha) {
						if ha[i][0] == x && ha[i][1] == y {
							has = true
						}
					}
					if !has {
						ha = append(ha, XY{x, y})
					}
				}
			}
		}
		h = append(h, ha)

		//debug(m, h, i)
	}

	fmt.Println(len(h[len(h)-1])+1) // +1 for initial S
}
