// Find board load after tilting north once and after 1000000000 all direction
// tilts.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func debug(m [][]rune) {
	for i := range m {
		fmt.Println(string(m[i]))
	}
	fmt.Println()
}

// Lol this is probably the ugliest thing I wrote this year, but what has to be
// done that has to be done. d=0-north, 1-west, 2-south, 3-east
func tilt(m [][]rune, d int) {

	if d == 0 || d == 1 {
		for y := range m {
			for x := range m[0] {
				if m[y][x] == 'O' {
					if d == 0 {
						for ry := y; ry > 0; ry-- {
							if m[ry-1][x] == '.' {
								m[ry-1][x] = 'O'
								m[ry][x] = '.'
							} else {
								break
							}
						}
					} else if d == 1 {
						for rx := x; rx > 0; rx-- {
							if m[y][rx-1] == '.' {
								m[y][rx-1] = 'O'
								m[y][rx] = '.'
							} else {
								break
							}
						}
					}
				}
			}
		}
	} else {
		for y := len(m) - 1; y >= 0; y-- {
			for x := len(m[0]) - 1; x >= 0; x-- {
				if m[y][x] == 'O' {
					if d == 2 {
						for ry := y; ry < len(m)-1; ry++ {
							if m[ry+1][x] == '.' {
								m[ry+1][x] = 'O'
								m[ry][x] = '.'
							} else {
								break
							}
						}
					} else if d == 3 {
						for rx := x; rx < len(m[0])-1; rx++ {
							if m[y][rx+1] == '.' {
								m[y][rx+1] = 'O'
								m[y][rx] = '.'
							} else {
								break
							}
						}
					}
				}
			}
		}
	}
}

func cp(m [][]rune) [][]rune {
	c := make([][]rune, len(m))
	for i := range m {
		c[i] = make([]rune, len(m[0]))
		copy(c[i], m[i])
	}
	return c
}

func cmp(a [][]rune, b [][]rune) bool {
	for y := range a {
		for x := range a[0] {
			if a[y][x] != b[y][x] {
				return false
			}
		}
	}
	return true
}

func load(m [][]rune) int {
	sum := 0
	load := len(m)
	for y := range m {
		sum += strings.Count(string(m[y]), "O") * (load - y)
	}
	return sum
}

func main() {
	m := [][]rune{}
	sum := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		m = append(m, []rune(l))
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	cycles := 1000000000 // part2

	// for part 2 *I could smell loop miles away* so I started caching and
	// figured out there's indeed a loop. So to figure out load we need to
	// remove items before loop (cycles - loop_sti), using length of the loop
	// we calculate remainding items within next round of loop (cycles -
	// loop_sti) % loop_l. Now we use index in cache where loop started
	// (loop_stk) and add on top of it which ultimately gives us cached item
	// that will be cycles-th and just call load() on it.

	cach := [][][]rune{}
	loop_sti := -1 // index where loop starts
	loop_stk := -1 // loop start index in cache
	loop_l := 0    // length of the loop
	done := false
	for i := 0; i < cycles; i++ {
		//fmt.Println("i=", i)
		for j := 0; j < 4; j++ {
			tilt(m, j)
			if i == 0 && j == 0 {
				sum = load(m) // part 1
			}
		}
		s := true
		for k := range cach {
			if cmp(cach[k], m) {
				//fmt.Println("i", i, "loop_st", loop_sti, "cachel", len(cach), "same at ", k, load(m))
				if loop_stk == k {
					done = true
					break
				}
				if loop_sti == -1 {
					loop_sti = i
					loop_stk = k
				}
				if loop_sti != -1 {
					loop_l++
				}
				s = false
				//m
				break
			}
		}
		if s {
			cach = append(cach, cp(m))
		}
		if done {
			break
		}
	}
	i := (cycles - loop_sti) % loop_l
	sum2 := load(cach[loop_stk+i])

	fmt.Println(sum)
	fmt.Println(sum2)
}
