// Parse each line as game, produce sum of possible game ids and sum of minimum
// power of minimum cube numbers.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Set struct {
	R int
	G int
	B int
}

type Game struct {
	Id   int
	Sets []Set
}

func parse(l string) Game {
	id := ""
	// Parse id
	ss := 0
	for i, c := range l {
		if i < 5 {
			continue
		}
		if c == ':' {
			ss = i
			break
		}
		id += string(c)
	}
	// Parse sets
	idd, _ := strconv.Atoi(id)
	gam := Game{
		Id:   idd,
		Sets: []Set{},
	}
	s := Set{}
	for i := ss + 1; i < len(l); i++ {
		c := l[i]
		if c == ';' || i == len(l)-1 {
			gam.Sets = append(gam.Sets, s)
			s = Set{}
			continue
		}
		if c >= 47 && c <= 58 {
			// Chomp all numbers from here
			j := 0
			for j = i; l[j] >= 47 && l[j] <= 58; j++ {
			}
			n, _ := strconv.Atoi(l[i:j])
			// Decide color
			switch l[j+1] {
			case 'r':
				s.R += n
				i = j + 2
			case 'g':
				s.G += n
				i = j + 4
			case 'b':
				s.B += n
				i = j + 3
			}
			continue
		}
	}
	//fmt.Println(gam)
	return gam
}

func main() {
	sum := 0
	sumpow := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()
		gam := parse(l)

		ok := true
		r := 0
		g := 0
		b := 0
		for _, s := range gam.Sets {
			// part 1
			if s.R > 12 || s.G > 13 || s.B > 14 {
				ok = false
			}
			// part 2
			if s.R > r {
				r = s.R
			}
			if s.G > g {
				g = s.G
			}
			if s.B > b {
				b = s.B
			}
		}
		if ok {
			sum += gam.Id
		}
		sumpow += r * g * b
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
	fmt.Println(sumpow)
}
