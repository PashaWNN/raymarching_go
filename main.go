package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"runtime"
)

const width = 512
const height = 512
const alpha = 25.0
const beta = 37.0
const iterations = 100
const dist = 4.6
const maxDist = 100
const fov = 39.0

var threads = runtime.NumCPU()


func getPixel(x, y int) float64 {
	var cx, cy, cz = camCoordinates.Get(0, 0), camCoordinates.Get(1, 0), camCoordinates.Get(2, 0)
	var rx, ry, rz = rayDirection(x, y).Get(0, 0), rayDirection(x, y).Get(1, 0), rayDirection(x, y).Get(2, 0)
	var dist0 = sdf(cx, cy, cz)
	var k = dist0 + sdf(cx + rx * dist0, cy + ry * dist0, cz + rz * dist0)
	for i := 0; i < iterations; i++ {
		k = iteratePixel(k, cx, cy, cz, rx, ry, rz)
	}
	if k > maxDist {
		return 0
	}
	return k
}

func iteratePixel(p float64, cx, cy, cz, rx, ry, rz float64) float64 {
	return p + sdf(cx + rx * p, cy + ry * p, cz + rz * p)
}

func main() {
	result := InitMatrix(width, height)
	palette := make([]color.Color, 256)
	for i := 0; i < 256; i++ {
		palette[i] = color.Gray{255 - uint8(i)}
	}
	rect := image.Rect(0, 0, width, height)
	anim := gif.GIF{}
	row := make(chan int, threads)
	chunk := height / threads
    for iter := 0; iter < iterations; iter++ {
    	img := image.NewPaletted(rect, palette)
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 1)
		for i := 0; i < threads; i++ {
			go func(start int) {
				end := start + chunk
				if end > height {
					end = height
				}
				for y := start; y < end; y++ {
					for x := 0; x < width; x++ {
						result[y][x].Iterate()
						col := result[y][x].Color()
						//img.Set(x, y, color.RGBA{col, col, col, 255})
						img.SetColorIndex(x, y, col)
					}
				}
				row <- 1
			}(i * chunk)
		}

		for i := 0; i < threads; i++ {
			<-row
		}
	}
	f, err := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_ = gif.EncodeAll(f, &anim)
}
