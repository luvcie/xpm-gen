package config

// holds configuration for the texture generation
// width/height: dimensions of the output
// algorithm: selected generation method
// colors: palette of hex codes
// chars: xpm mapping characters
type Config struct {
	Width     int
	Height    int
	Algorithm string
	Colors    []string
	Chars     []string
}