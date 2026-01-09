package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// LoadColors loads the color palette from a file.
func LoadColors(filePath string) error {
	palette.Colors = make(map[string]rl.Color)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not load colors from file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("%s:%d WARNING: Invalid line format\n", filePath, lineNumber)
			continue
		}

		key := parts[0]
		h, errH := strconv.ParseFloat(parts[1], 64)
		s, errS := strconv.ParseFloat(parts[2], 64)
		v, errV := strconv.ParseFloat(parts[3], 64)

		if errH != nil || errS != nil || errV != nil {
			fmt.Printf("%s:%d WARNING: Invalid HSV value\n", filePath, lineNumber)
			continue
		}

		hue := float32(h / 255.0 * 360.0)
		saturation := float32(s / 255.0)
		value := float32(v / 255.0)

		palette.Colors[key] = rl.ColorFromHSV(hue, saturation, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading colors file: %w", err)
	}

	return nil
}
