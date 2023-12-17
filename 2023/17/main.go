// Find most heat-loss eficient path for up to 3 consecutive straight moves and
// for no less than 4 but up to 10 consecutive straight moves.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type XY [2]int

// Priority queue item
type QI struct {
	XY  XY
	Dir XY
	L   int // length of straight block
	H   int // heat loss (cost)
}

// Priority queue
type Q struct {
	qis []QI
}

func (q Q) Len() int {
	return len(q.qis)
}

func (q *Q) Push(qi QI) {
	if q.qis == nil {
		q.qis = []QI{}
	}
	ii := sort.Search(len(q.qis), func(i int) bool {
		return q.qis[i].H < qi.H // make lowest last
	})
	q.qis = append(q.qis, QI{})
	copy(q.qis[ii+1:], q.qis[ii:])
	q.qis[ii] = qi
}

func (q *Q) Pop() QI {
	// sorted, just remove last
	qi := q.qis[len(q.qis)-1]
	q.qis = q.qis[:len(q.qis)-1]
	return qi
}

func pathf(m [][]int, minl int, maxl int) int {
	q := Q{}
	q.Push(QI{XY{0, 0}, XY{1, 0}, 1, 0}) // go right
	q.Push(QI{XY{0, 0}, XY{0, 1}, 1, 0}) // go down
	v := map[QI]bool{}

	for q.Len() > 0 {
		c := q.Pop()
		// for seen total cost is not important
		if _, seen := v[QI{c.XY, c.Dir, c.L, 0}]; seen {
			continue
		}
		v[QI{c.XY, c.Dir, c.L, 0}] = true

		// move
		xy := XY{c.XY[0] + c.Dir[0], c.XY[1] + c.Dir[1]}
		if xy[0] < 0 || xy[0] >= len(m[0]) || xy[1] < 0 || xy[1] >= len(m) {
			continue
		}
		h := c.H + m[xy[1]][xy[0]] // new heat loss

		// found end (for part 2 must be between minl and maxl)
		if c.L >= minl && c.L <= maxl && xy[0] == len(m[0])-1 && xy[1] == len(m)-1 {
			return h
		}

		// go to possible cells (forward, left and right turns)
		for _, d := range []XY{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			// no reverse
			if d[0]+c.Dir[0] == 0 && d[1]+c.Dir[1] == 0 {
				continue
			}
			// check if we can continue straight
			l := 1
			if d == c.Dir {
				l = c.L + 1
			}
			// if l > 3 {
			if (d != c.Dir && c.L < minl) || l > maxl { // part 2
				continue
			}
			q.Push(QI{xy, d, l, h})
		}
	}
	return 0
}

func main() {
	m := [][]int{}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ml := []int{}
		for _, c := range s.Text() {
			n, _ := strconv.Atoi(string(c))
			ml = append(ml, n)
		}
		m = append(m, ml)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(pathf(m, 1, 3))
	fmt.Println(pathf(m, 4, 10))
}
