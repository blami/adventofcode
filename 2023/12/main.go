// Find sum of valid arrangements of working and non-working springs according
// to given stencil. Also do it for 5 times unfold map.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var dp map[[3]int]int = map[[3]int]int{}

// Dynamic programming, li is position in l, ni is position in n and chl is
// length of current # chain.
func f(l []rune, n []int, li int, ni int, chl int) int {
	r := 0

	k := [3]int{li, ni, chl}
	r, ok := dp[k]
	if ok == true {
		return r
	}

	if li == len(l) {
		if ni == len(n) && chl == 0 {
			return 1
		} else if ni == len(n)-1 && n[ni] == chl {
			return 1
		} else {
			return 0
		}
	}

	for _, c := range []rune{'.', '#'} {
		if l[li] == c || l[li] == '?' {
			if c == '.' && chl == 0 {
				r += f(l, n, li+1, ni, 0)
			} else if c == '.' && chl > 0 && ni < len(n) && n[ni] == chl {
				r += f(l, n, li+1, ni+1, 0)
			} else if c == '#' {
				r += f(l, n, li+1, ni, chl+1)
			}
		}
	}

	dp[k] = r
	return r
}

func unfold(l []rune, n []int, t int) ([]rune, []int) {
	nl := []rune{}
	nn := []int{}
	for i := 0; i < t; i++ {
		nn = append(nn, n...)
		nl = append(nl, l...)
		if i != t-1 {
			// join by ?
			nl = append(nl, '?')
		}
	}
	return nl, nn
}

func main() {
	sum := 0
	sum2 := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := strings.Split(s.Text(), " ")

		// parse numbers
		n := []int{}
		for _, v := range strings.Split(l[1], ",") {
			i, _ := strconv.Atoi(v)
			n = append(n, i)
		}

		// part 1
		dp = map[[3]int]int{}
		sum += f([]rune(l[0]), n, 0, 0, 0)

		// part 2
		dp = map[[3]int]int{}
		ul, un := unfold([]rune(l[0]), n, 5)
		sum2 += f(ul, un, 0, 0, 0)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
