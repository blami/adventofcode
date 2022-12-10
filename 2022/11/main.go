package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Monkey struct {
	Items []int
	Op string // needs to be eval()'d
	Test int // seems to be always "divisble by"
	True int // is always throw to N
	False int // is always throw to N
	Insp int // number of inspected items
}

// Evaluate given operation given the old value and return new value.
func eval(op string, old int) int {
	a, b := 0, 0
	var err error
	l := strings.Split(op, " ")
	if l[0] == "old" { a = old } else { 
		a, err = strconv.Atoi(l[0]) 
		if err != nil { log.Fatal("wrong operand a: ", l) }
	}
	if l[2] == "old" { b = old } else {
		b, err = strconv.Atoi(l[2])
		if err != nil { log.Fatal("wrong operand b: ", l) }
	}
	switch l[1] {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		log.Fatal("operation ", l[1], " not implemented!")
	}
	return -1
}

// Run monkey business for given number of rounds and worry relieve factor.
func business(monkeys []Monkey, rounds int, relieve int) []Monkey {
	for rounds > 0 {
		for n, m := range(monkeys) {
			//log.Print("Monkey ", n)
			for _, oold := range(m.Items) {
				//log.Print(" inspecting ", oold)
				nnew := eval(m.Op, oold)
				//log.Print(" changed ", nnew)
				nnew = nnew / relieve // relieve factor
				//log.Print(" relieved ", nnew)
				if nnew % m.Test == 0 {
					//log.Print(" throwing to ",nnew,  m.True)
					monkeys[m.True].Items = append(monkeys[m.True].Items, nnew)
				} else {
					//log.Print(" throwing to ", nnew, m.False)
					monkeys[m.False].Items = append(monkeys[m.False].Items, nnew)
				}
				// count inspected items; DO NOT use m. as that is copy
				monkeys[n].Insp++
			}
			// empty items; DO NOT use m. as that is copy
			monkeys[n].Items = nil
		}
		rounds--
	}
	return monkeys
}

// Return product of two most active monkeys inspection counts.
func most(out []Monkey) int {
	mam := [2]int{}
	for _, m := range(out) {
		if m.Insp > mam[0] {
			mam[1] = mam[0]
			mam[0] = m.Insp
		} else if m.Insp > mam[1] {
			mam[1] = m.Insp
		}
	}
	return mam[0] * mam[1]
}


func main() {

	// NOTE: thought of using map but monkeys go in sequence...
	monkeys := []Monkey{}

	// Process the input
	s := bufio.NewScanner(os.Stdin)
	// BUG: this is sloppy in case line would be so long it does not fit
	// internal buffer but for input.txt is fine
	for s.Scan() {
		l := s.Text()

		// NOTE: assuming good input only
		if len(l) > 1 && l[0] != ' ' { // top level Monkey section
			n, _ := strconv.Atoi(strings.Split(l[:len(l)-1], " ")[1])
			m := Monkey{}
			// read Monkey block
			for {
				s.Scan()
				l := strings.Trim(s.Text(), " ")
				if l == "" {
					break
				}
				// NOTE: rather handle Atoi, etc. errors here as we don't know
				// the full syntax of input from examples...
				switch {
				case l[:2] == "St": // Starting items
					for _, v := range(strings.Split(strings.Trim(strings.Split(l, ":")[1], " "), ",")) {
						item, err := strconv.Atoi(strings.Trim(v, " "))
						if err != nil {
							log.Fatal("invalid start item: ", l)
						}
						m.Items = append(m.Items, item)
					}
				case l[:2] == "Op": // Operation
					m.Op = strings.Trim(strings.Split(strings.Split(l, ":")[1], "=")[1], " ")
				case l[:2] == "Te": // Test
					// Seems to be always "divisble by" so save only the integer
					var err error
					m.Test, err = strconv.Atoi(strings.Split(strings.Split(l, ":")[1], " ")[3])
					if err != nil {
						log.Fatal("invalid test divsible by: ", l)
					}
				case l[:2] == "If": // Conditions
					ll := strings.Split(l, " ")
					targ, err := strconv.Atoi(ll[5])
					if err != nil {
						log.Fatal("invalid throw target: ", l)
					}
					switch ll[1] {
					case "true:":
						m.True = targ
					case "false:":
						m.False = targ
					default:
						log.Fatal("invalid condition: ", l)
					}
				}
			}
			if n != len(monkeys) {
				log.Fatal("monkey sequence broken")
			}
			monkeys = append(monkeys, m)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// Run monkey business
	out1 := business(append([]Monkey{}, monkeys...), 20, 3)
	//out2 := business(append([]Monkey{}, monkeys...), 10000, 1)

	log.Print(most(out1))
	//log.Print(most(out2))
}
