// Find multiplication of record beating attempts for three rides and number of
// ways to beat final race.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	ts := []int{} // times
	ds := []int{} // destination

	// load times and distances
	i := 0
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		merg := ""
		for _, ns := range strings.Fields(l) {
			n, err := strconv.Atoi(ns)
			if err != nil {
				continue
			}
			merg += ns
			if i == 0 {
				ts = append(ts, n)
			}
			if i == 1 {
				ds = append(ds, n)
			}
		}
		n, _ := strconv.Atoi(merg)
		if i == 0 {
			ts = append(ts, n)
		}
		if i == 1 {
			ds = append(ds, n)
		}
		i++
	}

	// all races
	tot := 1
	ww := 0
	for r := 0; r < len(ts); r++ {
		sum := 0
		// simulate an attempt
		for sp := 0; sp <= ts[r]; sp++ {
			tr := ts[r] - sp
			dt := sp * tr
			if dt > ds[r] && r != len(ts)-1 {
				sum++
			}
			if dt > ds[r] && r == len(ts)-1 {
				ww++
			}
		}
		if sum > 0 {
			tot *= sum
		}
	}

	fmt.Println(tot)
	fmt.Println(ww)
}
