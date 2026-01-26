package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"xpm-synth/internal/config"
	"xpm-synth/internal/exporter"
	"xpm-synth/internal/generator"
)

// Version can be injected at build time via -ldflags
var Version = "v1.0-dev"

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
	algoPtr := flag.String("algo", "xor", "Algorithm: 'noise', 'xor', 'circles', 'mandelbrot', 'julia', 'melting', 'creature', 'pastel', 'attractor'")
	randColorsPtr := flag.Bool("randcolors", false, "Randomize the color palette")
	pngPtr := flag.Bool("png", false, "Convert output to PNG (requires ImageMagick)")
	versionPtr := flag.Bool("version", false, "Print version information")

	// custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "xpm-synth: advanced procedural texture synthesizer\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  xpm-synth [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// check version first
	if *versionPtr {
		fmt.Printf("xpm-synth %s\n", Version)
		os.Exit(0)
	}

	// validation
	validAlgos := map[string]bool{
		"noise": true, "xor": true, "circles": true,
		"mandelbrot": true, "julia": true, "melting": true,
		"creature": true, "pastel": true, "attractor": true,
	}

	if !validAlgos[*algoPtr] {
		fmt.Printf("Error: Unknown algorithm '%s'\n", *algoPtr)
		os.Exit(1)
	}

	// palette setup
	colors := []string{"#000000", "#39FF14", "#FF69B4", "#00FFFF", "#FFFF00", "#BF00FF"}
	chars := []string{"a", "b", "c", "d", "e", "f"}

	if *algoPtr == "creature" {
		colors = []string{"#000000", "#2b0000", "#660000", "#4a4a4a", "#e0e0e0", "#ffea00"}
	} else if *algoPtr == "pastel" {
		colors = []string{"#89CFF0", "#E6E6FA", "#98FF98", "#FFD1DC", "#FFDAB9", "#FFFDD0"}
	} else if *algoPtr == "attractor" {
		colors = []string{"#000000", "#111122", "#004488", "#0088CC", "#00FFFF", "#FFFFFF"}
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

	fmt.Printf("Generating %dx%d texture using '%s'\n", cfg.Width, cfg.Height, cfg.Algorithm)

	// execute pipeline
	grid := generator.GenerateGrid(cfg)
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