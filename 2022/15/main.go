// Read input from stdin and given Y coordinate as 1st command line argument
// find number of positions where can't be beacon on that row; given MAX value
// as 2nd command line argument find a position of hidden beacon and calculate
// its distress frequency.

// TODO: Did this as quick hack in between meetings. Get back to this and
// optimize it; runtime around 10s is terrible.

package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"

	"runtime/pprof"
)

// XY alias.
type XY [2]int

// An interval.
type I [2]int

// Manhattan distance between X and Y.
func madi(a, b XY) int {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx + dy
}

// Add interval to merged intervals slice...
// TODO: this is ugly and slow; optimize this
func add(ivs []I, iv I) []I {
	//log.Print(ivs, " + ", iv)

	out := []I{}
	ivs = append(ivs, iv)

	if len(ivs) == 1 {
		return ivs
	}
	sort.Slice(ivs, func(i, j int) bool {
		if ivs[i][0] == ivs[j][0] {
			return ivs[i][1] < ivs[j][1]
		}
		return ivs[i][0] < ivs[j][0]
	})
	out = append(out, ivs[0])
	for i := 1; i < len(ivs); i++ {
		e := out[len(out)-1][1]
		if ivs[i][0]-1 <= e {
			out[len(out)-1][1] = func(a, b int) int {
				if a > b {
					return a
				}
				return b
			}(e, ivs[i][1])
		} else {
			out = append(out, I{ivs[i][0], ivs[i][1]})
		}
	}

	//log.Print("  ", out)
	return out
}

func main() {

	if len(os.Args) < 3 {
		log.Fatal("arguments: <Y> <MAX>")
	}
	argy, _ := strconv.Atoi(os.Args[1])
	argmax, _ := strconv.Atoi(os.Args[2])

	f, _ := os.Create("prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	se := [][2]XY{} // 0-sensor, 1-closest beacon

	re := regexp.MustCompile(`Sensor at x=(?P<sx>\d+), y=(?P<sy>\d+): closest beacon is at x=(?P<bx>-?\d+), y=(?P<by>-?\d+)`)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		mtch := re.FindStringSubmatch(s.Text())
		sx, _ := strconv.Atoi(mtch[1])
		sy, _ := strconv.Atoi(mtch[2])
		bx, _ := strconv.Atoi(mtch[3])
		by, _ := strconv.Atoi(mtch[4])
		se = append(se, [2]XY{XY{sx, sy}, XY{bx, by}})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	bl := make([][]I, argmax+1) // better than map I guess...

	for i := range se {
		md := madi(se[i][0], se[i][1])
		cp := se[i][0] // current pos
		// generate diamond pattern
		// TODO: only need quarter of the loop as it is symetric...
		for dy := -md; dy <= md; dy++ {
			y := cp[1] + dy
			if y < 0 || y > argmax {
				continue
			}
			if len(bl[y]) == 1 && bl[y][0][0] == 0 && bl[y][0][1] == argmax && y != argy {
				continue
			}
			sgn := 1
			if dy > 0 {
				sgn = -1
			}
			f := cp[0] - md - (dy * sgn)
			t := cp[0] + md + (dy * sgn)

			// cut intervals (except for argy) to argmax boundaries; omit argy
			if y != argy {
				if t < 0 || f > argmax {
					continue
				}
				if f < 0 {
					f = 0
				}
				if t > argmax {
					t = argmax
				}
			}
			bl[y] = add(bl[y], I{f, t})
		}
	}

	ip := 0 // impossible positions
	for _, iv := range bl[argy] {
		ip += iv[1] - iv[0]
	}
	log.Print(ip)

	fq := 0 // distress beacon frequency
	for y := range bl {
		if len(bl[y]) < 2 {
			continue
		} // only one Y will have only one gap
		fq = ((bl[y][0][1] + 1) * 4000000) + y
		break
	}
	log.Print(fq)
}
