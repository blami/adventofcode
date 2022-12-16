// Visualization helper.

package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"log"
	"os"
)

// Global store
var imgs []*image.Paletted

// Render cave to image.
// BUG: Nice boundaries won't work with part 2 as it will stretch them on the
// go; need to figure better way.
func render(c [][]byte) {
	const s = 3

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

	h := len(c)
	w := r - l + 1

	img := image.NewPaletted(image.Rect(0, 0, w*s, h*s), palette.Plan9)
	for y := range c {
		for x, cc := range c[y][l : r+1] {
			// decide color, default is black 'cause cave
			cl := palette.Plan9[0].(color.RGBA)
			switch cc {
			case '#': // rocks
				cl = palette.Plan9[153].(color.RGBA)
			case 'o': // grains of sand
				cl = palette.Plan9[186].(color.RGBA)
			case '+': // sand pipe
				cl = palette.Plan9[222].(color.RGBA)
			}
			for ry := 0; ry < s; ry++ {
				for rx := 0; rx < s; rx++ {
					img.Set((x*s)+rx, (y*s)+ry, cl)
				}
			}
		}
	}

	imgs = append(imgs, img)
}

func clearGif() {
	imgs = make([]*image.Paletted, 0)
}

func saveGif(fn string) {
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	delay := []int{}
	for range imgs {
		delay = append(delay, 1) // 1/100s
	}

	err := gif.EncodeAll(f, &gif.GIF{
		Image:     imgs,
		Delay:     delay,
		LoopCount: 1, // will actually loop 2x; thanks Netscape
	})
	if err != nil {
		log.Println("len", len(imgs))
		log.Fatal(err)
	}
}
