// In given input find first and last digit (if only the first digit is present
// it is also the last digit) on each line and form a two digit number from
// them. Find sum of all those numbers. Second part is parsing spelled out
// numbers too.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	sum := 0
	sum2 := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := s.Text()

		flag := 0
		flag2 := 0
		var fd, ld, fd2, ld2 int

		for i := 0; i < len(l); i++ {
			var v, v2 int

			if l[i] >= 47 && l[i] <= 58 {
				v = int(l[i]) - 48
				// Second part
				// Spelled out numbers can't be just replaced because we
				// don't know order.
				v2 = int(l[i]) - 48
			} else if i+2 < len(l) && l[i:i+3] == "one" {
				v2 = 1
				i += 2
			} else if i+2 < len(l) && l[i:i+3] == "two" {
				v2 = 2
				i += 2
			} else if i+4 < len(l) && l[i:i+5] == "three" {
				v2 = 3
				i += 4
			} else if i+3 < len(l) && l[i:i+4] == "four" {
				v2 = 4
				i += 3
			} else if i+3 < len(l) && l[i:i+4] == "five" {
				v2 = 5
				i += 3
			} else if i+2 < len(l) && l[i:i+3] == "six" {
				v2 = 6
				i += 2
			} else if i+4 < len(l) && l[i:i+5] == "seven" {
				v2 = 7
				i += 4
			} else if i+4 < len(l) && l[i:i+5] == "eight" {
				v2 = 8
				i += 4
			} else if i+3 < len(l) && l[i:i+4] == "nine" {
				v2 = 9
				i += 3
			} else {
				continue
			}

			// Need v != 0 condition because we do both part 1 and part 2
			if v != 0 && flag == 0 {
				fd = v
				flag = 1
			} else if v != 0 {
				ld = v
				flag = 2
			}
			if v2 != 0 && flag2 == 0 {
				fd2 = v2
				flag2 = 1
			} else if v2 != 0 {
				ld2 = v2
				flag2 = 2
			}

			//fmt.Printf("  i=%d flag=%d flag2=%d v=%d v2=%d fdld=%d%d fd2ld2=%d%d\n", i, flag, flag2, v, v2, fd, ld, fd2, ld2)
		}
		if flag == 1 {
			ld = fd
		}
		if flag2 == 1 {
			ld2 = fd2
		}

		n := fd*10 + ld
		n2 := fd2*10 + ld2
		//fmt.Printf("  n=%d n2=%d\n", n, n2)

		sum += n
		sum2 += n2
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
	fmt.Println(sum2)
}
