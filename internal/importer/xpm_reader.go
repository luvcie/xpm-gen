package importer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type XPMData struct {
	Width         int
	Height        int
	NumColors     int
	CharsPerPixel int
	Colors        map[string]string // char -> hex color
	PaletteKeys   []string          // ordered list of chars (to preserve order)
	Pixels        []string          // raw pixel rows
}

// parses a simple xpm file
func ReadXPM(filename string) (*XPMData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var contentLines []string

	// regex to find content inside double quotes
	re := regexp.MustCompile(`"([^"]+)"`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			contentLines = append(contentLines, matches[1])
		}
	}

	if len(contentLines) == 0 {
		return nil, fmt.Errorf("no xpm data found in file")
	}

	// 1. header: "width height num_colors chars_per_pixel"
	headerParts := strings.Fields(contentLines[0])
	if len(headerParts) < 4 {
		return nil, fmt.Errorf("invalid xpm header: %s", contentLines[0])
	}

	w, _ := strconv.Atoi(headerParts[0])
	h, _ := strconv.Atoi(headerParts[1])
	nc, _ := strconv.Atoi(headerParts[2])
	cpp, _ := strconv.Atoi(headerParts[3])

	data := &XPMData{
		Width:         w,
		Height:        h,
		NumColors:     nc,
		CharsPerPixel: cpp,
		Colors:        make(map[string]string),
		PaletteKeys:   make([]string, 0, nc),
		Pixels:        make([]string, h),
	}

	// 2. palette
	// we expect 'nc' lines of palette definition
	// format usually: "char c color"
	for i := 0; i < nc; i++ {
		if 1+i >= len(contentLines) {
			return nil, fmt.Errorf("unexpected end of file reading palette")
		}
		entry := contentLines[1+i]

		// the char code is the first 'cpp' characters
		if len(entry) < cpp {
			return nil, fmt.Errorf("palette entry too short: %s", entry)
		}
		charCode := entry[:cpp]
		
		// find the color definition. we look for "c " (color visual)
		// this is a simplification; xpm can have 'm', 'g', 'g4' etc.
		// but for our generator's output, 'c' is standard.
		parts := strings.Split(entry[cpp:], "c ")
		if len(parts) < 2 {
			// fallback: try splitting by spaces and taking the last part
			fields := strings.Fields(entry)
			if len(fields) >= 2 {
				color := fields[len(fields)-1]
				data.Colors[charCode] = color
			} else {
				return nil, fmt.Errorf("could not parse palette entry: %s", entry)
			}
		} else {
			color := strings.TrimSpace(parts[1])
			data.Colors[charCode] = color
		}
		data.PaletteKeys = append(data.PaletteKeys, charCode)
	}

	// 3. pixels
	for i := 0; i < h; i++ {
		idx := 1 + nc + i
		if idx >= len(contentLines) {
			return nil, fmt.Errorf("not enough pixel data")
		}
		data.Pixels[i] = contentLines[idx]
	}

	return data, nil
}