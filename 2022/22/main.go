package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"regexp"
	"fmt"
)

type XY [2]int

var D [4]XY = [4]XY {
	XY{1, 0}, // 0 right
	XY{0, 1}, // 1 down
	XY{-1, 0}, // 2 left
	XY{0, -1}, // 3 up
}

var F map[int]byte = map[int]byte{
	0: '>',
	1: 'v',
	2: '<',
	3: '^',
}


func turn(cd int, d string) int {
	if d == "L" { cd -=1 } else { cd+=1 }
	if cd < 0 { cd = len(D) - 1 } else if cd > len(D) - 1 { cd = 0 } // always +-1
	return cd
}

func render(m [][]byte, cp XY, f int) {
	for y := range m {
		l := ""
		for x, ch := range m[y] {
			if cp[0] == x && cp[1] == y { ch = F[f] }
			l += string(ch)
		}
		fmt.Printf("%3d%s\n", y, l)
	}
	log.Printf("%c %v", F[f], cp)
}

func main() {
	m := [][]byte{}
	p := []string{}
	cp := XY{0,0}
	f := 0 // facing right


	re := regexp.MustCompile(`[0-9]+|[LR]`)
	s := bufio.NewScanner(os.Stdin)
	inm := true
	for s.Scan() {
		l := s.Text()
		if l == "" {
			inm = false
			continue
		}
		if inm {
			m = append(m, []byte(l))
			// mark start position
			if len(m) == 1 {
				for x, ch := range m[0] {
					if ch == '.' {
						cp[0] = x
						break
					}
				}
			}
		} else {
			p = append(p, re.FindAllString(l, -1)...)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// move step by step
	for _, s := range p {
		//log.Print(s)
		if s == "L" || s == "R" {
			f = turn(f, s)
			m[cp[1]][cp[0]] = F[f]
		} else {
			n, _ := strconv.Atoi(s)
			for n > 0 {
				np := XY{cp[0] + D[f][0], cp[1] + D[f][1]} // new position
				//log.Print(np)
				// handle wraps 
				if D[f][1] != 0 { // Y
					//log.Print("Y")
					if  np[1] < 0 || np[1] > len(m) - 1 || len(m[np[1]]) - 1 < np[0] || m[np[1]][np[0]] == ' ' {
						// wrapping from top to down
						if D[f][1] == -1 {
							// BUG: what if there's nothing free here
							ny := 0
							for ny = len(m)-1; ny >= 0; ny-- {
								if len(m[ny]) - 1 < np[0] { continue } // line not long enough
								//log.Print(ny, " ", len(m[ny])-1, " ", np[1])
								if m[ny][np[0]] != ' ' {
									// found = true
									break
								}
							}
							// transition func
							np[1] = ny
						} else {
						// wrapping from down to up
							ny := 0
							for ny = 0; ny < len(m); ny++ {
								if len(m[ny]) - 1 < np[0] { continue } // line not long enough
								if m[ny][np[0]] != ' ' {
									// found = true
									break
								}
							}
							// transition func
							np[1] = ny
						}
					}
				} else if D[f][0] != 0 { // X
					if  np[0] < 0 || np[0] > len(m[np[1]]) - 1 || m[np[1]][np[0]] == ' '{
						// wrapping from left to right
						if D[f][0] == -1 {
							// BUG: what if there's nothing free here
							nx := 0
							for nx = len(m[np[1]])-1; nx >= 0; nx-- {
								if m[np[1]][nx] != ' ' {
									// found = true
									break
								}
							}
							// transition func
							np[0] = nx
						} else {
						// wrapping from right to left
							nx := 0
							for nx = 0; nx < len(m[np[1]]); nx++ {
								if m[np[1]][nx] != ' ' {
									// found = true
									break
								}
							}
							// transition func
							np[0] = nx
						}
					}
				}

				// handle wall hits
				if m[np[1]][np[0]] != '#' {
					cp = np
					m[cp[1]][cp[0]] = F[f]
				} else {
					break // stop here...
				}
				n--
			}
		}
		//render(m, cp, f)
	}

	render(m, cp, f)
	//log.Print(cp[1]+1, " ", cp[0]+1, " ", f)
	log.Print(((cp[1]+1) * 1000) + ((cp[0]+1) * 4) + f)
	//log.Print(m,p,cp)

}
