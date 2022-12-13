
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)


// Split packet list by top level ','s to parts.
func split(s string) []string {
	// remove outer []
	if s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1:len(s)-1]
	}

	out := []string{}

	i := 0 // current part start
	b := 0 // bracket level
	for j, r := range(s) {
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

// Returns true if given string is list.
func isl(s string) bool {
	return s[0] == '['
}

// Compare two packets recursively (d is recursion depth).
func cmp(ls, rs string, d int) int {
	l, r := split(ls), split(rs)
	cl := len(l)
	if len(r) < cl { cl = len(r) }

	// debug print
	ind := "" 
	for i := 0; i < d; i++ {ind += " "}
	log.Printf("%scmp %v(%d) ? %v(%d) cl=%d {", ind, l, len(l), r, len(r), cl)

	out := 0
F:
	for i := 0; i < cl; i++ {
		switch {
		case !isl(l[i]) && !isl(r[i]): // int int comparison
			lint, _ := strconv.Atoi(l[i])
			rint, _ := strconv.Atoi(r[i])
			log.Printf("%s%d ? %d=%v", ind, lint, rint, lint <= rint)
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
	if out == 0 {
		switch {
			case len(r) < len(l):
				out = -1
			case len(r) > len(l):
				out = 1
			}
	}
	log.Printf("%s}=%v", ind, out)

	return out
}


func main() {

	n := 1
	gp := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		s.Scan()
		r := s.Text()

		log.Printf("")
		log.Printf("<< pair %d, %v ? %v", n, l, r)
		out := cmp(l, r, 0)
		log.Printf(">> pair %d out=%v", n, out)
		if out == 0 {
			log.Fatal("oops")
		}
		if out == 1 {
			gp += n
		}

		n++
		s.Scan() // scan empty line
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(gp)
}
