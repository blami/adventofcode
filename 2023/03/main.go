// Find sum of all numbers adjacent (diagonally too) to non '.' characters and
// sum of all gear ratios (2 numbers adjacent to '*' multiplied).
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

func main() {
	sch := []string{}
	w:= 0
	h:= 0
	sum := 0
	sum2 := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		sch = append(sch, s.Text())
		// assuming all lines have same length
		if w == 0 { w = len(s.Text()) }
		h += 1
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			//fmt.Print(string(sch[y][x]))
			// if there is part start fishing numbers
			if sch[y][x] != '.' && !unicode.IsDigit(rune(sch[y][x])) {
				/*
				fmt.Printf("found %s at %d:%d\n", string(sch[y][x]), x, y)
				// debug printout
				for ddy := y - 1; ddy <= y + 1; ddy++ {
					if ddy < 0 || ddy > h - 1 { break }
					for i := 0; i < len(sch[ddy]); i++ {
						if i == x && ddy == y {
							fmt.Print("\033[31m")
						}
						fmt.Print(string(sch[ddy][i]))
						if i == x && ddy == y {
							fmt.Print("\033[0m")
						}
					}
					fmt.Println()
				}
				*/

				// try all 8 directions
				seen := [][]int{}
				pts := []int{}
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						// skip the part itself
						if dx == 0 && dy == 0 { continue }
						// bounds check
						if x + dx < 0 || y + dy < 0 || x + dx >= w || y + dy >= h { continue }
						// skip anything we saw already
						skip := false
						for _, s := range seen {
							if s[0] == x + dx && s[1] == y + dy {
								skip = true
								break
							}
						}
						if skip { continue }
						if unicode.IsDigit(rune(sch[y + dy][x + dx])) {
							// found number, seek to its begining and parse it
							// to the end
							n := 0
							ns := 0 // number start x
							for i := x+dx; unicode.IsDigit(rune(sch[y+dy][i])) && i >= 0; i-- {
								ns = i
								if i == 0 {
									break
								}
							}
							nx := ns
							for nx = ns; unicode.IsDigit(rune(sch[y+dy][nx])) && nx < w; nx++ {
								// add coords to seen
								seen = append(seen, []int{nx, y+dy})
								n = n * 10 + int(sch[y + dy][nx] - 48)
								if nx == w - 1 {
									break
								}
							}
							pts = append(pts, n)
							sum += n
							//fmt.Printf("  %s at %d:%d (start %d:%d, n: %d), sum: %d\n", string(sch[y+dy][x +dx]), dx, dy, ns, y+dy, n, sum)
						}
					}
				}
				// part 2
				if sch[y][x] == '*' && len(pts) == 2 {
					r := pts[0] * pts[1]
					//fmt.Printf("  gear r: %d\n", r)
					sum2 += r
				}
				//fmt.Println()
			}
		}
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
