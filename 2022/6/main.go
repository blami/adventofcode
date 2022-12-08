// Read input from stdin and find the first packet and message marker in
// signal.

package main

import (
	"bufio"
	"log"
	"os"
)

// Check given window for unique characters.
func check(win []byte) bool {
	seen := make(map[byte]bool)
	for _, b := range win {
		if _, ok := seen[b]; ok {
			return false
		}
		seen[b] = true
	}
	return true
}

func main() {

	fp := 0 // position of first packet marker (begin)
	fm := 0 // position of first message marker (begin)

	// NOTE: This assumes the entire signal can be read in single s.Scan()
	// which worked for the excercise. If it was not case transition between
	// buf needs to be handled explicitely.
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	buf := s.Bytes()
	for i, _ := range buf {
		// packet
		if fp == 0 && i < len(buf)-4 && check(buf[i:i+4]) {
			if fp == 0 {
				fp = i
			}
		}
		// message
		if fm == 0 && i < len(buf)-14 && check(buf[i:i+14]) {
			if fm == 0 {
				fm = i
			}
		}
		if fp != 0 && fm != 0 {
			break
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// NOTE: add length of marker itself to get correct answer
	log.Print(fp + 4)
	log.Print(fm + 14)
}
