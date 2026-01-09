package game

import (
	"github.com/engpetarmarinov/eepers-go/pkg/entities"
	"github.com/engpetarmarinov/eepers-go/pkg/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// State represents the entire state of the game.
type State struct {
	Map                [][]world.Cell
	Player             entities.PlayerState
	Items              []entities.Item
	Eepers             []entities.EeperState
	Bombs              []entities.BombState
	TurnAnimation      float32
	Camera             rl.Camera2D
	Tutorial           TutorialState
	DurationOfLastTurn float64
}
