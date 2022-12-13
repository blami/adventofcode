// Read input from stdin and process packets. Using custom packet comparator
// find number of packet pairs out of order and also decoder key.

package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// Split packet list by top level ','s to separate top level components (which
// can still be lists).
func split(s string) []string {
	// remove outer []
	if s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}

	out := []string{}
	i := 0 // current part start
	b := 0 // bracket level
	for j, r := range s {
		//log.Printf("split r=%s i=%d j=%d b=%d", string(r), i, j, b)
		switch {
		case r == '[':
			b++
		case r == ']': // can be end
			b--
			if b == 0 && j == len(s)-1 {
				out = append(out, s[i:j+1])
			}
		case b == 0 && r == ',': // cannot be end
			out = append(out, s[i:j])
			i = j + 1
		case j == len(s)-1:
			out = append(out, s[i:j+1])
		}
	}
	return out
}

// Check if given string (as in packet component) is a list.
func isl(s string) bool {
	return s[0] == '['
}

// Compare two packets recursively. Outcome can be either 1 which is good order
// left side being "lower" than right side or -1 bad order or 0 for equal.
func cmp(ls, rs string, d int) bool {
	l, r := split(ls), split(rs)
	cl := len(l)
	if len(r) < cl {
		cl = len(r)
	}

	// debug print
	ind := ""
	for i := 0; i < d; i++ {
		ind += " "
	}
	//log.Printf("%scmp %v(%d) ? %v(%d) cl=%d {", ind, l, len(l), r, len(r), cl)

	out := 0
	// true 90s programmers use labels
F:
	for i := 0; i < cl; i++ {
		// NOTE: This part is terribly complex and could be written simpler. I
		// still could not get it right...
		switch {
		case !isl(l[i]) && !isl(r[i]): // int int comparison
			lint, _ := strconv.Atoi(l[i])
			rint, _ := strconv.Atoi(r[i])
			//log.Printf("%s%d ? %d=%v", ind, lint, rint, lint <= rint)
			if lint < rint {
				out = 1
				break F
			}
			if lint > rint {
				out = -1
				break F
			}

		case isl(l[i]) && !isl(r[i]): // [] int comparison
			if out = cmp(l[i], "["+r[i]+"]", d+1); out != 0 {
				break F
			}

		case !isl(l[i]) && isl(r[i]): // int [] comparison
			if out = cmp("["+l[i]+"]", r[i], d+1); out != 0 {
				break F
			}

		case isl(l[i]) && isl(r[i]): // [] [] comparison
			if out = cmp(l[i], r[i], d+1); out != 0 {
				break F
			}
		}
	}
	// run outs
	if out == 0 {
		switch {
		case len(r) < len(l):
			out = -1
		case len(r) > len(l):
			out = 1
		}
		// NOTE: this still can end up as len(l) == len(r) hence out will be 0
		// and one needs to look at next item.
	}
	//log.Printf("%s}=%v", ind, out)

	return out
}

func main() {

	n := 1
	gp := 0
	sp := []string{"[[2]]", "[[6]]"} // sorted packets

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		s.Scan()
		r := s.Text()
		sp = append(sp, l, r) // store both packets for further sorting

		//log.Printf("")
		//log.Printf("<< pair %d, %v ? %v", n, l, r)
		out := cmp(l, r, 0)
		//log.Printf(">> pair %d out=%v", n, out)
		if out == 0 {
			log.Fatal("this shouldn't happen...")
		}
		if out == 1 {
			gp += n
		}
		n++
		s.Scan() // scan empty line at the end of pair
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// sort packets using our comparator
	sort.Slice(sp, func(i, j int) bool {
		return cmp(sp[i], sp[j], 0) == 1
	})

	dk := 1 // decryption key
	for i, p := range sp {
		if p == "[[2]]" || p == "[[6]]" {
			dk *= i + 1
		}
	}

	log.Print(gp)
	log.Print(dk)
}
