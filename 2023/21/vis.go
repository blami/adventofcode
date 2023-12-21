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

// Render map to image.
func renderGif(m [][]rune, h []XY, hn []XY, s int) {
	img := image.NewPaletted(image.Rect(0, 0, len(m[0])*s, len(m)*s), palette.Plan9)
	for y := range m {
		for x := range m[y] {
			c := color.RGBA{0, 0, 0, 0xff}
			if m[y][x] != '.' {
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
			// draw sxs rectangle so the image is "readable"
			for ry := 0; ry < s; ry++ {
				for rx := 0; rx < s; rx++ {
					img.Set((x*s)+rx, (y*s)+ry, c)
				}
			}
		}
	}

	imgs = append(imgs, img)
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
		log.Fatal(err)
	}
	imgs = nil
}
