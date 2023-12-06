// Parse almanac and find lowest seed destination.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

func main() {
	seeds := []int{}
	alm := [][][]int{}

	// parse almanac
	m := -1
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if strings.HasPrefix(l, "seeds:") {
			for _, ns := range strings.Fields(strings.Split(l, ":")[1]) {
				n, _ := strconv.Atoi(ns)
				seeds = append(seeds, n)
			}
			continue
		}
		if l == "" {
			continue
		}
		if !unicode.IsDigit(rune(l[0])) {
			alm = append(alm, [][]int{})
			m++
		} else {
			v := []int{}
			for _, ns := range strings.Fields(l) {
				n, _ := strconv.Atoi(ns)
				v = append(v, n)
			}
			alm[m] = append(alm[m], v)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// find location
	locmin := -1
	for _, seed := range seeds {
		loc := seed
		// go through map
		for _, sec := range alm {
			for _, row := range sec {
				if loc >= row[1] && loc < (row[1]+row[2]) {
					df := (loc - row[1])
					loc = row[0] + df
					break
				}
			}
		}
		if loc < locmin || locmin == -1 {
			locmin = loc
		}
	}

	// part 2
	// find location over ranges
	locmin2 := -1
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	for i := 0; i < len(seeds)-2; i += 2 {
		// parallel bruteforce
		wg.Add(1)
		go func(s int, l int) {
			defer wg.Done()
			locminsub := -1
			for seed := s; seed < s+l; seed++ {
				loc := seed
				for _, sec := range alm {
					for _, row := range sec {
						if loc >= row[1] && loc < (row[1]+row[2]) {
							df := (loc - row[1])
							loc = row[0] + df
							break
						}
					}
				}
				if loc < locminsub || locminsub == -1 {
					locminsub = loc
				}
			}
			mu.Lock()
			if locminsub < locmin2 || locmin2 == -1 {
				locmin2 = locminsub
			}
			mu.Unlock()
		}(seeds[i], seeds[i+1])
	}
	wg.Wait()

	fmt.Println(locmin)
	fmt.Println(locmin2)
}
