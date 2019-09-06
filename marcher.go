package main

import (
	mat "github.com/skelterjohn/go.matrix"
	"math"
)

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
	//return math.Min(sphere(x, y, z), cube(x, y, z))
	return math.Max(-sphere(x, y, z), cube(x, y, z))
}


func rayDirection(x, y int) *mat.DenseMatrix {
	uv := mat.Sum(mat.Scaled(u, float64(x)), mat.Scaled(v, float64(y)))
	return mat.Sum(nullVector, uv)
}


type Pixel struct {
	p float64
	cx, cy, cz float64
	rx, ry, rz float64
}

func (p *Pixel) Iterate() {
	p.p += sdf(p.cx + p.rx * p.p, p.cy + p.ry * p.p, p.cz + p.rz * p.p)
}

func (p *Pixel) Color() uint8 {
	if p.p > maxDist {
		return 0
	}
	return uint8(p.p * 100)
}

func InitMatrix(width, height int) [][]Pixel {
	result := make([][]Pixel, 0)
	var cx, cy, cz = camCoordinates.Get(0, 0), camCoordinates.Get(1, 0), camCoordinates.Get(2, 0)
	var dist0 = sdf(cx, cy, cz)
	for h := 0; h < height; h++ {
		row := make([]Pixel, 0)
		for w := 0; w < width; w++ {
			x := w - width/2
			y := h - height/2
			var rx, ry, rz = rayDirection(x, y).Get(0, 0), rayDirection(x, y).Get(1, 0), rayDirection(x, y).Get(2, 0)
			row = append(row, Pixel{dist0, cx, cy, cz, rx, ry, rz})
		}
		result = append(result, row)
	}
	return result
}