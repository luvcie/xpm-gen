package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"xpm-gen/internal/config"
	"xpm-gen/internal/exporter"
	"xpm-gen/internal/generator"
	"xpm-gen/internal/importer"
)

// version can be injected at build time via -ldflags
var Version = "v1.0-dev"

// returns an ansi string that paints the background with the given hex color
func colorBlock(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return ""
	}
	r, _ := strconv.ParseInt(hex[1:3], 16, 64)
	g, _ := strconv.ParseInt(hex[3:5], 16, 64)
	b, _ := strconv.ParseInt(hex[5:7], 16, 64)
	// ansi truecolor background: \033[48;2;r;g;bm
	return fmt.Sprintf("\033[48;2;%d;%d;%dm      \033[0m", r, g, b)
}

// generates random hex palette
// takes: size n
// returns: slice of hex strings
func generateRandomPalette(n int) []string {
	palette := make([]string, n)
	for i := 0; i < n; i++ {
		r := rand.Intn(256)
		g := rand.Intn(256)
		b := rand.Intn(256)
		palette[i] = fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return palette
}

// main entry point
// orchestrates configuration, generation, and saving
func main() {
	rand.Seed(time.Now().UnixNano())

	// cli flags setup
	widthPtr := flag.Int("w", 128, "Width of the texture")
	heightPtr := flag.Int("h", 128, "Height of the texture")
	algoPtr := flag.String("algo", "xor", "Algorithm: 'noise', 'xor', 'circles', 'mandelbrot', 'julia', 'melting', 'creature', 'pastel', 'attractor', 'cute', 'cutebunny', 'physarum', 'coral', 'random'")
	randColorsPtr := flag.Bool("randcolors", false, "Randomize the color palette")
	randomGenPtr := flag.Bool("random", false, "Generate a unique random algorithm")
	recolorPtr := flag.String("recolor", "", "Recolor an existing XPM file (interactive)")
	pngPtr := flag.Bool("png", false, "Convert output to PNG (requires ImageMagick)")
	versionPtr := flag.Bool("version", false, "Print version information")

	// custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "xpm-gen: advanced procedural texture synthesizer\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  xpm-gen [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// check version first
	if *versionPtr {
		fmt.Printf("xpm-gen %s\n", Version)
		os.Exit(0)
	}

	// recolor mode
	if *recolorPtr != "" {
		fmt.Printf("Reading %s...\n", *recolorPtr)
		data, err := importer.ReadXPM(*recolorPtr)
		if err != nil {
			fmt.Printf("Error reading XPM: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Recoloring %s (%dx%d, %d colors)\n", *recolorPtr, data.Width, data.Height, data.NumColors)
		newColors := make([]string, len(data.PaletteKeys))

		// initialize readline
		rl, err := readline.New("")
		if err != nil {
			fmt.Printf("Error initializing readline: %v\n", err)
			os.Exit(1)
		}
		defer rl.Close()

		for i, charCode := range data.PaletteKeys {
			oldColor := data.Colors[charCode]
			preview := colorBlock(oldColor)
			
			prompt := fmt.Sprintf("Color %d: %s %s (mapped to '%s') -> New hex color (e.g., #FF0000) [Enter to keep]: ", i+1, preview, oldColor, charCode)
			rl.SetPrompt(prompt)
			
			line, err := rl.Readline()
			if err != nil { // eof or ctrl+c
				break
			}
			
			input := strings.TrimSpace(line)
			if input == "" {
				newColors[i] = oldColor
			} else {
				// auto-prepend '#' if missing
				if !strings.HasPrefix(input, "#") && len(input) == 6 {
					input = "#" + input
				}
				newColors[i] = input
			}
		}

		// reconstruct grid
		fmt.Println("Reconstructing grid...")
		charMap := make(map[string]int)
		for i, k := range data.PaletteKeys {
			charMap[k] = i
		}

		grid := make([][]int, data.Height)
		for y, row := range data.Pixels {
			grid[y] = make([]int, data.Width)
			for x := 0; x < data.Width; x++ {
				start := x * data.CharsPerPixel
				end := start + data.CharsPerPixel
				if end > len(row) {
					// safe fallback
				end = len(row)
				}
			char := row[start:end]
			if idx, ok := charMap[char]; ok {
				grid[y][x] = idx
			} else {
				grid[y][x] = 0 // default to 0 if unknown
			}
			}
		}

		// create config for exporter
		cfg := config.Config{
			Width:     data.Width,
			Height:    data.Height,
			Algorithm: "recolored",
			Colors:    newColors,
			Chars:     data.PaletteKeys,
		}

		// export
		xpmContent := exporter.GridToXPM(grid, cfg)
		// we'll use the original filename base + _recolored
		fileName := exporter.SaveUniqueFile("recolored", xpmContent)
		fmt.Printf("Success! Generated %s\n", fileName)
		
		if *pngPtr {
			if err := exporter.ConvertToPNG(fileName); err != nil {
				fmt.Printf("Error converting to PNG: %v\n", err)
			} else {
				fmt.Printf("Success! PNG created.\n")
			}
		}

		os.Exit(0)
	}

	// validation
	validAlgos := map[string]bool{
		"noise": true, "xor": true, "circles": true,
		"mandelbrot": true, "julia": true, "melting": true,
		"creature": true, "pastel": true, "attractor": true,
		"cute": true, "cutebunny": true, "physarum": true, "coral": true,
		"random_gen": true, // allow our internal name
	}

	if *algoPtr == "random" {
		keys := make([]string, 0, len(validAlgos))
		for k := range validAlgos {
			if k != "random_gen" { // exclude the internal one from selection
				keys = append(keys, k)
			}
		}
		*algoPtr = keys[rand.Intn(len(keys))]
	}

	if !validAlgos[*algoPtr] && !*randomGenPtr {
		fmt.Printf("Error: Unknown algorithm '%s'\n", *algoPtr)
		os.Exit(1)
	}

	// palette setup
	colors := []string{"#000000", "#39FF14", "#FF69B4", "#00FFFF", "#FFFF00", "#BF00FF"}
	chars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p"}

	if *algoPtr == "creature" {
		colors = []string{"#000000", "#2b0000", "#660000", "#4a4a4a", "#e0e0e0", "#ffea00"}
	} else if *algoPtr == "pastel" {
		colors = []string{"#89CFF0", "#E6E6FA", "#98FF98", "#FFD1DC", "#FFDAB9", "#FFFDD0"}
	} else if *algoPtr == "attractor" {
		colors = []string{"#000000", "#111122", "#004488", "#0088CC", "#00FFFF", "#FFFFFF"}
	} else if *algoPtr == "cute" {
		// procedural color harmony (hsv)
		baseHue := float64(rand.Intn(360))

		// body: base hue, low sat (50), high val (95) -> gives us that pastel look
		bodyColor := hsvToHex(baseHue, 50, 95)

		// eyes: complementary hue (+180), high sat (80), med val (50) -> high contrast to pop out
		eyeHue := math.Mod(baseHue+180, 360)
		eyeColor := hsvToHex(eyeHue, 80, 50)

		// background: transparent
		bgColor := "None"

		colors = []string{bgColor, bodyColor, eyeColor}
	} else if *algoPtr == "cutebunny" {
		// soft whites, pinks, browns
		palettes := [][]string{
			{"None", "#FFFFFF", "#FF69B4"}, // white bunny, pink eyes
			{"None", "#FFC0CB", "#000000"}, // pink bunny, black eyes
			{"None", "#D2B48C", "#5C4033"}, // brown bunny, dark eyes
			{"None", "#E6E6FA", "#4B0082"}, // lavender bunny, indigo eyes
		}
		colors = palettes[rand.Intn(len(palettes))]
	} else if *algoPtr == "physarum" {
		// generate a random neon gradient
		// black -> dark color -> bright color -> white
		baseHue := rand.Float64() * 360.0
		colors = make([]string, 16)
		colors[0] = "#000000" // background
		for i := 1; i < 16; i++ {
			// ramp up value and saturation
			t := float64(i) / 15.0
			// hue shifts slightly for interest
			h := math.Mod(baseHue + (t * 30.0), 360.0)
			s := 100.0 - (t * 20.0) // desaturate slightly towards white
			v := 30.0 + (t * 70.0)  // get brighter
			
			// push the last few colors to pure white
			if i > 13 { s = 0; v = 100 }
			
			colors[i] = hsvToHex(h, s, v)
		}
	} else if *algoPtr == "coral" {
		// electric blue / cyan / magenta gradient
		colors = []string{"#000000", "#000033", "#000066", "#000099", "#0000CC", "#0000FF", "#0055FF", "#00AAFF", "#00FFFF", "#55FFFF", "#AAFFFF", "#FFFFFF", "#FF00FF", "#FF55FF"}
	}

	if *randColorsPtr {
		colors = generateRandomPalette(6)
	}

	cfg := config.Config{
		Width:     *widthPtr,
		Height:    *heightPtr,
		Algorithm: *algoPtr,
		Colors:    colors,
		Chars:     chars,
	}

	var grid [][]int
	
	if *randomGenPtr {
		cfg.Algorithm = "random_gen"
		// generate a new random expression
		expr := generator.GenerateRandomExpression(5 + rand.Intn(5)) // depth 5-10
		algoString := expr.String()
		fmt.Printf("Generated Algorithm: %s\n", algoString)
		
		// save the algorithm to a file
		// use a timestamp to ensure uniqueness and match the image filename pattern approximately
		timestamp := time.Now().Unix()
		algoFilename := fmt.Sprintf("xpmgen_random_%d.algo", timestamp)
		if err := os.WriteFile(algoFilename, []byte(algoString), 0644); err != nil {
			fmt.Printf("Error saving algorithm file: %v\n", err)
		} else {
			fmt.Printf("Saved algorithm to %s\n", algoFilename)
		}

		// generate the grid using this expression
		grid = generator.GenerateFromExpression(cfg, expr)
	} else {
		fmt.Printf("Generating %dx%d texture using '%s'\n", cfg.Width, cfg.Height, cfg.Algorithm)
		// execute pipeline
		grid = generator.GenerateGrid(cfg)
	}

	xpmContent := exporter.GridToXPM(grid, cfg)
	fileName := exporter.SaveUniqueFile(cfg.Algorithm, xpmContent)

	fmt.Printf("Success! Generated %s\n", fileName)

	if *pngPtr {
		if err := exporter.ConvertToPNG(fileName); err != nil {
			fmt.Printf("Error converting to PNG: %v\n", err)
		} else {
			fmt.Printf("Success! PNG created.\n")
		}
	}
}

// just a helper to convert hsv values to a hex string
func hsvToHex(h, s, v float64) string {
	s /= 100
	v /= 100
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c
	var r, g, b float64
	if 0 <= h && h < 60 {
		r, g, b = c, x, 0
	} else if 60 <= h && h < 120 {
		r, g, b = x, c, 0
	} else if 120 <= h && h < 180 {
		r, g, b = 0, c, x
	} else if 180 <= h && h < 240 {
		r, g, b = 0, x, c
	} else if 240 <= h && h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}
	return fmt.Sprintf("#%02X%02X%02X", int((r+m)*255), int((g+m)*255), int((b+m)*255))
}
