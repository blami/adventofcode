// Find sums of extrapolated history numbers on right and left side of dataset.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	sum := 0
	sum2 := 0
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		r := 0 // round
		dat := [][]int{}

		l := strings.Fields(s.Text())
		dat = append(dat, []int{})
		for _, v := range l {
			n, _ := strconv.Atoi(v)
			dat[0] = append(dat[0], n)
		}

		for {
			r++
			dat = append(dat, []int{})
			z := 0
			for i := 0; i < len(dat[r-1])-1; i++ {
				v := dat[r-1][i+1] - dat[r-1][i]
				dat[r] = append(dat[r], v)
				if v == 0 {
					z++
				}
			}
			if z == len(dat[r]) {
				break
			}
		}

		// extrapolate
		for i := len(dat) - 1; i >= 0; i-- {
			suf := 0
			pre := 0
			if i < len(dat)-1 {
				suf = dat[i][len(dat[i])-1] + dat[i+1][len(dat[i+1])-1]
				pre = dat[i][0] - dat[i+1][0]
			}
			dat[i] = append(dat[i][:1], dat[i][1:]...)
			dat[i][0] = pre
			dat[i] = append(dat[i], suf)
		}
		sum += dat[0][len(dat[0])-1]
		sum2 += dat[0][0]
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
