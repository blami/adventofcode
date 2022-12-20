package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ins(sl [][2]int, s [2]int, i int) [][2]int {
	if i < 0 {
		i = len(sl) + i
	}
	return append(sl[:i], append([][2]int{s}, sl[i:]...)...)
}

func main() {

	nn := [][2]int{} // value, order
	z := 0

	s := bufio.NewScanner(os.Stdin)
	i := 0
	for s.Scan() {
		v, _ := strconv.Atoi(s.Text())
		nn = append(nn, [2]int{v, i})
		if v == 0 { z = i }
		i++
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	for i:=0; i<len(nn); i++ {
		oi := -1 // find index to move
		for j := range nn {
			if nn[j][1] == i {
				oi = j
				break
			}
		}
		if oi == -1 { log.Fatal("xxx") }

		v := nn[oi]
		nn = append(nn[:oi], nn[oi+1:]...) // remove v 
		ni := (oi + v[0]) % len(nn) // new index
		if ni == 0 {
			ni = len(nn)
		}
		nn = ins(nn, v, ni)
	}

	co := 0
	for i := range nn {
		if nn[i][0] == 0 { z = i ; break }
	}
	for i := 1000; i <= 3000; i+=1000 {
		co += nn[(z + i) % len(nn)][0]
	}
	log.Print(co)

}
