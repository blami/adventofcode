// Simulate circuit and find number of high and low signals after 1000 button
// presses; find number of button presses needed to get rx low.
//
// NOTE: This solution is not universal and assumes that one & module leads to
// rx and N modules leads to that & module...
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	//"strconv"
	"strings"
)

type Module struct {
	Name string
	Type rune
	Dest []string
	On   bool            // %
	Mem  map[string]bool // &
}

type Signal struct {
	From string
	To   string
	Val  bool
}

// Generate graphviz dot file
func dot(ms []Module) {
	fmt.Println("#dot -Tpng out.dot > out.png")
	style := ""
	fmt.Println("digraph D {")
	for _, m := range ms {
		t := string(m.Type)
		if t == "b" {
			t = ""
		}
		if t == "%" {
			t = "\\%"
		}
		if t == "&" {
			style += fmt.Sprintf("\"%s%s\" [style=filled, color=\"#ffb3b2\"]\n", t, m.Name)
		}
		fmt.Printf("\"%s%s\" -> {", t, m.Name)
		for _, d := range m.Dest {
			mi := find(ms, d)
			if mi == -1 {
				fmt.Printf(" \"%s\" ", d)
			} else {
				t := string(ms[mi].Type)
				if t == "b" {
					t = ""
				}
				if t == "%" {
					t = "\\%"
				}
				fmt.Printf(" \"%s%s\" ", t, ms[mi].Name)
			}
		}
		fmt.Println("}")
	}
	fmt.Print(style)
	fmt.Println("}")
}

func debug(si Signal) {
	lh := "lo"
	if si.Val {
		lh = "hi"
	}
	fmt.Printf("%s -%s-> %s\n", si.From, lh, si.To)
}

func find(ms []Module, name string) int {
	for i := range ms {
		if ms[i].Name == name {
			return i
		}
	}
	return -1
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(x ...int) int {
	// BUG: len(x) >= 2
	r := x[0] * x[1] / gcd(x[0], x[1])
	for i := 0; i < len(x[2:]); i++ {
		r = lcm(r, x[i+2])
	}
	return r
}

func main() {
	ms := []Module{}
	s := bufio.NewScanner(os.Stdin)
	has_rx := false
	for s.Scan() {
		l := s.Text()
		f := strings.Split(l, " -> ")
		nam := f[0]
		typ := rune(f[0][0])
		if typ == '%' || typ == '&' {
			nam = f[0][1:]
		}
		f = strings.Split(f[1], ", ")
		for i := range f { // make test.txt work (part 1 only)
			if f[i] == "output" {
				f[i] = "rx" // fix for test2.txt
			}
			if f[i] == "rx" {
				has_rx = true
				break
			}
		}
		ms = append(ms, Module{nam, typ, f, false, nil})
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	for i := range ms {
		if ms[i].Type == '&' {
			ms[i].Mem = make(map[string]bool)
			// find all inputs and add them to Ins
			for _, m := range ms {
				for _, d := range m.Dest {
					if d == ms[i].Name {
						ms[i].Mem[m.Name] = false
					}
				}
			}
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "-dot" {
		dot(ms)
		return
	}

	// TODO: Find something more generic
	// part 2
	sist := map[string]map[string]bool{} // signal state on inputs of each part
	ins := map[string][]string{}         // inverted dict (inputs to parts)
	conj := []string{}
	if has_rx {
		for _, m := range ms {
			for _, d := range m.Dest {
				if _, ok := ins[d]; !ok {
					ins[d] = []string{}
				}
				mi := find(ms, d)
				if mi != -1 && ms[mi].Type == '&' {
					if _, ok := sist[d]; !ok {
						sist[d] = map[string]bool{}
					}
					sist[d][m.Name] = false
				}
				ins[d] = append(ins[d], m.Name)
			}
		}
		// BUG: add assertions as this expects one & module input to rx
		conj = ins[ins["rx"][0]] // list of & modules above & module leading to rx
	}

	// simulate circuit
	bn := 0              // button presses
	lohi := [2]int{0, 0} // part 1

	// part 2
	rxlow := 0 // button presses to make rx low
	// cycle detection
	prevlow := map[string]int{} // previous bn when module was low
	nlow := map[string]int{}    // number of times module was seen low
	loop := []int{}

	for {
		bn += 1
		q := []Signal{
			{"button", "broadcaster", false},
		}
		//fmt.Println("---")
		for len(q) > 0 {
			// pop
			si := q[0]
			q = q[1:]

			if has_rx {
				if si.Val == false {
					_, inprev := prevlow[si.To]
					inconj := false
					for _, c := range conj {
						if c == si.To {
							inconj = true
							break
						}
					}
					if inprev && nlow[si.To] == 2 && inconj {
						loop = append(loop, bn-prevlow[si.To])
					}
					prevlow[si.To] = bn
					nlow[si.To] += 1
				}
				if len(loop) == len(conj) {
					rxlow = lcm(loop...)
					if bn > 1000 {
						break // have everything
					}
				}
			}

			// bruteforce part2
			// takes too long but works...
			/*
				if has_rx && (si.To == "rx" || si.To == "output") && si.Val == false {
					rxlow = bn
					if bn > 1000 {
						break // have everything
					}
				}
			*/

			// part 1
			//debug(si)
			if bn <= 1000 {
				if si.Val {
					lohi[1] += 1
				} else {
					lohi[0] += 1
				}
			}

			mi := find(ms, si.To)
			if mi == -1 {
				continue
			}
			switch ms[mi].Type {
			case 'b': // broadcaster, just broadcast
				for _, d := range ms[mi].Dest {
					q = append(q, Signal{ms[mi].Name, d, si.Val})
				}
			case '%': // flipflop
				if si.Val {
					continue
				}
				ms[mi].On = !ms[mi].On
				for _, d := range ms[mi].Dest {
					q = append(q, Signal{ms[mi].Name, d, ms[mi].On})
				}
			case '&': // conjunction
				ms[mi].Mem[si.From] = si.Val
				val := false
				for _, v := range ms[mi].Mem {
					if !v {
						val = true
						break
					}
				}
				for _, d := range ms[mi].Dest {
					q = append(q, Signal{ms[mi].Name, d, val})
				}
			}
		}

		if bn > 1000 && (rxlow > 0 || !has_rx) { // have everything
			break
		}
	}

	fmt.Println(lohi[0] * lohi[1])
	if has_rx {
		fmt.Println(rxlow)
	} else {
		fmt.Println("circuit has no rx/output")
	}
}
