package generator

import (
	"math/rand"
	"xpm-gen/internal/config"
)

// basically the cute generator but with guaranteed long ears
func runCuteBunnyGenerator(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	for i := range grid {
		grid[i] = make([]int, cfg.Width)
	}

	balls := []Point{}
	centerX := float64(cfg.Width) / 2.0
	
	// 1. the body (just like cute.go)
	// mostly concentrated in bottom half
	numBodyParts := 3 + rand.Intn(3)
	
	spawnWidth := float64(cfg.Width) * 0.4
	spawnHeight := float64(cfg.Height) * 0.4
	spawnOffsetY := float64(cfg.Height) * 0.4 // lower down

	for i := 0; i < numBodyParts; i++ {
		px := (centerX - spawnWidth/2) + rand.Float64()*spawnWidth
		py := spawnOffsetY + rand.Float64()*spawnHeight
		
		minR := float64(cfg.Width) * 0.08
		maxR := float64(cfg.Width) * 0.18
		r := minR + rand.Float64()*(maxR-minR)

		balls = append(balls, Point{px, py, r})
		mx := centerX + (centerX - px)
		balls = append(balls, Point{mx, py, r})
	}

	// 2. the ears (the important part)
	// we stack circles to make them long
	earLen := 3 + rand.Intn(3) // how many balls tall the ear is
	earBaseX := centerX - (float64(cfg.Width) * 0.15) // offset from center
	earBaseY := spawnOffsetY // start where body starts
	earRadius := float64(cfg.Width) * 0.06

	for i := 0; i < earLen; i++ {
		// go up!
		yPos := earBaseY - (float64(i) * earRadius * 1.5)
		// maybe tilt them out slightly?
		xPos := earBaseX - (float64(i) * earRadius * 0.2)
		
		balls = append(balls, Point{xPos, yPos, earRadius})
		
		// mirror ear
		mx := centerX + (centerX - xPos)
		balls = append(balls, Point{mx, yPos, earRadius})
	}

	// 3. render field
	minY := cfg.Height
	maxY := -1
	
	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < cfg.Width; x++ {
			influence := 0.0
			fx, fy := float64(x), float64(y)
			
			for _, b := range balls {
				distSq := (fx-b.x)*(fx-b.x) + (fy-b.y)*(fy-b.y)
				if distSq < 1.0 { distSq = 1.0 }
				influence += (b.r * b.r) / distSq
			}

			if influence > 1.2 {
				grid[y][x] = 1
				if y < minY { minY = y }
				if y > maxY { maxY = y }
			} else {
				grid[y][x] = 0
			}
		}
	}

	if maxY == -1 { return grid }

	// 4. face logic
	// bunny eyes need to be wider apart maybe?
	creatureHeight := maxY - minY
	
	// eyes roughly in the middle of the "head" part (not the ears)
	// let's estimate head top is below the ears
	// assuming ears take up top 30-40%
	headTop := minY + int(float64(creatureHeight)*0.3) 
	
	eyeY := headTop + int(float64(creatureHeight)*0.2)
	eyeSpacing := int(float64(cfg.Width) * 0.14)
	
	lx := int(centerX) - eyeSpacing
	rx := int(centerX) + eyeSpacing
	eyeRadius := int(float64(cfg.Width) * 0.025)
	if eyeRadius < 1 { eyeRadius = 1 }

	drawEye(grid, lx, eyeY, eyeRadius, cfg)
	drawEye(grid, rx, eyeY, eyeRadius, cfg)

	return grid
}
