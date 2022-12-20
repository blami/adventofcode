package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Monkey struct {
	l, r string // left right monkeys
	op byte
	num int
	done bool
}


func main() {

	ms := map[string]Monkey{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := strings.Split(s.Text(), ":")

		m := Monkey{}
		nam := l[0] // name
		job := strings.Trim(l[1], " ") // job
		num, err := strconv.Atoi(job) // number
		if err != nil {
			lrop := strings.Split(job, " ")
			m.l = lrop[0]
			m.op = lrop[1][0] // +-*/
			m.r = lrop[2]
		} else {
			m.num = num
			m.done = true
		}

		ms[nam] = m
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// resolve
	i := 0
	done := false
	for !done {
		done = true
		for nam := range ms {
			m := ms[nam]
			if m.done { continue }
			if ms[m.l].done && ms[m.r].done {
				switch ms[nam].op {
				case '+':
					m.num = ms[m.l].num + ms[m.r].num
				case '-':
					m.num = ms[m.l].num - ms[m.r].num
				case '*':
					m.num = ms[m.l].num * ms[m.r].num
				case '/':
					m.num = ms[m.l].num / ms[m.r].num
				default:
					log.Fatal(nam, ".op: ", m.op)
				}
				m.done = true
				ms[nam] = m
			}
			if !ms[nam].done { done = false }
		}
		i++
	}

	log.Print(i, " ", ms, " ", ms["root"].num)
}
