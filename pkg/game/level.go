package game

import (
	"image"
	_ "image/png" // import the png decoder
	"os"

	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// LevelCell represents the different types of cells in a level file.
type LevelCell int

const (
	LevelNone LevelCell = iota
	LevelGnome
	LevelMother
	LevelGuard
	LevelFloor
	LevelWall
	LevelDoor
	LevelCheckpoint
	LevelBombRefill
	LevelBarricade
	LevelKey
	LevelPlayer
	LevelFather
	LevelBombSlot
)

// LevelCellColor maps level cell types to their corresponding colors.
var LevelCellColor = map[LevelCell]rl.Color{
	LevelNone:       rl.NewColor(0, 0, 0, 0),
	LevelGnome:      rl.NewColor(255, 150, 0, 255),
	LevelMother:     rl.NewColor(150, 255, 0, 255),
	LevelGuard:      rl.NewColor(0, 255, 0, 255),
	LevelFloor:      rl.NewColor(255, 255, 255, 255),
	LevelWall:       rl.NewColor(0, 0, 0, 255),
	LevelDoor:       rl.NewColor(0, 255, 255, 255),
	LevelCheckpoint: rl.NewColor(255, 0, 255, 255),
	LevelBombRefill: rl.NewColor(255, 0, 0, 255),
	LevelBarricade:  rl.NewColor(255, 0, 150, 255),
	LevelKey:        rl.NewColor(255, 255, 0, 255),
	LevelPlayer:     rl.NewColor(0, 0, 255, 255),
	LevelFather:     rl.NewColor(38, 95, 218, 255),
	LevelBombSlot:   rl.NewColor(188, 83, 83, 255),
}

// LoadGameFromImage loads a game state from an image file.
func LoadGameFromImage(filePath string, gs *State, updatePlayer bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	gs.Map = make([][]world.Cell, height)
	for i := range gs.Map {
		gs.Map[i] = make([]world.Cell, width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			color := rl.NewColor(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))

			levelCell := LevelNone
			for cell, cellColor := range LevelCellColor {
				if cellColor.R == color.R && cellColor.G == color.G && cellColor.B == color.B && cellColor.A == color.A {
					levelCell = cell
					break
				}
			}

			switch levelCell {
			case LevelFloor:
				gs.Map[y][x] = world.CellFloor
			case LevelWall:
				gs.Map[y][x] = world.CellWall
			case LevelDoor:
				gs.Map[y][x] = world.CellDoor
			case LevelBarricade:
				gs.Map[y][x] = world.CellBarricade
			case LevelCheckpoint:
				gs.Map[y][x] = world.CellFloor
				gs.AllocateItem(world.IVector2{X: x, Y: y}, entities.ItemCheckpoint)
			case LevelBombRefill:
				gs.Map[y][x] = world.CellFloor
				gs.AllocateItem(world.IVector2{X: x, Y: y}, entities.ItemBombRefill)
			case LevelBombSlot:
				gs.Map[y][x] = world.CellFloor
				gs.AllocateItem(world.IVector2{X: x, Y: y}, entities.ItemBombSlot)
			case LevelKey:
				gs.Map[y][x] = world.CellFloor
				gs.AllocateItem(world.IVector2{X: x, Y: y}, entities.ItemKey)
			case LevelGuard:
				gs.Map[y][x] = world.CellFloor
				gs.SpawnGuard(world.IVector2{X: x, Y: y})
			case LevelMother:
				gs.Map[y][x] = world.CellFloor
				gs.SpawnMother(world.IVector2{X: x, Y: y})
			case LevelGnome:
				gs.Map[y][x] = world.CellFloor
				gs.SpawnGnome(world.IVector2{X: x, Y: y})
			case LevelFather:
				gs.Map[y][x] = world.CellFloor
				// Father eepers will be implemented later
			case LevelPlayer:
				gs.Map[y][x] = world.CellFloor
				if updatePlayer {
					gs.Player.Position = world.IVector2{X: x, Y: y}
					// Initialize eyes target to look down (default direction)
					gs.Player.EyesTarget = world.IVector2{X: x, Y: y + 1}
				}
			default:
				gs.Map[y][x] = world.CellFloor
			}
		}
	}

	return nil
}
