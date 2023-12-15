// Find sum of step hashes and sum of lens powers.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func hash(s string) int {
	cur := 0
	for _, c := range s {
		cur += int(c)
		cur *= 17
		cur %= 256
	}
	return cur
}

type Lens struct {
	Label string
	FL    int
}

func debug(step string, bi int, ls [256][]Lens) {
	fmt.Println("after ", step, "bi", bi)
	for i, l := range ls {
		if len(l) > 0 {
			fmt.Println("  box", i, l)
		}
	}
}

func main() {
	sum := 0
	box := [256][]Lens{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		for _, step := range strings.Split(s.Text(), ",") {
			sum += hash(step)
			del := strings.Index(step, "=")
			if del == -1 {
				del = strings.Index(step, "-")
			}
			label := step[0:del]

			bi := hash(label)
			//fmt.Print(step, " box ", bi)
			switch step[del] {
			case '-':
				for li := range box[bi] {
					// removal
					if box[bi][li].Label == label {
						//fmt.Print(" removal ", box[bi], "->")
						box[bi] = append(box[bi][:li], box[bi][li+1:]...)
						break
					}
				}
			case '=':
				fls := step[del+1:]
				fl, _ := strconv.Atoi(fls)

				l := Lens{label, fl}
				// swap or append
				add := true
				for li := range box[bi] {
					if box[bi][li].Label == label {
						//fmt.Print(" swap ", box[bi], "->")
						box[bi][li] = l
						add = false
						break
					}
				}
				if add {
					//fmt.Print(" append ", box[bi], "->")
					box[bi] = append(box[bi], l)
				}
			}
			//fmt.Println(box[bi])
			//debug(step, bi, box)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// find power
	sum2 := 0
	for bi := range box {
		for li, l := range box[bi] {
			pow := (bi + 1) * (li + 1) * l.FL
			//fmt.Println(l.Label, "bi", bi, "li", li, "FL", l.FL, "pow", pow)
			sum2 += pow
		}
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
