package generator

import (
	"math"
	"math/rand"
	"xpm-gen/internal/config"
	"github.com/schollz/progressbar/v3"
)

type Agent struct {
	x, y  float64
	angle float64
}

// simulates physarum polycephalum (slime mold) behavior
// creates organic transport networks and vein-like structures
func runPhysarum(cfg config.Config) [][]int {
	width, height := cfg.Width, cfg.Height
	
	// 1. init simulation state
	trail := make([][]float64, height)
	nextTrail := make([][]float64, height)
	for y := 0; y < height; y++ {
		trail[y] = make([]float64, width)
		nextTrail[y] = make([]float64, width)
	}

	// 2. spawn agents uniformly (no voids)
	// standard density
	numAgents := int(float64(width*height) * 0.12)
	agents := make([]Agent, numAgents)
	
	for i := range agents {
		agents[i] = Agent{
			x: rand.Float64() * float64(width),
			y: rand.Float64() * float64(height),
			angle: rand.Float64() * 2 * math.Pi,
		}
	}

	// simulation parameters tuned for ULTRA THIN lines
	sensorAngle := 45.0 * (math.Pi / 180.0)
	sensorDist := 4.0 
	turnAngle := 45.0 * (math.Pi / 180.0)
	decayFactor := 0.9 
	depositAmount := 0.2 // very low deposit => only heavy traffic survives
	
	steps := 500
	
	bar := progressbar.Default(int64(steps), "simulating physarum")

	for step := 0; step < steps; step++ {
		bar.Add(1)

		// a. move and sense
		for i := range agents {
			a := &agents[i]
			
			l := sense(a, sensorDist, -sensorAngle, trail, width, height)
			c := sense(a, sensorDist, 0, trail, width, height)
			r := sense(a, sensorDist, sensorAngle, trail, width, height)
			
			if c > l && c > r {
				// straight
			} else if c < l && c < r {
				a.angle += (rand.Float64() - 0.5) * 2 * turnAngle
			} else if l > r {
				a.angle -= turnAngle
			} else if r > l {
				a.angle += turnAngle
			}
			
			nextX := a.x + math.Cos(a.angle)
			nextY := a.y + math.Sin(a.angle)
			
			if nextX < 0 { nextX += float64(width) }
			if nextX >= float64(width) { nextX -= float64(width) }
			if nextY < 0 { nextY += float64(height) }
			if nextY >= float64(height) { nextY -= float64(height) }
			
			a.x = nextX
			a.y = nextY
			
			ix, iy := int(nextX), int(nextY)
			trail[iy][ix] += depositAmount
			if trail[iy][ix] > 1.0 { trail[iy][ix] = 1.0 }
		}
		
		// b. diffuse and decay
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				sum := 0.0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						ny := (y + dy + height) % height
						nx := (x + dx + width) % width
						sum += trail[ny][nx]
					}
				}
				avg := sum / 9.0
				nextTrail[y][x] = avg * decayFactor
			}
		}
		for y := 0; y < height; y++ {
			copy(trail[y], nextTrail[y])
		}
	}

	// convert
	grid := make([][]int, height)
	numColors := len(cfg.Colors)
	
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		for x := 0; x < width; x++ {
			val := trail[y][x]
			
			// sharp threshold: cutoff at 0.2 to make lines thin
			if val < 0.2 { 
				val = 0 
			} else { 
				// stretch remainder
				val = (val - 0.2) * 1.25 
			}
			
			idx := int(val * float64(numColors))
			if idx >= numColors { idx = numColors - 1 }
			if idx < 0 { idx = 0 }
			grid[y][x] = idx
		}
	}

	return grid
}

func sense(a *Agent, dist, angleOffset float64, trail [][]float64, w, h int) float64 {
	sensorAngle := a.angle + angleOffset
	sx := a.x + math.Cos(sensorAngle)*dist
	sy := a.y + math.Sin(sensorAngle)*dist
	
	ix := (int(sx) + w) % w
	iy := (int(sy) + h) % h
	
	return trail[iy][ix]
}