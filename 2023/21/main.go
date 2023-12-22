// Find number of covered garden plots given the number of steps. With and
// without map wrapping at sides.
// NOTE: no test.txt for this one as it does not work with it. Test data in
// case of this puzzle was just to give partial hints to get implementation
// right.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type XY [2]int

func debug(m [][]rune, h []XY, i int) {
	fmt.Println(i)
	for y := range m {
		for x := range m[0] {
			p := string(m[y][x])
			for _, st := range h {
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

func tran(i, l int) int {
	if i < 0 {
		i += l
	}
	r := i % l
	if r < 0 {
		r = r + l
	}
	return r
}

// quadratic equation
func qe(v []int, x int) int {
	a := (v[0] - 2*v[1] + v[2]) / 2
	b := (-3*v[0] + 4*v[1] - v[2]) / 2
	c := v[0]
	fmt.Printf("solve: %d*x^2 + %d*x + %d (x=%d)\n", a, b, c, x)
	return a*x*x + b*x + c
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

	cov1 := 0

	h := []XY{{sx, sy}}
	for i := 0; i < 64; i++ {
		hn := []XY{}
		for _, c := range h {
			for _, d := range []XY{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				x := c[0] + d[0]
				y := c[1] + d[1]

				if x < 0 || x >= len(m[0]) || y < 0 || y >= len(m) {
					continue
				}
				if m[y][x] == '.' || m[y][x] == 'S' {
					has := false
					for i := range hn {
						if hn[i][0] == x && hn[i][1] == y {
							has = true
						}
					}
					if !has {
						hn = append(hn, XY{x, y})
					}
				}
			}
		}
		renderGif(m, h, hn, 2, [4]int{0, 1, 0, 1})
		h = hn
		// DIAMONDS DIAMONDS DIAMONDS
		//debug(m, h, i)
	}
	cov1 = len(h)
	saveGif("out.gif", true)

	// part 2
	// solve points for 0, 1 and 2... see q()
	coef := []int{} // coeficients for quadratic equation
	h = []XY{{sx, sy}}
	st := 26501365
	for i := 0; i < st; i++ {
		hn := []XY{}
		hncach := map[XY]bool{} // slice finds cost, cache in map
		for _, c := range h {
			for _, d := range []XY{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				x := c[0] + d[0]
				y := c[1] + d[1]
				wx, wy := tran(x, len(m[0])), tran(y, len(m))
				if m[wy][wx] == '.' || m[wy][wx] == 'S' {
					xy := XY{x, y}
					if _, has := hncach[xy]; !has {
						hn = append(hn, xy)
						hncach[xy] = true
					}
				}
			}
		}
		//h = hn AAARGH this took me good hour to comment and move down!
		// the period of repeat from 65 is width of map (131 for my input)
		// so a,b and c are at i=65, i=65+131 and i=65+2*131
		if i%len(m[0]) == st%len(m[0]) {
			fmt.Println("i", i, "len(h)", len(h))
			coef = append(coef, len(h))
			saveGif("out2.gif", false)
		}
		if len(coef) == 3 {
			break
		}
		renderGif(m, h, hn, 2, [4]int{1, 2, 1, 2})
		h = hn
		//debug(m, h, i)
	}
	cov2 := qe(coef, st/len(m[0]))
	saveGif("out2.gif", true)

	fmt.Println(cov1)
	fmt.Println(cov2)
}
