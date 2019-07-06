package main

import (
	"fmt"
	mat "github.com/skelterjohn/go.matrix"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const width = 1024
const height = 1024
const alpha = 25.0
const beta = 37.0
const iterations = 150
const dist = 4.6
const fov = 39.0

var ca = math.Cos(alpha)
var sa = math.Sin(alpha)
var cb = math.Cos(beta)
var sb = math.Sin(beta)
var rotationMatrix = mat.MakeDenseMatrixStacked([][]float64{
	{ca * cb, -ca * sb, sa},
	{sb, cb, 0},
	{-cb * sa, sa * sb, ca},
})

var camCoordinates = mat.Product(rotationMatrix, mat.MakeDenseMatrixStacked([][]float64{{dist}, {0}, {0}}))
var nullVector = mat.Scaled(mat.MakeDenseMatrixStacked([][]float64{
	{ca * cb},
	{sb},
	{sa * cb},
}), -1)
var pixelSize = 2 * math.Tan(fov / 2) / height
var u = mat.Scaled(mat.Product(rotationMatrix, mat.MakeDenseMatrixStacked([][]float64{{0}, {0}, {1}})), pixelSize)
var v = mat.Scaled(mat.Product(rotationMatrix, mat.MakeDenseMatrixStacked([][]float64{{0}, {1}, {0}})), pixelSize)

func sphere(x, y, z float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2) + math.Pow(z, 2)) - 1.2
}

func cube(x, y, z float64) float64 {
	return math.Max(math.Abs(x), math.Max(math.Abs(y), math.Abs(z))) - 1
}

func sdf(x, y, z float64) float64 {
	return math.Min(sphere(x, y, z), cube(x, y, z))
	//return math.Max(-sphere(x, y, z), cube(x, y, z))
}

func rayDirection(x, y int) *mat.DenseMatrix {
	uv := mat.Sum(mat.Scaled(u, float64(x)), mat.Scaled(v, float64(y)))
	return mat.Sum(nullVector, uv)
}

func ray(x, y int, d float64) *mat.DenseMatrix{
	rayDir := rayDirection(x, y)
	normRayDir := mat.Scaled(rayDir, math.Pow(rayDir.OneNorm(), -1))
	return mat.Sum(camCoordinates, mat.Scaled(normRayDir, d))
}


func getPixel(x, y int) float64 {
	var cx, cy, cz = camCoordinates.Get(0, 0), camCoordinates.Get(1, 0), camCoordinates.Get(2, 0)
	var rx, ry, rz = rayDirection(x, y).Get(0, 0), rayDirection(x, y).Get(1, 0), rayDirection(x, y).Get(2, 0)
	var dist0 = sdf(cx, cy, cz)
	var k = dist0 + sdf(cx + rx * dist0, cy + ry * dist0, cz + rz * dist0)
	for i := 0; i < iterations; i++ {
		k += sdf(cx + rx * k, cy + ry * k, cz + rz * k)
	}
	return k
}

func main() {
	var result [width][height]float64
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var h2 = int(math.Round(width / 2))
	var w2 = int(math.Round(height / 2))

	var max float64 = 0;
	var min float64 = math.MaxFloat64;

	for y := -h2; y < h2; y++ {
		for x := -w2; x < w2; x++ {
			pixel := getPixel(y, x)
			if pixel < min {
				min = pixel
			}
			if pixel > max {
				max = pixel
			}
			result[y + h2][x + w2] = pixel

			img.Set(x + w2, y + h2,
				color.RGBA{uint8(result[y + h2][x + w2] * 100),
					uint8(result[y + h2][x + w2] * 100),
					uint8(result[y + h2][x + w2] * 100) ,
					255})
		}
	}
	fmt.Println("Max: %f\nMin: %f", max, min)
	f, err := os.OpenFile("img.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}