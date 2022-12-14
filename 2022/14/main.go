// Read input from stdin and build cave. Simulate sand flowing grain by grain
// in cave until it either falls below into abyss or fills up the cave.

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// XY alias
type XY [2]int

// Print cave current state.
func debug(c [][]byte) {
	// find nice boundaries
	l, r := 500, -1
	for i := range c[:len(c)-1] {
		for j := range c[i] {
			if c[i][j] != '.' {
				if j < l {
					l = j
				}
				if j > r {
					r = j
				}
			}
		}
	}

	log.Print()
	log.Print("cave size: ", len(c[0]), "x", len(c))
	for i := range c {
		log.Printf("%03d %s", i, string(c[i][l:r+1]))
	}
	log.Printf("    %03d", l)
}

// Add grain of sand to cave and let it fall until it rests. Takes slice to
// mutate as an argument and abyss flag which toggles between (true) stopping
// at moment when grain falls to abyss below rocks or until when there's no
// space for another grain of sand (false).
func sandfall(c [][]byte, abyss bool) bool {
	w := len(c[0])
	h := len(c)
	sp := XY{500, 0}

	// if there's no more space to add grain, return
	if c[sp[1]][sp[0]] != '+' {
		return false
	}

	rest := false
	for !rest {
		switch {
		case sp[1]+1 < h && c[sp[1]+1][sp[0]] == '.': // down
			sp = XY{sp[0], sp[1] + 1}
		case sp[1]+1 < h && sp[0]-1 >= 0 && c[sp[1]+1][sp[0]-1] == '.': // down left
			sp = XY{sp[0] - 1, sp[1] + 1}
		case sp[1]+1 < h && sp[0]+1 < w && c[sp[1]+1][sp[0]+1] == '.': // down right
			sp = XY{sp[0] + 1, sp[1] + 1}
		case sp[1]+1 >= h-2 && abyss: // abyss
			return false
		default: // rests where it is
			rest = true
		}
	}
	c[sp[1]][sp[0]] = 'o'
	return true
}

func main() {

	w, h := 500, 0 // sand flows from 500,0 so we need at least that size
	rl := [][]XY{} // rock lines

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		rl = append(rl, make([]XY, 0))

		l := strings.Split(s.Text(), " -> ")
		for _, li := range l {
			xy := strings.Split(li, ",")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			if x > w {
				w = x + 1
			}
			if y > h {
				h = y + 3
			} // in part 2 there's a floor on h + 2
			rl[len(rl)-1] = append(rl[len(rl)-1], XY{x, y})
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	w = w + h // need up to pyramid space for part 2

	// Build cave
	// TODO: function eh?
	c := make([][]byte, h)
	for i := range c { // initialize an empty cave
		c[i] = make([]byte, w)
		for j := range c[i] {
			ch := byte('.')
			if i == h-1 {
				ch = '#'
			}
			c[i][j] = ch
		}
	}
	for i := range rl { // add rocks
		for j := range rl[i] {
			if j == 0 {
				continue
			}
			a := 0 // axis
			switch {
			case rl[i][j-1][0] == rl[i][j][0]: // |
				a = 1
			case rl[i][j-1][1] == rl[i][j][1]: // -
				a = 0
			default:
				log.Fatal("wrong line")
			}
			d := 1                           // direction
			l := rl[i][j][a] - rl[i][j-1][a] // length
			if l < 0 {
				d = -1
				l *= d
			}
			ip := XY{rl[i][j-1][0], rl[i][j-1][1]}
			for l >= 0 {
				c[ip[1]][ip[0]] = '#'
				ip[a] += 1 * d
				l--
			}
		}
	}
	c[0][500] = '+' // sand pipe

	// Let the ~snow~ sand fall...
	na, nf := 0, 0
	for sandfall(c, true) { // part 1 find how much until abyss
		na++
	}
	debug(c)
	// remove sand from part 1
	for i := range c {
		for j := range c[i] {
			if c[i][j] == 'o' {
				c[i][j] = '.'
			}
		}
	}
	for sandfall(c, false) { // part 2 find how many fill up space
		nf++
	}
	debug(c)

	log.Print(na)
	log.Print(nf)
}
