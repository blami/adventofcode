package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type R struct {
	Var  string
	Op   rune
	Val  int
	Dest string
}

type W struct {
	Name  string
	Rules []R
}

type V struct {
	Var string
	Val int
}

func parseRule(rule string) R {
	r := R{}
	f := strings.Split(rule, ":")
	// fallback at the end of list has no :
	if len(f) == 1 {
		return R{f[0], '*', 0, ""}
	}

	r.Dest = f[1]
	opi := -1
	for _, op := range []rune{'<', '>'} {
		opi = strings.IndexRune(f[0], op)
		if opi != -1 {
			break
		}
	}
	if opi == -1 {
		panic("unknown operator in: " + rule)
	}
	r.Op = rune(f[0][opi])
	r.Var = f[0][0:opi]
	r.Val, _ = strconv.Atoi(f[0][opi+1:])
	return r
}

func parseWorkflow(wflow string) W {
	w := W{}
	for i := range wflow {
		if wflow[i] == '{' {
			w.Name = wflow[0:i]
			rs := strings.Split(wflow[i+1:len(wflow)-1], ",")
			for _, r := range rs {
				w.Rules = append(w.Rules, parseRule(r))
			}
			break
		}
	}
	return w
}

func parseVars(vars string) []V {
	v := []V{}
	//fmt.Println(vars)
	for _, as := range strings.Split(vars[1:len(vars)-1], ",") {
		//fmt.Println(as)
		f := strings.Split(as, "=")
		n, _ := strconv.Atoi(f[1])
		v = append(v, V{f[0], n})
	}
	return v
}

func comp(a, b int, op rune) bool {
	switch op {
	case '<':
		return a < b
	case '>':
		return a > b
	}
	panic("comp oops")
}

func run(ws map[string]W, v []V) (int, string) {
	sum := 0
	p := ""
	cur := "in"
	for cur != "A" && cur != "R" {
		//fmt.Println(ws[cur])
		p += cur + "->"
		dest := ""
		for _, r := range ws[cur].Rules {
			match := false
			if r.Op == '*' {
				dest = r.Var
				break
			}
			for _, vv := range v {
				if vv.Var == r.Var && comp(vv.Val, r.Val, r.Op) {
					dest = r.Dest
					//fmt.Println("  match", r, string(r.Op), "vars", v)
					match = true
					break
				}
			}
			if match {
				break
			}
		}
		cur = dest
		//fmt.Println(cur)
	}
	p += cur
	if cur == "A" {
		for _, vv := range v {
			sum += vv.Val
		}
	}
	//fmt.Println(v, p)
	return sum, p
}

func rangel(r map[string][2]int) int {
	l := 1
	for _, v := range r {
		l *= 1 + v[1] - v[0]
	}
	return l
}

func rangecp(ra map[string][2]int) map[string][2]int {
	rc := map[string][2]int{}
	for k, v := range ra {
		rc[k] = v
	}
	return rc
}

// part 2
func runR(ws map[string]W, ra map[string][2]int, cur string) int {
	n := 0
	//fmt.Println(ws[cur])
	for _, r := range ws[cur].Rules {
		if r.Op == '*' {
			if r.Var == "A" {
				n += rangel(ra)
			} else if r.Var != "R" {
				n += runR(ws, ra, r.Var)
			}
		} else {
			if r.Op == '>' {
				if ra[r.Var][1] > r.Val {
					rc := rangecp(ra)
					rc[r.Var] = [2]int{max(ra[r.Var][0], r.Val+1), ra[r.Var][1]}

					if r.Dest == "A" {
						n += rangel(rc)
					} else if r.Dest != "R" {
						n += runR(ws, rc, r.Dest)
					}
				}
				ra[r.Var] = [2]int{ra[r.Var][0], r.Val} // continue!
			} else { // know from part 1 its only <>
				if ra[r.Var][0] < r.Val {
					rc := rangecp(ra)
					rc[r.Var] = [2]int{ra[r.Var][0], min(ra[r.Var][1], r.Val-1)}

					if r.Dest == "A" {
						n += rangel(rc)
					} else if r.Dest != "R" {
						n += runR(ws, rc, r.Dest)
					}
				}
				ra[r.Var] = [2]int{r.Val, ra[r.Var][1]} // continue!
			}
		}
	}
	return n
}

func main() {
	ws := map[string]W{}
	s := bufio.NewScanner(os.Stdin)
	wsdone := false
	sum := 0
	for s.Scan() {
		l := s.Text()
		if l == "" {
			wsdone = true
			continue
		}
		if !wsdone {
			w := parseWorkflow(l)
			ws[w.Name] = w
		} else {
			v := parseVars(l)
			n, _ := run(ws, v)
			sum += n
			//fmt.Println(v, p, n)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// part2
	ra := map[string][2]int{
		"x": {1, 4000},
		"m": {1, 4000},
		"a": {1, 4000},
		"s": {1, 4000},
	}
	sum2 := runR(ws, ra, "in")

	fmt.Println(sum)
	fmt.Println(sum2)
}
