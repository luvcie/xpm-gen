package generator

import (
	"math"
	"math/rand"
	"xpm-gen/internal/config"
)

type Point struct {
	x, y, r float64
}

// doing the metaballs thing for blobs and neoteny for the cute faces.
func runCuteGenerator(cfg config.Config) [][]int {
	grid := make([][]int, cfg.Height)
	for i := range grid {
		grid[i] = make([]int, cfg.Width)
	}

	// 1. spawn some metaballs (the hearts of the creature)
	// random locations, but mirrored across the y-axis so it looks symmetric
	numHearts := 3 + rand.Intn(3) // 3 to 5
	balls := []Point{}

	centerX := float64(cfg.Width) / 2.0
	// keep them somewhat in the middle so they don't drift off screen
	spawnWidth := float64(cfg.Width) * 0.4
	spawnHeight := float64(cfg.Height) * 0.6
	spawnOffsetY := float64(cfg.Height) * 0.2

	for i := 0; i < numHearts; i++ {
		// spawn on the left (or center)
		px := (centerX - spawnWidth/2) + rand.Float64()*spawnWidth
		py := spawnOffsetY + rand.Float64()*spawnHeight
		
		// random size
		// scale it based on the image size, like 5-15% of the width
		minR := float64(cfg.Width) * 0.05
		maxR := float64(cfg.Width) * 0.15
		r := minR + rand.Float64()*(maxR-minR)

		balls = append(balls, Point{px, py, r})
		
		// mirror logic:
		// just mirroring everything to make sure it's perfectly symmetric
		mx := centerX + (centerX - px)
		balls = append(balls, Point{mx, py, r})
	}

	// 2. render the threshold (the metaballs field)
	minY := cfg.Height
	maxY := -1
	
	// pre-calculate squares to save on sqrt calls
	// field function: sum( r^2 / dist^2 )
	for y := 0; y < cfg.Height; y++ {
		for x := 0; x < cfg.Width; x++ {
			influence := 0.0
			fx, fy := float64(x), float64(y)
			
			for _, b := range balls {
				distSq := (fx-b.x)*(fx-b.x) + (fy-b.y)*(fy-b.y)
				if distSq < 1.0 {
					distSq = 1.0 // avoid dividing by zero
				}
				influence += (b.r * b.r) / distSq
			}

			// threshold
			if influence > 1.2 { // magic number to tune how blobby it is
				grid[y][x] = 1 // body color (index 1)
				
				// track bounds so we know where the head is
				if y < minY { minY = y }
				if y > maxY { maxY = y }
			} else {
				grid[y][x] = 0 // background (index 0)
			}
		}
	}

	// if nothing got drawn, just bail out (unlikely though)
	if maxY == -1 {
		return grid
	}

	// 3. neoteny ratio (making the face look cute)
	// head is the top part. assuming the whole blob is the body.
	// we want eyes "low" on the head to look like a baby.
	
	creatureHeight := maxY - minY
	
	// "vertical: the eyes must be located below the vertical center line of the head."
	// putting eyes at roughly 45% down from the top.
	eyeY := minY + int(float64(creatureHeight)*0.45)
	
	// "horizontal: the distance between eyes should be relatively wide."
	eyeSpacing := int(float64(cfg.Width) * 0.12)
	
	lx := int(centerX) - eyeSpacing
	rx := int(centerX) + eyeSpacing
	
	// eye size
	eyeRadius := int(float64(cfg.Width) * 0.03)
	if eyeRadius < 1 { eyeRadius = 1 }

	drawEye(grid, lx, eyeY, eyeRadius, cfg)
	drawEye(grid, rx, eyeY, eyeRadius, cfg)

	return grid
}

func drawEye(grid [][]int, cx, cy, r int, cfg config.Config) {
	// simple filled circle
	// using color index 2 for the eyes
	colorIdx := 2
	
	for y := cy - r; y <= cy + r; y++ {
		for x := cx - r; x <= cx + r; x++ {
			if x >= 0 && x < cfg.Width && y >= 0 && y < cfg.Height {
				dist := math.Sqrt(float64((x-cx)*(x-cx) + (y-cy)*(y-cy)))
				if dist <= float64(r) {
					grid[y][x] = colorIdx
				}
			}
		}
	}
}
