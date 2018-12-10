package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
)

func main() {
	pp, err := readInput("../inputs/10.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	/*
		for i := 0; i < 100000; i++ {
			bounds, offset := computeBounds(pp, i)
			fmt.Printf("%d %d bounds x: %d y: %d offset x: %d y: %d\n", bounds.Max.X*bounds.Max.Y, i, bounds.Max.X, bounds.Max.Y, offset.x, offset.y)
		}
	*/

	i := 10813
	bounds, offset := computeBounds(pp, i)
	fmt.Printf("%d %d bounds x: %d y: %d offset x: %d y: %d\n", bounds.Max.X*bounds.Max.Y, i, bounds.Max.X, bounds.Max.Y, offset.x, offset.y)
	img := render(bounds, offset, pp, i)
	for _, row := range img {
		fmt.Println(string(row))
	}
}

type point struct {
	x, y   int
	dx, dy int
}

func computeBounds(pp []point, t int) (image.Rectangle, point) {
	var (
		bounds *image.Rectangle
	)
	for _, p := range pp {
		cp := point{x: p.x + t*p.dx, y: p.y + t*p.dy}
		if bounds == nil {
			bounds = &image.Rectangle{
				Min: image.Point{X: cp.x, Y: cp.y},
				Max: image.Point{X: cp.x, Y: cp.y},
			}
		} else {
			if cp.x < bounds.Min.X {
				bounds.Min.X = cp.x
			}
			if cp.x > bounds.Max.X {
				bounds.Max.X = cp.x
			}

			if cp.y < bounds.Min.Y {
				bounds.Min.Y = cp.y
			}
			if cp.y > bounds.Max.Y {
				bounds.Max.Y = cp.y
			}
		}
	}

	offX := -bounds.Min.X
	bounds.Min.X += offX
	bounds.Max.X += offX
	offY := -bounds.Min.Y
	bounds.Min.Y += offY
	bounds.Max.Y += offY

	return *bounds, point{x: offX, y: offY}
}

func render(bounds image.Rectangle, offset point, pp []point, t int) [][]rune {
	img := make([][]rune, bounds.Max.Y+1)
	for y := 0; y <= bounds.Max.Y; y++ {
		img[y] = make([]rune, bounds.Max.X+1)
		for x := 0; x <= bounds.Max.X; x++ {
			img[y][x] = ' '
		}
	}

	for _, p := range pp {
		y := p.y + t*p.dy + offset.y
		x := p.x + t*p.dx + offset.x
		img[y][x] = '#'
	}

	return img
}

func readInput(filename string) ([]point, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %v", filename, err)
	}

	var pp []point
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		l := scan.Text()
		// position=<-21400,  54189> velocity=< 2, -5>
		var p point
		_, err := fmt.Sscanf(l, "position=<%d, %d> velocity=<%d, %d>", &p.x, &p.y, &p.dx, &p.dy)
		if err != nil {
			return nil, fmt.Errorf("error parsing line '%s': %v", l, err)
		}
		pp = append(pp, p)
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}

	return pp, nil
}
