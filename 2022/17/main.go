// Read the input from stdin and simulate falling rocks. Find how tall will be
// tower of fallend rocks after 2022 and 1000000000000 rocks.

// TODO: Explore idea of only storing current height of each chamber column to
// make maxh() faster and less memory used by ch. Anything that is beyond
// highest block in each column is unnecessary...

package main

import (
	"bufio"
	"log"
	"os"
)

var rs = [][]complex64{
	[]complex64{0, 1, 2, 3},                // -
	[]complex64{1, 0 + 1i, 2 + 1i, 1 + 2i}, // +
	[]complex64{0, 1, 2, 2 + 1i, 2 + 2i},   // _|
	[]complex64{0, 0 + 1i, 0 + 2i, 0 + 3i}, // |
	[]complex64{0, 1, 0 + 1i, 1 + 1i},      // #
}

// Find maximum height in given chamber
func maxh(ch []complex64) float32 {
	h := float32(0)
	for i := range ch {
		if imag(ch[i]) > h {
			h = imag(ch[i])
		}
	}
	return h
}

// Check whether given rock fits at position p in given chamber.
func fits(ch []complex64, rock []complex64, p complex64) bool {
	for _, rp := range rock {
		q := p + rp
		in := false
		for i := range ch { // check if q is in ch
			if ch[i] == q {
				in = true
			}
		}
		if !(real(q) >= 0 && real(q) < 7) || imag(q) <= 0 || in {
			return false
		}
	}
	return true
}

// Return quotient and remainder of integer division.
func divmod(n, d int) (int, int) {
	return n / d, n % d
}

func main() {

	js := []byte{}        // jets
	ch := []complex64{-1} // chamber positions

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		for _, b := range s.Bytes() {
			js = append(js, b)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	ri := 0 // current rock index
	ji := 0 // current jet index
	n := 0  // block counter

	var th [2]int // tower heights

	cach := map[[2]int][2]int{} // repeating pattern cache
	maxn := 1000000000000
	for n <= maxn {
		h := maxh(ch)
		p := complex(2, 4+h) // start position
		if n == 2022 {
			th[0] = int(h) // part 1
		}
		// XXX: trick, there's a repeating pattern (alignment of jet directions
		// and rocks; this is generic way to find it in any input.
		if v, ok := cach[[2]int{ri, ji}]; ok {
			d, r := divmod(1000000000000-n, v[0]-n)
			if r == 0 {
				th[1] = int(h) + (v[1]-int(h))*d // part 2
				if n < 2022 {
					// just to make it play nicely with test.txt
					maxn = 2022
				}
				// break
			}
		} else {
			cach[[2]int{ri, ji}] = [2]int{n, int(h)}
		}
		// just to make it play nicely with test.txt
		if th[0] > 0 && th[1] > 0 {
			break
		}

		for { // let the rock hit the floor (or other rock)
			dx := complex64(-1 + 0i)
			if js[ji] == '>' {
				dx = 1 + 0i
			}
			ji = (ji + 1) % len(js) // next jet

			if fits(ch, rs[ri], p+dx) {
				p += dx
			} // sideways
			if fits(ch, rs[ri], p-1i) { // fall
				p -= 1i
			} else {
				break // rest
			}
		}

		for i := range rs[ri] {
			ch = append(ch, rs[ri][i]+p) // add rock to chamber
		}

		ri = (ri + 1) % len(rs) // next rock
		n++
	}

	log.Print(th[0])
	log.Print(th[1])
}
