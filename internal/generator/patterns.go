package generator

import (
	"math"
	"math/rand"
	"xpm-synth/internal/config"
)

// generates simple static noise
// takes: config
// returns: random color index
func noise(cfg config.Config) int {
	return rand.Intn(len(cfg.Colors))
}

// generates bitwise xor fractal pattern
// classic munching squares effect with coordinate offset
// takes: x/y coords, random offsets, config
// returns: color index
func xorPattern(x, y, randX, randY int, cfg config.Config) int {
	val := ((x + randX) ^ (y + randY))
	return val % len(cfg.Colors)
}

// generates hypnotic concentric ripples
// uses distance field with varying thickness and offsets
// takes: x/y coords, random offsets, config
// returns: color index
func circles(x, y, randX, randY, randOffset int, cfg config.Config) int {
	offsetX := (float64(randX%100) / 50.0) * float64(cfg.Width) * 0.5
	offsetY := (float64(randY%100) / 50.0) * float64(cfg.Height) * 0.5
	cx, cy := (float64(cfg.Width)/2)+offsetX, (float64(cfg.Height)/2)+offsetY
	ringThickness := 1.0 + (float64(randX%10) / 2.0)
	dist := math.Sqrt(math.Pow(float64(x)-cx, 2) + math.Pow(float64(y)-cy, 2))
	val := int(dist/ringThickness) + randOffset
	return val % len(cfg.Colors)
}

// generates domain-warped aesthetic textures
// uses sine wave interference to create shiny/glassy look
// takes: x/y coords, random offsets, config
// returns: color index
func pastel(x, y, randX, randY int, cfg config.Config) int {
	scale := 50.0
	dx := float64(x + randX)
	dy := float64(y + randY)
	warpX := dx + 20.0*math.Sin(dy/60.0)
	warpY := dy + 20.0*math.Cos(dx/60.0)
	h := 0.5 + 0.5*math.Sin((warpX+warpY)/scale)
	h = math.Pow(h, 0.8)
	colorIdx := int(h * float64(len(cfg.Colors)))
	if colorIdx < 0 {
		return 0
	}
	if colorIdx >= len(cfg.Colors) {
		return len(cfg.Colors) - 1
	}
	return colorIdx
}