package exporter

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"xpm-gen/internal/config"
)

// converts grid to xpm string
// formats header, palette, and pixel data for xpm3 standard
// takes: grid (2d array), config
// returns: xpm content string
func GridToXPM(grid [][]int, cfg config.Config) string {
	header := "/* XPM */\n"
	header += "static char * texture[] = {\n"
	header += fmt.Sprintf("\"%d %d %d %d\",\n", cfg.Width, cfg.Height, len(cfg.Colors), 1)

	for i, color := range cfg.Colors {
		header += fmt.Sprintf("\"%s c %s\",\n", cfg.Chars[i], color)
	}

	for y := 0; y < cfg.Height; y++ {
		line := "\""
		for x := 0; x < cfg.Width; x++ {
			line += cfg.Chars[grid[y][x]]
		}
		line += "\",\n"
		header += line
	}

	header += "};\n"
	return header
}

// saves content to a unique filename
// checks existing files to prevent overwrites, increments counter
// takes: algorithm name, content string
// returns: filename used
// mutates: filesystem (creates new .xpm file)
func SaveUniqueFile(algo string, content string) string {
	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("trippy_%s_%d.xpm", algo, i)
		if _, err := os.Stat(name); os.IsNotExist(err) {
			err := os.WriteFile(name, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error writing file: %v\n", err)
				os.Exit(1)
			}
			return name
		}
	}
	timestamp := time.Now().Unix()
	name := fmt.Sprintf("trippy_%s_%d.xpm", algo, timestamp)
	os.WriteFile(name, []byte(content), 0644)
	return name
}

// converts xpm to png using imagemagick
// calls external 'convert' command
// takes: filename string
// returns: error or nil
// mutates: filesystem (creates new .png file via subprocess)
func ConvertToPNG(fileName string) error {
	pngName := fileName[:len(fileName)-4] + ".png"
	cmd := exec.Command("convert", fileName, pngName)
	return cmd.Run()
}