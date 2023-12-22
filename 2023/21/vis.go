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

var imgs []*image.Paletted

// Render map to image. Argument wr is wrapping range, how many full maps to 0
// left, 1 right, 2 up, 3 bottom should be rendered.
func renderGif(m [][]rune, h []XY, hn []XY, s int, ra [4]int) {
	if os.Getenv("DEBUG") == "" {
		return
	}
	ax := 0 - ra[0]*len(m[0])
	aw := ra[1] * len(m[0])
	ay := 0 - ra[2]*len(m)
	ah := ra[3] * len(m)

	// adjust in-image width and height so they are always positive
	iw := (ra[0] + ra[1]) * len(m[0])
	ih := (ra[2] + ra[3]) * len(m)
	// offset
	ox := ra[0] * len(m[0])
	oy := ra[2] * len(m)

	//fmt.Println("size", aw,ah, "xy", ax, ay, "isize", iw, ih, "oxy", ox, oy)
	img := image.NewPaletted(image.Rect(0, 0, iw*s, ih*s), palette.Plan9)

	for y := ay; y < ah; y++ {
		for x := ax; x < aw; x++ {
			c := color.RGBA{0, 0, 0, 0xff}
			if m[tran(y, len(m))][tran(x, len(m[0]))] != '.' {
				c = color.RGBA{92, 92, 92, 0xff}
			}
			for _, xy := range h {
				if xy[0] == x && xy[1] == y {
					c = color.RGBA{128, 0, 0, 0xff}
				}
			}
			for _, xy := range hn {
				if xy[0] == x && xy[1] == y {
					c = color.RGBA{255, 0, 0, 0xff}
				}
			}
			// borders
			if (aw > len(m[0]) && x%len(m[0]) == 0) || (ah > len(m) && y%len(m) == 0) {
				c = color.RGBA{0, 255, 0, 0xff}
			}
			// draw sxs rectangle so the image is "readable"
			for ry := 0; ry < s; ry++ {
				for rx := 0; rx < s; rx++ {
					img.Set((ox+x)*s+rx, (oy+y)*s+ry, c)
				}
			}
		}
	}

	imgs = append(imgs, img)
}

func saveGif(fn string, empty bool) {
	if os.Getenv("DEBUG") == "" {
		return
	}
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
		log.Fatal(err)
	}
	if empty {
		imgs = nil
	}
}
