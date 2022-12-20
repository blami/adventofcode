package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// A Monkey.
type Monkey struct {
	l, r string // left right monkeys
	op   byte
	num  float64
}

// Map of monkeys.
type Monkeys map[string]Monkey

// Resolve given monkey's number without mutating underlying map.
func (ms *Monkeys) resolve(n string, guess bool, humn float64) float64 {
	if guess && n == "humn" {
		return humn
	}
	m := (*ms)[n]
	if m.op == '#' {
		return m.num
	}
	l := ms.resolve(m.l, guess, humn)
	r := ms.resolve(m.r, guess, humn)
	switch m.op {
	case '+':
		return l + r
	case '-':
		return l - r
	case '*':
		return l * r
	case '/':
		return l / r
	default:
		log.Fatalf("%s: unknown op %c", n, m.op)
		return 0 //duh
	}
}

func main() {

	ms := Monkeys{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := strings.Split(s.Text(), ":")

		m := Monkey{}
		nam := l[0]                    // name
		job := strings.Trim(l[1], " ") // job
		num, err := strconv.Atoi(job)  // number
		if err != nil {
			lrop := strings.Split(job, " ")
			m.l = lrop[0]
			m.op = lrop[1][0] // +-*/
			m.r = lrop[2]
		} else {
			m.num = float64(num)
			m.op = '#' // monkey just tells number
		}

		ms[nam] = m
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	root := ms.resolve("root", false, 0)

	// this is to get halving condition right
	gt := ms.resolve(ms["root"].r, true, 0) > ms.resolve(ms["root"].l, true, 0)

	e := ms.resolve(ms["root"].r, true, 0) // resolve root.l with humn=0; expected value
	l, h := int(-root*5), int(root*5)      // low, high for binsearch; TOGGLE if needed
	humn := 0
	for l < h {
		humn = (l + h) / 2                                 // shoot in mid
		a := ms.resolve(ms["root"].l, true, float64(humn)) // resolve root.r with humn=(l+h)/2 as a guess
		//log.Printf("l=%d h=%d humn=%d e=%f a=%f", l, h, humn, e, a)
		if a == e {
			break // humn causes match
		}
		if (!gt && a < e) || (gt && a > e) {
			h = humn
		} else {
			l = humn
		} // binary halving
	}

	log.Print(int(root))
	log.Print(humn)
}
