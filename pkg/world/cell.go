package world

import (
	"github.com/engpetarmarinov/eepers-go/pkg/palette"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Cell represents the type of cell in the game map.
type Cell int

const (
	CellNone Cell = iota
	CellFloor
	CellWall
	CellDoor
	CellBarricade
	CellExplosion
)

// CellColor returns the color for a given cell type.
func CellColor(c Cell) rl.Color {
	switch c {
	case CellNone:
		return palette.Colors["COLOR_BACKGROUND"]
	case CellFloor:
		return palette.Colors["COLOR_FLOOR"]
	case CellWall:
		return palette.Colors["COLOR_WALL"]
	case CellDoor:
		return palette.Colors["COLOR_DOOR"]
	case CellBarricade:
		return palette.Colors["COLOR_BARRICADE"]
	case CellExplosion:
		return palette.Colors["COLOR_EXPLOSION"]
	default:
		return rl.Black
	}
}
