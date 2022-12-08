// Read input from stdin and find overlapping sections.

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Simple datatype to express an interval for easier work than with arrays
type Interval struct {
	from int
	to   int
}

// Create new interval from "FROM-TO" formatted string
func NewInterval(s string) Interval {
	ss := strings.Split(s, "-")
	from, err := strconv.Atoi(ss[0])
	if err != nil {
		log.Fatal(err)
	}
	to, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatal(err)
	}

	return Interval{from, to}
}

// Check whether interval i fully contains other interval j
func (i Interval) Contains(j Interval) bool {
	return j.from >= i.from && j.to <= i.to
}

// Check whether interval i overlaps with interval j
func (i Interval) Overlaps(j Interval) bool {
	return (j.to >= i.from && j.from <= i.to) || (j.from >= i.to && j.to <= i.from)
}

func main() {

	nf := 0 // number of fully contained sections between pairs
	no := 0 // number of overlapping sections between pairs

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		i := strings.Split(s.Text(), ",")
		i1 := NewInterval(i[0])
		i2 := NewInterval(i[1])

		// count fully contained intervals
		if i1.Contains(i2) || i2.Contains(i1) {
			nf++
		}

		// count overlapping intervals
		if i1.Overlaps(i2) {
			no++
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(nf)
	log.Print(no)
}
