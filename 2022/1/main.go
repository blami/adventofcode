// Read input from stdin and calculate total calories per elf. Store
// top three highest number of calories in rank and print first position
// (part 1) and sum of first three positions (part 2).

package main

import (
	"os"
	"bufio"
	"strconv"
	"log"
)

type Ranking struct {
	pos [3]int	
}

func (r *Ranking) Push(v int) {
	switch {
	case v >= r.pos[0]:
		r.pos[2] = r.pos[1]
		r.pos[1] = r.pos[0]
		r.pos[0] = v
	case v >= r.pos[1]:
		r.pos[2] = r.pos[1]
		r.pos[1] = v
	case v >= r.pos[2]:
		r.pos[2] = v
	}
}

func (r *Ranking) Sum() int {
	return	r.pos[0] + r.pos[1] + r.pos[2]
}

func (r *Ranking) Top() int {
	return r.pos[0]
}

func main() {

	rank := Ranking{}
	cur := 0	// current elf acummulator

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		if l == "" {
			rank.Push(cur)
			cur = 0
		} else {
			v, err := strconv.Atoi(l)
			if err != nil {
				log.Fatal(err)
			}
			cur += v
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	log.Print(rank.Top())
	log.Print(rank.Sum())
}
