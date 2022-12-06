// Read the input signal from stdin and find first packet and message marker.

package main

import (
	"bufio"
	"log"
	"os"
)

// Check given window for unique characters
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

	// NOTE: this algorithm is incorrect in case input will be bigger than
	// scanner buffer (there will be more s.Scan()s) and marker will be at
	// overlap of them as we use slice of the buffer itself as moving window.

	p := 0 // position in scanned buffer
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		buf := s.Bytes()
		for i, _ := range buf {
			// packet
			if fp == 0 && i < len(buf)-4 && check(buf[i:i+4]) {
				if fp == 0 {
					fp = p + i
				}
			}

			// message
			if fm == 0 && i < len(buf)-14 && check(buf[i:i+14]) {
				if fm == 0 {
					fm = p + i
				}
			}

			if fp != 0 && fm != 0 {
				break
			}
		}
		p = len(buf)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	// NOTE: add length of marker itself to get correct answer
	log.Print(fp + 4)
	log.Print(fm + 14)
}
