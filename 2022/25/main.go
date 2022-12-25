package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// Integer power.
func pow(b, e int) int {
	v := 1
	for e > 0 {
		v *= b
		e--
	}
	return v
}

// Convert SNAFU to int.
func stoi(s string) int {
	v := 0
	for i := len(s)-1; i >= 0; i-- {
		vv := 0
		switch s[i] {
		case '-':
			vv = -1
		case '=':
			vv = -2
		case '0','1','2':
			vv, _ = strconv.Atoi(string(s[i]))
		default:
			log.Fatalf("wrong SNAFU: %s, digit: %c", s, s[i])
		}
		v += vv * pow(5, (len(s)-1) - i)
	}
	return v
}

// Convert int to SNAFU.
func itos(i int) string {
	s := ""
	for i > 0 {
		s = []string{"=","-","0","1","2"}[(i + 2) % 5] + s
		i = (i + 2) / 5
	}
	return s
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	t := 0 // total
	for s.Scan() {
		l := s.Text()
		t += stoi(l)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	log.Print(itos(t))
}
