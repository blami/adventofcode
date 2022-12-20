package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// Run one round of mix on slice nn where first item is value and second
// original position. Optionally use decryption key dk.
func mix(nn [][2]int, dk int) [][2]int {
	dkm := dk % (len(nn) - 1)
	for i := 0; i < len(nn); i++ {
		oi := -1 // find index to move
		for j := range nn {
			if nn[j][1] == i {
				oi = j
				break
			}
		}

		v := nn[oi]
		nn = append(nn[:oi], nn[oi+1:]...)  // remove v
		ni := (oi + (v[0] * dkm)) % len(nn) // new index
		if ni == 0 {
			ni = len(nn)
		}
		if ni < 0 {
			ni = len(nn) + ni // negative index is relative to tail
		}
		nn = append(nn[:ni], append([][2]int{v}, nn[ni:]...)...) // add v
	}

	return nn
}

func main() {

	in := [][2]int{} // input as tuple of value, index
	dk := 811589153

	s := bufio.NewScanner(os.Stdin)
	i := 0
	for s.Scan() {
		v, _ := strconv.Atoi(s.Text())
		in = append(in, [2]int{v, i})
		i++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	for _, p := range [][2]int{[2]int{1, 1}, [2]int{10, dk}} { // rounds, key
		// copy input for each part
		pin := make([][2]int, len(in))
		copy(pin, in)

		for r := 0; r < p[0]; r++ {
			pin = mix(pin, p[1])
		}
		z := -1
		for z = range pin {
			if pin[z][0] == 0 {
				break
			}
		}
		co := 0
		for i := 1000; i <= 3000; i += 1000 {
			co += (pin[(z+i)%len(pin)][0] * p[1])
		}
		log.Print(co)
	}
}
