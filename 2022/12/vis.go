// Visualization helper.

package main

import (
	"os"
	"log"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
)

// Render map to image.
func render(m [][]int, v map[XY]bool, sp, ep, cp XY) *image.Paletted {
	h := len(m)
	w := len(m[0])
	s := 5

	img := image.NewPaletted(image.Rect(0, 0, w*s, h*s), palette.Plan9)
	for y, l := range m {
		for x, _ := range l {
			xy := XY{x, y}
			h := uint8((255/26) * m[y][x])
			c := color.RGBA{255-h, 255-h, 255-h, 0xff}
			if _, ok := v[XY{x, y}]; ok {
				c = color.RGBA{0, 255-h, 0, 0xff}
			}
			if xy == cp {
				c = color.RGBA{0, 0, 0xff, 0xff}
			}
			if xy == sp || xy == ep {
				c = color.RGBA{0xff, 0, 0, 0xff}
			}
			// draw sxs rectangle so the image is "readable"
			for ry := 0; ry < s; ry++ {
				for rx := 0; rx < s; rx++ {
					img.Set((x*s) + rx, (y*s) + ry, c)
				}
			}
		}
	}

	return img
}

func saveGif(fn string, imgs []*image.Paletted) {
	f, _ := os.OpenFile(fn, os.O_WRONLY | os.O_CREATE, 0644)
	defer f.Close()

	delay := []int{}
	for range(imgs) {
		delay = append(delay, 1) // 1/100s
	}

	err := gif.EncodeAll(f, &gif.GIF{
		Image: imgs,
		Delay: delay,
		LoopCount: 1, // will actually loop 2x; thanks Netscape
	})
	if err != nil {
		log.Fatal(err)
	}
}
