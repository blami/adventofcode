// Find area of polygon given as set of direction + length plotting
// instructions.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type XY [2]int64

type T struct {
	Dir XY
	N   int64
}

func dir(s string) XY {
	switch s {
	case "R", "0":
		return XY{1, 0}
	case "D", "1":
		return XY{0, 1}
	case "L", "2":
		return XY{-1, 0}
	case "U", "3":
		return XY{0, -1}
	}
	panic("oops")
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func shoelace(ts []T) int64 {
	v := []XY{{0, 0}} // vertices
	var l int64 = 0   // length
	for i, t := range ts {
		l += t.N
		v = append(v, XY{v[i][0] + (t.Dir[0] * t.N), v[i][1] + (t.Dir[1] * t.N)})
	}

	// shoelace
	var sum int64 = 0
	for i := 0; i < len(v)-2; i++ {
		sum += (v[i][0] + v[i+1][0]) * (v[i][1] - v[i+1][1]) / 2
	}
	return abs(sum) + l/2 + 1
}

func main() {
	ts1 := []T{}
	ts2 := []T{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		f := strings.Fields(l)

		d := dir(f[0])
		n, _ := strconv.Atoi(f[1])
		ts1 = append(ts1, T{d, int64(n)})
		// part 2
		d = dir(f[2][len(f[2])-2 : len(f[2])-1])
		nn, err := strconv.ParseInt(f[2][2:len(f[2])-2], 16, 64)
		if err != nil {
			panic("err")
		}
		ts2 = append(ts2, T{d, nn})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(shoelace(ts1))
	fmt.Println(shoelace(ts2))
}
