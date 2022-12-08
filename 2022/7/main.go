// Read input from stdin, build directory tree based on captured output and
// find the smallest directory sizes and smallest enough directory to remove.

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// A directory.
type Dir struct {
	Name  string // directory name
	FSize int    // file size
	Dirs  []*Dir // subdirectories
	Par   *Dir   // parent directory
}

// Change directory to n and return pointer to it.
func (d *Dir) Cd(n string) *Dir {
	//if d.Dirs == nil { return nil } // no subdirectories
	for _, sd := range d.Dirs {
		if sd.Name == n {
			return sd
		}
	}
	return nil
}

// Make a directory n in current directory
func (d *Dir) Mkdir(n string) {
	d.Dirs = append(d.Dirs, &Dir{Name: n, Par: d})
}

// Add file size to directory
func (d *Dir) Grow(s int) {
	d.FSize += s
}

/*
// Print full path to directory for debugging purposes
func (d* Dir) Path() string {
	fp := d.Name
	p := d.Par
	for p != nil {
		fp = p.Name + "/" + fp
		p = p.Par
	}
	return fp
}
*/

// Recursively measure size of directory as sum of its own size and sizes of
// all nested subdirectories.
func (d *Dir) Size() int {
	ts := d.FSize
	for _, sd := range d.Dirs {
		ts += sd.Size()
	}
	return ts
}

// Recursively walk directory down and execute f on it.
func (d *Dir) Walk(f func(dd *Dir)) {
	f(d)
	for _, sd := range d.Dirs {
		sd.Walk(f)
	}
}

func main() {

	r := Dir{}
	c := &r // current directory

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		l := strings.Split(s.Text(), " ")
		switch {
		// cd
		case l[0] == "$" && l[1] == "cd":
			switch {
			case l[2] == "/":
				c = &r
			case l[2] == "..":
				if c.Par == nil {
					log.Fatal("nothing above ", c.Name)
				}
				c = c.Par
			default:
				c = c.Cd(l[2])
				if c == nil {
					log.Fatal("invalid dir")
				}
			}
		// directory
		case l[0] == "dir":
			c.Mkdir(l[1])
		// file size
		case l[0] != "$" && l[0] != "dir":
			fsz, _ := strconv.Atoi(l[0])
			c.Grow(fsz)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	totsz := 0 // total size of directories smaller than 100000

	rmdsz := r.Size()          // smallest directory to remove size found
	remain := 70000000 - rmdsz // remaining free space

	r.Walk(func(d *Dir) {
		if d.Size() <= 100000 {
			totsz += d.Size()
		}

		dsz := d.Size()
		if remain+dsz >= 30000000 {
			if rmdsz > dsz {
				rmdsz = dsz
			}
		}
	})

	log.Print(totsz)
	log.Print(rmdsz)

}
