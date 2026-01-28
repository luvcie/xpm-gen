package generator

import (
	"math/rand"
	"xpm-gen/internal/config"
	"github.com/schollz/progressbar/v3"
)

// gray-scott reaction diffusion simulation
// generates biological patterns like coral, fingerprints, and spots
func runCoral(cfg config.Config) [][]int {
	width, height := cfg.Width, cfg.Height
	
	// grids for chemicals A and B
	// a = feed, b = kill
	gridA := make([][]float64, height)
	gridB := make([][]float64, height)
	nextA := make([][]float64, height)
	nextB := make([][]float64, height)
	
	for y := 0; y < height; y++ {
		gridA[y] = make([]float64, width)
		gridB[y] = make([]float64, width)
		nextA[y] = make([]float64, width)
		nextB[y] = make([]float64, width)
		
		for x := 0; x < width; x++ {
			gridA[y][x] = 1.0 // fill world with 'feed'
			// heavy noise seeding to ensure it doesn't die out
			if rand.Float64() < 0.10 { 
				gridB[y][x] = 1.0
			} else {
				gridB[y][x] = 0.0
			}
		}
	}
	
	// reverted to "Standard Coral" parameters which are very robust
	// these are guaranteed to grow and fill the screen
	feed := 0.0545
	k := 0.062
	diffA := 1.0
	diffB := 0.5
	
	steps := 1000
	bar := progressbar.Default(int64(steps), "growing coral")

	for step := 0; step < steps; step++ {
		bar.Add(1)
		
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				a := gridA[y][x]
				b := gridB[y][x]
				
				// laplacian (diffusion) using 3x3 convolution
				lapA := 0.0
				lapB := 0.0
				
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						ny := (y + dy + height) % height
						nx := (x + dx + width) % width
						
						weight := 0.0
						if dx == 0 && dy == 0 {
							weight = -1.0
						} else if dx == 0 || dy == 0 {
							weight = 0.2
						} else {
							weight = 0.05
						}
						
						lapA += gridA[ny][nx] * weight
						lapB += gridB[ny][nx] * weight
					}
				}
				
				// reaction-diffusion formula
				abb := a * b * b
				
				newA := a + (diffA * lapA) - abb + (feed * (1.0 - a))
				newB := b + (diffB * lapB) + abb - ((k + feed) * b)
				
				// clamp
				if newA < 0 { newA = 0 }
				if newA > 1 { newA = 1 }
				if newB < 0 { newB = 0 }
				if newB > 1 { newB = 1 }
				
				nextA[y][x] = newA
				nextB[y][x] = newB
			}
		}
		
		for y := 0; y < height; y++ {
			copy(gridA[y], nextA[y])
			copy(gridB[y], nextB[y])
		}
	}

	// convert (visualize B concentration)
	outGrid := make([][]int, height)
	numColors := len(cfg.Colors)
	
	for y := 0; y < height; y++ {
		outGrid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			val := gridB[y][x]
			
			// lower visual threshold so we can see the pattern even if it's faint
			// val is usually 0.0 to 0.4
			
			// normalize typical range
			val = val * 2.5 
			
			if val < 0.05 { 
				val = 0 
			}
			
			if val > 1.0 { val = 1.0 }
			
			idx := int(val * float64(numColors))
			if idx >= numColors { idx = numColors - 1 }
			if idx < 0 { idx = 0 }
			outGrid[y][x] = idx
		}
	}

	return outGrid
}
