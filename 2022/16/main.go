// Read input from stdin and find the combination of open valves in 30 minutes
// and in 26 minutes and 2 actors.

package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// A valve as parsed from input.
type Valve struct {
	rate int
	next []string
}

// Max of a,b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Whether or not s exists in slice s.
func in(sl []string, s string) bool {
	for _, t := range sl {
		if s == t {
			return true
		}
	}
	return false
}

var cach map[string]int // poor man's cache; key is cv + ov + t + el as string

// Recursive (and slow) DFS with poor man's caching. vs is map of Valves with
// their neighbour names and flow rate, cv is current valve, ov is slice of
// opened valves, t is remaining time until eruption, el is whether or not
// elephant is involved and d is depth of recursion.
func maxflow(vs map[string]Valve, cv string, ov []string, t int, el bool, d int) int {
	if d == 0 {
		cach = nil
	}
	if t <= 0 {
		if el {
			//log.Printf("ele %d cv=%s ov=%v", d, cv, ov)
			// if elephant is involved at the end of human work run elephant
			return maxflow(vs, "AA", append([]string{}, ov...), 26, false, d)
		} else {
			return 0
		}
	}
	sort.Strings(ov) // sort needed to check for cache hit
	if bf, ok := cach[cv+strings.Join(ov, ",")+strconv.Itoa(t)+strconv.FormatBool(el)]; ok {
		return bf
	}

	bf := 0
	for _, nv := range vs[cv].next {
		bf = max(bf, maxflow(vs, nv, ov, t-1, el, d+1))
	}
	if !in(ov, cv) && vs[cv].rate > 0 && t > 0 /* have time left for t-2? */ {
		f := (t - 1) * vs[cv].rate
		for _, nv := range vs[cv].next {
			//log.Print(ov, " ", cov)
			// BUG: slice copy is needed otherwise we mutate underlying array
			// of d-1
			bf = max(bf, f+maxflow(vs, nv, append([]string{cv}, ov...), t-2, el, d+1))
		}
	}

	if cach == nil {
		cach = make(map[string]int)
	}
	sort.Strings(ov) // sorted ov makes cache more effective (vs cost of sort of such short slice)
	cach[cv+strings.Join(ov, ",")+strconv.Itoa(t)+strconv.FormatBool(el)] = bf
	//log.Printf("%d cv=%s ov=%v bf=%d", t, cv, ov, bf)
	return bf
}

func main() {
	vs := map[string]Valve{}

	re := regexp.MustCompile(`Valve (?P<v>[A-Z]{2}) has flow rate=(?P<r>\d+); tunnel[s]? lead[s]? to valve[s]? (?P<t>.*)$`)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		mtch := re.FindStringSubmatch(s.Text())
		if len(mtch) == 0 {
			log.Fatalf("invalid input: %s", s.Text())
		}
		r, _ := strconv.Atoi(mtch[2])
		v := Valve{
			r,
			strings.Split(strings.ReplaceAll(mtch[3], " ", ""), ","),
		}
		vs[mtch[1]] = v
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(maxflow(vs, "AA", []string{}, 30, false, 0)) // part1
	log.Print(maxflow(vs, "AA", []string{}, 26, true, 0))  // part2
}
