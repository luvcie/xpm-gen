package generator

import (
	"math"
	"math/rand"
	"xpm-gen/internal/config"
)

// executes cyclic cellular automaton simulation
// evolves a random grid over generations to create liquid patterns
// takes: config
// returns: full 2d grid of color indices
func runMeltingSimulation(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	nextGrid := make([][]int, cfg.Height)
	for y := 0; y < cfg.Height; y++ {
		grid[y] = make([]int, cfg.Width)
		nextGrid[y] = make([]int, cfg.Width)
		for x := 0; x < cfg.Width; x++ {
			grid[y][x] = rand.Intn(len(cfg.Colors))
		}
	}

	generations := 50 + rand.Intn(100)
	threshold := 1

	for g := 0; g < generations; g++ {
		for y := 0; y < cfg.Height; y++ {
			for x := 0; x < cfg.Width; x++ {
				currentVal := grid[y][x]
				nextVal := (currentVal + 1) % len(cfg.Colors)
				neighbors := 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						ny := (y + dy + cfg.Height) % cfg.Height
						nx := (x + dx + cfg.Width) % cfg.Width
						if grid[ny][nx] == nextVal {
							neighbors++
						}
					}
				}
				if neighbors >= threshold {
					nextGrid[y][x] = nextVal
				} else {
					nextGrid[y][x] = currentVal
				}
			}
		}
		for y := 0; y < cfg.Height; y++ {
			copy(grid[y], nextGrid[y])
		}
	}
	return grid
}

// generates symmetric rorschach-style creatures
// uses random walkers, gravity simulation, and mirroring
// takes: config
// returns: full 2d grid of color indices
func runCreatureGenerator(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	for y := 0; y < cfg.Height; y++ {
		grid[y] = make([]int, cfg.Width)
		for x := 0; x < cfg.Width; x++ {
			grid[y][x] = 0
		}
	}

	centerX := cfg.Width / 2
	blobs := 5 + rand.Intn(10)

	for i := 0; i < blobs; i++ {
		cx := centerX + (rand.Intn(20) - 10)
		cy := rand.Intn(cfg.Height-20) + 10
		radius := 5 + rand.Intn(20)
		colorType := 1 + rand.Intn(3)

		for y := 0; y < cfg.Height; y++ {
			for x := 0; x < centerX; x++ {
				dx := x - cx
				dy := y - cy
				dist := math.Sqrt(float64(dx*dx + dy*dy))
				noise := rand.Float64() * 5.0
				if dist < (float64(radius) + noise) {
					grid[y][x] = colorType
				}
			}
		}
	}

	for i := 0; i < 500; i++ {
		x := rand.Intn(centerX)
		y := rand.Intn(cfg.Height - 10)
		if grid[y][x] != 0 {
			length := rand.Intn(20)
			for d := 0; d < length; d++ {
				if y+d < cfg.Height {
					grid[y+d][x] = grid[y][x]
				}
			}
		}
	}

	numEyes := 1 + rand.Intn(3)
	for i := 0; i < numEyes; i++ {
		ex := rand.Intn(centerX - 5)
		ey := rand.Intn(cfg.Height/2) + 10
		if grid[ey][ex] != 0 {
			grid[ey][ex] = 5
			grid[ey][ex+1] = 5
			grid[ey+1][ex] = 5
			grid[ey+1][ex+1] = 5
		}
	}

	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < centerX; x++ {
			mirrorX := cfg.Width - 1 - x
			grid[y][mirrorX] = grid[y][x]
		}
	}

	return grid
}

// simulates clifford attractor with density mapping
// searches for chaotic parameters to ensure good spread
// takes: config
// returns: full 2d grid of color indices
func runAttractor(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	for y := 0; y < cfg.Height; y++ {
		grid[y] = make([]int, cfg.Width)
	}

	density := make([][]float64, cfg.Height)
	for y := 0; y < cfg.Height; y++ {
		density[y] = make([]float64, cfg.Width)
	}

	var a, b, c, d float64
	foundGoodParams := false
	for attempt := 0; attempt < 100; attempt++ {
		a = rand.Float64()*4.0 - 2.0
		b = rand.Float64()*4.0 - 2.0
		c = rand.Float64()*4.0 - 2.0
		d = rand.Float64()*4.0 - 2.0

		x, y := 0.0, 0.0
		minX, maxX := 10.0, -10.0
		minY, maxY := 10.0, -10.0
		goodSpread := true
		for i := 0; i < 1000; i++ {
			xn := math.Sin(a*y) + c*math.Cos(a*x)
			yn := math.Sin(b*x) + d*math.Cos(b*y)
			x, y = xn, yn
			if i > 100 {
				if x < minX { minX = x }
				if x > maxX { maxX = x }
				if y < minY { minY = y }
				if y > maxY { maxY = y }
			}
			if math.IsInf(x, 0) || math.IsNaN(x) {
				goodSpread = false
				break
			}
		}
		if !goodSpread { continue }
		if (maxX-minX) > 0.5 && (maxY-minY) > 0.5 {
			foundGoodParams = true
			break
		}
	}
	if !foundGoodParams {
		a, b, c, d = -1.4, 1.6, 1.0, 0.7
	}

	x, y := 0.0, 0.0
	iterations := 5000000
	for i := 0; i < iterations; i++ {
		xn := math.Sin(a*y) + c*math.Cos(a*x)
		yn := math.Sin(b*x) + d*math.Cos(b*y)
		x, y = xn, yn
		screenX := int((x + 2.5) / 5.0 * float64(cfg.Width))
		screenY := int((y + 2.5) / 5.0 * float64(cfg.Height))
		if screenX >= 0 && screenX < cfg.Width && screenY >= 0 && screenY < cfg.Height {
			density[screenY][screenX] += 1.0
		}
	}

	maxDensity := 0.0
	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < cfg.Width; x++ {
			if density[y][x] > maxDensity {
				maxDensity = density[y][x]
			}
		}
	}

	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < cfg.Width; x++ {
			if density[y][x] == 0 {
				grid[y][x] = 0
			} else {
				val := math.Log(density[y][x]) / math.Log(maxDensity)
				colorIdx := int(val * float64(len(cfg.Colors)))
				if colorIdx < 0 { colorIdx = 0 }
				if colorIdx >= len(cfg.Colors) { colorIdx = len(cfg.Colors) - 1 }
				grid[y][x] = colorIdx
			}
		}
	}
	return grid
}