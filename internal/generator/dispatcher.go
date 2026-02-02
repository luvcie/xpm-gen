package generator

import (
	"math"
	"math/rand"
	"xpm-gen/internal/config"
)

// allocates and populates the color grid based on configuration
// routes execution to specific algorithms or simulations
// takes: cfg (configuration struct)
// returns: 2d array of color indices
func GenerateGrid(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	for y := 0; y < cfg.Height; y++ {
		grid[y] = make([]int, cfg.Width)
	}

	// pre-calculate randomizers
	randX := rand.Intn(1000)
	randY := rand.Intn(1000)
	randColorOffset := rand.Intn(len(cfg.Colors))
	juliaCx := (rand.Float64() * 2.0) - 1.0
	juliaCy := (rand.Float64() * 2.0) - 1.0
	mandelZoom := 0.5 + rand.Float64()

	// stateful simulations
	if cfg.Algorithm == "melting" {
		return runMeltingSimulation(cfg)
	}
	if cfg.Algorithm == "creature" {
		return runCreatureGenerator(cfg)
	}
	if cfg.Algorithm == "cute" {
		return runCuteGenerator(cfg)
	}
	if cfg.Algorithm == "cutebunny" {
		return runCuteBunnyGenerator(cfg)
	}
	if cfg.Algorithm == "physarum" {
		return runPhysarum(cfg)
	}
	if cfg.Algorithm == "coral" {
		return runCoral(cfg)
	}
	if cfg.Algorithm == "attractor" {
		return runAttractor(cfg)
	}

	// stateless pixel-by-pixel generation
	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < cfg.Width; x++ {
			var colorIdx int
			switch cfg.Algorithm {
			case "noise":
				colorIdx = noise(cfg)
			case "xor":
				colorIdx = xorPattern(x, y, randX, randY, cfg)
			case "circles":
				colorIdx = circles(x, y, randX, randY, randColorOffset, cfg)
			case "pastel":
				colorIdx = pastel(x, y, randX, randY, cfg)
			case "mandelbrot":
				colorIdx = mandelbrot(x, y, cfg, mandelZoom, randColorOffset)
			case "julia":
				colorIdx = julia(x, y, cfg, juliaCx, juliaCy, randColorOffset)
			default:
				colorIdx = 0
			}
			grid[y][x] = colorIdx
		}
	}
	return grid
}

// GenerateFromExpression generates a grid using a custom Expression
func GenerateFromExpression(cfg config.Config, expr Expression) [][]int {
	grid := make([][]int, cfg.Height)
	w, h := float64(cfg.Width), float64(cfg.Height)
	
	for y := 0; y < cfg.Height; y++ {
		grid[y] = make([]int, cfg.Width)
		for x := 0; x < cfg.Width; x++ {
			// Evaluate expression
			val := expr.Eval(float64(x), float64(y), w, h)
			
			// Map result to color index
			// We use Mod and Abs to ensure it stays within bounds
			// Multiply by length to use the full range
			idx := int(math.Abs(val) * float64(len(cfg.Colors)))
			idx = idx % len(cfg.Colors)
			
			grid[y][x] = idx
		}
	}
	return grid
}