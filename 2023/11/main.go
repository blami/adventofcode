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

// lol
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// expand space
func exp(gs []XY, e [][]int, f int) []XY {
	ns := []XY{}
	for _, g := range gs {
		n := XY{g.X, g.Y}
		for _, y := range e[0] {
			if y > g.Y {
				break
			}
			n.Y += f - 1
		}
		for _, x := range e[1] {
			if x > g.X {
				break
			}
			n.X += f - 1
		}
		ns = append(ns, n)
	}
	return ns
}

// find sum of shortest distances between galaxies
func sdsum(gs []XY) int {
	done := map[XY]bool{}
	sum := 0
	// part1
	for srci, srcxy := range gs {
		for dsti, dstxy := range gs {
			if dsti == srci {
				continue
			}
			_, sd := done[XY{srci, dsti}]
			_, ds := done[XY{dsti, srci}]
			if sd || ds {
				continue
			}
			// shortest path length
			d := abs(dstxy.Y-srcxy.Y) + abs(dstxy.X-srcxy.X)
			done[XY{srci, dsti}] = true
			sum += d
		}
	}
	return sum
}

func main() {
	w := 0
	e := [][]int{{}, {}} // empty 0-rows, 1-columns
	gs := []XY{}

	s := bufio.NewScanner(os.Stdin)
	sy := 0
	var hcc []int // hash count in column
	for s.Scan() {
		l := s.Text()
		w = len(l)
		if hcc == nil {
			hcc = make([]int, w)
		}
		if strings.Index(l, "#") == -1 {
			e[0] = append(e[0], sy)
		}
		for sx, c := range l {
			if c == '#' {
				hcc[sx] += 1
				gs = append(gs, XY{sx, sy})
			}
		}
		sy++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	for i := range hcc {
		if hcc[i] == 0 {
			e[1] = append(e[1], i)
		}
	}

	fmt.Println(sdsum(exp(gs, e, 2)))
	fmt.Println(sdsum(exp(gs, e, 1000000)))
}
