// Read input from stdin and run monkey business while tracking how many items
// each monkey inspected.

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// A monkey.
type Monkey struct {
	Items []int  // Currently held item worry levels.
	Op    string // Operation applied to worry level while inspecting item.
	Test  int    // Divisor in test.
	True  int    // Target if test passes.
	False int    // Target if test fails.
	Insp  int    // Number of inspected items.
}

// Evaluate given operation given the old value and return new value.
func eval(op string, old int) int {
	opr := [2]int{} // operands
	l := strings.Split(op, " ")
	for i, j := range []int{0, 2} {
		if l[j] == "old" {
			opr[i] = old
		} else {
			opr[i], _ = strconv.Atoi(l[j])
		}
	}
	switch l[1] {
	case "+":
		return opr[0] + opr[1]
	case "-":
		return opr[0] - opr[1]
	case "*":
		return opr[0] * opr[1]
	case "/":
		return opr[0] / opr[1]
	default:
		log.Fatal("operation ", l[1], " not implemented!")
	}
	return -1
}

// Run monkey business for given number of rounds using given worry level
// management function.
func business(monkeys []Monkey, rounds int, w func(int) int) []Monkey {
	for rounds > 0 {
		for i, m := range monkeys {
			for _, o := range m.Items {
				n := eval(m.Op, o)
				n = w(n) // change worry level
				if n%m.Test == 0 {
					monkeys[m.True].Items = append(monkeys[m.True].Items, n)
				} else {
					monkeys[m.False].Items = append(monkeys[m.False].Items, n)
				}
				monkeys[i].Insp++
			}
			// empty items; DO NOT use m.Items = nil as that is a copy
			monkeys[i].Items = nil
		}
		rounds--
	}
	return monkeys
}

// Find two most active monkeys and return product of number of items they
// inspected.
func most(out []Monkey) int {
	mam := [2]int{}
	for _, m := range out {
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

	monkeys := []Monkey{}
	// In part 2 worry level numbers easily overflow int, to manage we
	// calculate least common multiplier of all divisors and instead of /3
	// (part 1) we %lcm to keep worry levels "right" but within int range.
	lcm := 1 // least common multiple of test divisors

	s := bufio.NewScanner(os.Stdin)
	// BUG: this is sloppy in case line would be so long it does not fit
	// internal buffer but for input.txt is fine
	for s.Scan() {
		l := s.Text()

		// NOTE: assuming good input only and monkeys going in order from 0 to
		// N; not particularly proud about this.
		if len(l) > 1 && l[0] != ' ' { // top level Monkey section
			m := Monkey{}
			for {
				s.Scan()
				l := strings.Trim(s.Text(), " ")
				if l == "" {
					break
				}
				switch {
				case l[:2] == "St": // Starting items
					for _, v := range strings.Split(strings.Trim(strings.Split(l, ":")[1], " "), ",") {
						item, _ := strconv.Atoi(strings.Trim(v, " "))
						m.Items = append(m.Items, item)
					}
				case l[:2] == "Op": // Operation
					m.Op = strings.Trim(strings.Split(strings.Split(l, ":")[1], "=")[1], " ")
				case l[:2] == "Te": // Test
					// Seems to be always "divisble by" so save only the integer
					m.Test, _ = strconv.Atoi(strings.Split(strings.Split(l, ":")[1], " ")[3])
					lcm *= m.Test
				case l[:2] == "If": // Conditions
					ll := strings.Split(l, " ")
					targ, _ := strconv.Atoi(ll[5])
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
			// works for my input, not sure whether monkeys can come in
			// different order...
			monkeys = append(monkeys, m)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// run monkey business
	out1 := business(append([]Monkey{}, monkeys...), 20, func(w int) int { return w / 3 })
	out2 := business(append([]Monkey{}, monkeys...), 10000, func(w int) int { return w % lcm })

	log.Print(most(out1))
	log.Print(most(out2))
}
