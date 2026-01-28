package generator

import (
	"xpm-gen/internal/config"
)

// calculates pixel color for mandelbrot set
// applies zoom and random color offset for variety
// takes: x/y coords, config, zoom factor, random offset
// returns: color index
func mandelbrot(x, y int, cfg config.Config, zoom float64, randOffset int) int {
	scaleX := 3.5 / zoom
	scaleY := 2.0 / zoom
	offsetX := -2.5 + (1.0 - (1.0 / zoom))
	jx := float64(x)/float64(cfg.Width)*scaleX + offsetX
	jy := float64(y)/float64(cfg.Height)*scaleY - (scaleY / 2.0)

	zx, zy := 0.0, 0.0
	iter := 0
	maxIter := 50
	for iter < maxIter && (zx*zx+zy*zy) < 4.0 {
		tmp := zx*zx - zy*zy + jx
		zy = 2.0*zx*zy + jy
		zx = tmp
		iter++
	}
	if iter == maxIter {
		return 0
	}
	return ((iter + randOffset) % (len(cfg.Colors) - 1)) + 1
}

// calculates pixel color for julia set
// uses random complex constants cx/cy to vary shape
// takes: x/y coords, config, complex constants, random offset
// returns: color index
func julia(x, y int, cfg config.Config, cx, cy float64, randOffset int) int {
	zx := float64(x)/float64(cfg.Width)*3.0 - 1.5
	zy := float64(y)/float64(cfg.Height)*3.0 - 1.5
	iter := 0
	maxIter := 50
	for iter < maxIter && (zx*zx+zy*zy) < 4.0 {
		tmp := zx*zx - zy*zy + cx
		zy = 2.0*zx*zy + cy
		zx = tmp
		iter++
	}
	if iter == maxIter {
		return 0
	}
	return ((iter + randOffset) % (len(cfg.Colors) - 1)) + 1
}