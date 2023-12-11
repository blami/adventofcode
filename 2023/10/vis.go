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

// Render map to image.
func render(m [][]rune, a [][]rune, sp, cp XY, ap XY, s int) *image.Paletted {
	h := len(m)
	w := len(m[0])

	img := image.NewPaletted(image.Rect(0, 0, w*s, h*s), palette.Plan9)
	for y := range m {
		for x := range m[y] {
			xy := XY{x, y}
			c := color.RGBA{0, 0, 0, 0xff}
			if m[xy.Y][xy.X] != '.' {
				c = color.RGBA{92, 92, 92, 0xff}
			}
			if a[xy.Y][xy.X] == '#' {
				c = color.RGBA{0, 128, 0, 0xff}
			}
			if xy == cp {
				c = color.RGBA{0, 255, 0, 0xff}
			}
			if xy == sp {
				c = color.RGBA{255, 0, 0, 0xff}
			}
			if m[xy.Y][xy.X] == 'O' {
				c = color.RGBA{33, 33, 33, 0xff}
			}
			if m[xy.Y][xy.X] == 'I' {
				c = color.RGBA{128, 128, 0, 0xff}
			}
			if xy == ap {
				c = color.RGBA{255, 255, 0, 0xff}
			}
			// draw sxs rectangle so the image is "readable"
			for ry := 0; ry < s; ry++ {
				for rx := 0; rx < s; rx++ {
					img.Set((x*s)+rx, (y*s)+ry, c)
				}
			}
		}
	}

	return img
}

func saveGif(fn string, imgs []*image.Paletted) {
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
}
