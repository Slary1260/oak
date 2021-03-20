package render

import (
	"image"
	"image/color"
	"testing"
)

func TestSimpleCircle(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))
	DrawCircle(rgba, color.RGBA{255, 255, 255, 255}, 5, 0, 0, 0)
	circleRGBACheck(t, rgba)
	c := NewCircle(color.RGBA{255, 255, 255, 255}, 5, 0, 0, 0)
	circleRGBACheck(t, c.GetRGBA())
}

func circleRGBACheck(t *testing.T, rgba *image.RGBA) {
	// For better or for worse, the current implementation produces
	// . . . . . . . . . .
	// . . . x x x x x . .
	// . . x x       x x .
	// . x x           x x
	// . x               x
	// . x               x
	// . x               x
	// . x x           x x
	// . . x x       x x .
	// . . . x x x x x . .
	// This should change in the future, probably leaning towards using Beziers.
	// 3 years later thoughts: Circle using a circle algorithm and bezier using a
	// bezier algorithm is perfectly fine, the generated circle is fine
	boolExpected := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
		{0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
	}
	for x, col := range boolExpected {
		for y, b := range col {
			if b == 0 {
				if (color.RGBA{0, 0, 0, 0}) != rgba.At(x, y) {
					t.Fatalf("circle not unset where expected")
				}
			} else {
				if (color.RGBA{255, 255, 255, 255}) != rgba.At(x, y) {
					t.Fatalf("circle not set where expected")
				}
			}
		}
	}
}
