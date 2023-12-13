// Find reflection line in given pattern either without smudge or with one
// smudge.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func findm(m [][]rune, smudge bool) ([2]int, [2]int) {
	// go column by column and try going both sides
	f := [2]int{-1, -1}
	s := [2]int{-1, -1}
	for x := 0; x < len(m[0]); x++ {
		d := 0
		ss := [][2]int{}
		for x-d >= 0 && x+d+1 < len(m[0]) {
			for y := 0; y < len(m); y++ {
				if m[y][x-d] != m[y][x+d+1] {
					ss = append(ss, [2]int{x - d, y})
					if !smudge {
						break
					}
				}
			}
			if (!smudge && len(ss) != 0) || (smudge && len(ss) > 1) {
				break
			}
			if ((!smudge && len(ss) == 0) || (smudge && len(ss) == 1)) && (x-d == 0 || x+d+1 == len(m[0])-1) {
				f[0] = x
				if smudge {
					s = ss[0]
				}
				break
			}
			d++
		}
		if f[0] != -1 {
			break
		}
	}

	// anybody like spaghetti?
	if f[0] == -1 {
		for y := 0; y < len(m); y++ {
			d := 0
			ss := [][2]int{}
			for y-d >= 0 && y+d+1 < len(m) {
				for x := 0; x < len(m[0]); x++ {
					if m[y-d][x] != m[y+d+1][x] {
						ss = append(ss, [2]int{x, y - d})
						if !smudge {
							break
						}
					}
				}
				if (!smudge && len(ss) != 0) || (smudge && len(ss) > 1) {
					break
				}
				if ((!smudge && len(ss) == 0) || (smudge && len(ss) == 1)) && (y-d == 0 || y+d+1 == len(m)-1) {
					f[1] = y
					if smudge {
						s = ss[0]
					}
					break
				}
				d++
			}
			if f[1] != -1 {
				break
			}
		}
	}

	return f, s
}

func debug(m [][]rune, f [2]int, s [2]int) {
	for y := -2; y < len(m); y++ {
		if y < 0 {
			fmt.Printf("   ")
		} else {
			if f[1] == y {
				fmt.Printf("v %d", y%10)
			} else if f[1]+1 == y && f[1] != -1 {
				fmt.Printf("^ %d", y%10)
			} else {
				fmt.Printf("  %d", y%10)
			}
		}
		for x := range m[0] {
			if y == -1 {
				fmt.Printf("%d", x%10)
			} else if y == -2 {
				if f[0] != -1 {
					if f[0] == x {
						fmt.Printf(">")
					} else if f[0]+1 == x {
						fmt.Printf("<")
					} else {
						fmt.Printf(" ")
					}
				}
			} else {
				if s[0] == x && s[1] == y {
					if m[y][x] == '.' {
						fmt.Printf("\033[31m#\033[0m")
					} else {
						fmt.Printf("\033[31m.\033[0m")
					}
				} else {
					fmt.Printf("%s", string(m[y][x]))
				}
			}
		}
		fmt.Println()
	}
}

func addsum(m [][]rune, smudge bool) int {
	f, _ /*s*/ := findm(m, smudge)
	//debug(m, f, s)
	if f[1] == -1 && f[0] != -1 {
		//fmt.Println("horz match", smudge, "at", f[0], f[0]+1)
		return f[0] + 1
	} else if f[0] == -1 && f[1] != -1 {
		//fmt.Println("vert match", smudge, "at", f[1], f[1]+1)
		return (f[1] + 1) * 100
	} else {
		panic("is this case?")
	}
}

func main() {
	m := [][]rune{}

	sum := 0
	sum2 := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			sum += addsum(m, false)
			sum2 += addsum(m, true)
			m = [][]rune{}
			continue
		}
		m = append(m, []rune(l))
	}
	sum += addsum(m, false)
	sum2 += addsum(m, true)
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
