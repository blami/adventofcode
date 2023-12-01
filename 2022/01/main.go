// Read input from stdin and calculate total calories per elf. Store
// top three highest number of calories in rank and print first position
// and sum of first three positions.

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// Type alias to hold first 3 positions in elf ranking.
type Ranking [3]int

// Push value into ranking, if not good enough it will fall through.
func (r *Ranking) Push(v int) {
	switch {
	case v >= (*r)[0]:
		(*r)[2] = (*r)[1]
		(*r)[1] = (*r)[0]
		(*r)[0] = v
	case v >= (*r)[1]:
		(*r)[2] = (*r)[1]
		(*r)[1] = v
	case v >= (*r)[2]:
		(*r)[2] = v
	}
}

// Return top place.
func (r *Ranking) Top() int {
	return (*r)[0]
}

// Return sum of top three places.
func (r *Ranking) Sum() int {
	return (*r)[0] + (*r)[1] + (*r)[2]
}

func main() {

	r := Ranking{}
	cur := 0 // current elf acummulator

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			r.Push(cur)
			cur = 0
		} else {
			v, _ := strconv.Atoi(l) // assuming good input
			cur += v
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(r.Top())
	log.Print(r.Sum())
}
